package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"avito/internal/cerr"
	"avito/internal/config"
	"avito/internal/gen"
)

var (
	Host           string
	RequestTimeout = 10 * time.Second
	BasePath       string
)

var (
	basePathPR    string
	basePathUsers string
	basePathTeam  string
)

func init() {
	cfg := config.InitConfig()
	Host = fmt.Sprintf("%v:%v", cfg.ServiceHost, cfg.ServicePort)
	BasePath = "http://" + Host
	basePathPR = BasePath + "/pullRequest"
	basePathUsers = BasePath + "/users"
	basePathTeam = BasePath + "/team"
}

func DoWebRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	return resp, nil
}

func GetError(err cerr.ErrorType) gen.ErrorResponse {
	_, genErr := cerr.HandleErrs(cerr.CustomError{ErrType: err})

	return genErr
}

type TestData struct {
	teamForTest      *gen.Team                          //nolint:unused
	prForTest        *gen.PostPullRequestCreateJSONBody //nolint:unused
	mergeForTest     *gen.PostPullRequestMergeJSONBody  //nolint:unused
	setActiveForTest *gen.PostUsersSetIsActiveJSONBody  //nolint:unused
	path             string                             //nolint:unused
	description      string                             //nolint:unused
	body             any                                //nolint:unused
	expectedCode     int                                //nolint:unused
	expectedBody     any                                //nolint:unused
}

func CreateTeamForTest(team *gen.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	jsonData, err := json.Marshal(team)
	if err != nil {
		return err
	}

	resp, err := DoWebRequest(ctx, http.MethodPost, basePathTeam+"/add", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func CreatePRForTest(pr *gen.PostPullRequestCreateJSONBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	jsonData, err := json.Marshal(pr)
	if err != nil {
		return err
	}

	resp, err := DoWebRequest(ctx, http.MethodPost, basePathPR+"/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func MergePRForTest(pr *gen.PostPullRequestMergeJSONBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	jsonData, err := json.Marshal(pr)
	if err != nil {
		return err
	}

	resp, err := DoWebRequest(ctx, http.MethodPost, basePathPR+"/merge", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func SetIsActiveForTest(user *gen.PostUsersSetIsActiveJSONBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := DoWebRequest(ctx, http.MethodPost, basePathUsers+"/setIsActive", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
