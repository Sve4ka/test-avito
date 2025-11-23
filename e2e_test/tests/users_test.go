//go:build e2e

package tests

import (
	"avito/internal/cerr"
	"avito/internal/gen"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

// TestSetIsActive test /user/setIsActive
func TestSetIsActive(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			teamForTest: &gen.Team{
				TeamName: "testSetIsActiveSuccess",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testSetIsActiveSuccess",
						UserId:   "testSetIsActiveSuccess",
					},
				},
			},
			path:        basePathUsers + "/setIsActive",
			description: "Set Is Active Success",

			body: gen.PostUsersSetIsActiveJSONBody{
				IsActive: false,
				UserId:   "testSetIsActiveSuccess",
			},
			expectedCode: http.StatusOK,
			expectedBody: gen.PostUsersSetIsActive200JSONResponse{
				User: &gen.User{
					UserId:   "testSetIsActiveSuccess",
					Username: "testSetIsActiveSuccess",
					TeamName: "testSetIsActiveSuccess",
					IsActive: false,
				},
			},
		},
		{
			path:        basePathUsers + "/setIsActive",
			description: "Set Is Active NotFound",

			body: gen.PostUsersSetIsActiveJSONBody{
				IsActive: false,
				UserId:   "testSetIsActiveNotFound",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: gen.PostUsersSetIsActive404JSONResponse{
				Error: GetError(cerr.NOT_FOUND).Error,
			},
		},
	}
	var errorData gen.ErrorResponse

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.teamForTest != nil {
				err := CreateTeamForTest(test.teamForTest)
				if err != nil {
					t.Fatalf("error: %s", err)
				}
			}

			jsonData, err := json.Marshal(test.body)
			if err != nil {
				t.Fatalf("error: %s", err)
			}

			resp, err := DoWebRequest(ctx, http.MethodPost, test.path, bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("error: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				err = json.NewDecoder(resp.Body).Decode(&errorData)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}

				t.Errorf("got %d : %v", resp.StatusCode, errorData)
			} else {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				expectedBytes, err := json.Marshal(test.expectedBody)
				require.NoError(t, err)

				assert.JSONEq(t, string(expectedBytes), string(bodyBytes))
			}
		})
	}
}

// TestGetReview test /user/getReview
func TestGetReview(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			teamForTest: &gen.Team{
				TeamName: "testGetReviewSuccessNil",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testGetReviewSuccessNil",
						UserId:   "testGetReviewSuccessNil",
					},
				},
			},
			path:         basePathUsers + "/getReview?user_id=",
			description:  "Get Review Success 0 PR",
			body:         "testGetReviewSuccessNil",
			expectedCode: http.StatusOK,
			expectedBody: gen.GetUsersGetReview200JSONResponse{
				UserId:       "testGetReviewSuccessNil",
				PullRequests: []gen.PullRequestShort{},
			},
		},
		{
			teamForTest: &gen.Team{
				TeamName: "testGetReviewSuccess",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testGetReviewSuccess",
						UserId:   "testGetReviewSuccess",
					},
					{
						IsActive: true,
						Username: "testGetReviewSuccess",
						UserId:   "testGetReviewSuccess_1",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "testGetReviewSuccess",
				PullRequestName: "testGetReviewSuccess",
				PullRequestId:   "testGetReviewSuccess",
			},
			path:        basePathUsers + "/getReview?user_id=",
			description: "Get Review Success 0 PR",

			body:         "testGetReviewSuccess_1",
			expectedCode: http.StatusOK,
			expectedBody: gen.GetUsersGetReview200JSONResponse{
				UserId: "testGetReviewSuccess_1",
				PullRequests: []gen.PullRequestShort{
					{
						AuthorId:        "testGetReviewSuccess",
						PullRequestName: "testGetReviewSuccess",
						PullRequestId:   "testGetReviewSuccess",
						Status:          gen.PullRequestShortStatusOPEN,
					},
				},
			},
		},
		{
			path:        basePathUsers + "/getReview?user_id=",
			description: "Get Review Success NotFound",

			body:         "testGetReviewNotFound",
			expectedCode: http.StatusNotFound,
			expectedBody: gen.GetTeamGet404JSONResponse{
				Error: GetError(cerr.NOT_FOUND).Error,
			},
		},
	}
	var errorData gen.ErrorResponse

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.teamForTest != nil {
				err := CreateTeamForTest(test.teamForTest)
				if err != nil {
					t.Fatalf("error: %s", err)
				}
			}
			if test.prForTest != nil {
				err := CreatePRForTest(test.prForTest)
				if err != nil {
					t.Fatalf("error: %s", err)
				}
			}
			jsonData, err := json.Marshal(test.body)
			if err != nil {
				t.Fatalf("error: %s", err)
			}

			resp, err := DoWebRequest(ctx, http.MethodGet, test.path+test.body.(string), bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("error: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Log(resp.StatusCode)
				err = json.NewDecoder(resp.Body).Decode(&errorData)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}

				t.Errorf("got %d : %v", resp.StatusCode, errorData)
			} else {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				expectedBytes, err := json.Marshal(test.expectedBody)
				require.NoError(t, err)

				assert.JSONEq(t, string(expectedBytes), string(bodyBytes))
			}
		})
	}
}
