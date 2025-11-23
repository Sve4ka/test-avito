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

// TestAdd test /team/add
func TestAdd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			path:        basePathTeam + "/add",
			description: "Add success",
			body: gen.Team{
				TeamName: "testAddSuccess",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testAddSuccess",
						UserId:   "testAddSuccess_1",
					},
				},
			},
			expectedCode: http.StatusCreated,
			expectedBody: gen.PostTeamAdd201JSONResponse{
				Team: &gen.Team{
					TeamName: "testAddSuccess",
					Members: []gen.TeamMember{
						{
							IsActive: true,
							Username: "testAddSuccess",
							UserId:   "testAddSuccess_1",
						},
					},
				},
			},
		}, {
			path:        basePathTeam + "/add",
			description: "Add team_name already exists",
			teamForTest: &gen.Team{
				TeamName: "testAddTeamNameExist",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testAddTeamNameExist",
						UserId:   "testAddTeamNameExist_1",
					},
				},
			},
			body: gen.Team{
				TeamName: "testAddTeamNameExist",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testAddTeamNameExist",
						UserId:   "testAddTeamNameExist_1",
					},
				},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gen.PostTeamAdd400JSONResponse{
				Error: GetError(cerr.TEAM_EXISTS).Error,
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

// TestGet test /team/get
func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			path:        basePathTeam + "/get?team_name=",
			description: "Get team_name Success",
			teamForTest: &gen.Team{
				TeamName: "testGetTeamNameSuccess",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testGetTeamNameSuccess",
						UserId:   "testGetTeamNameSuccess_1",
					},
				},
			},
			body:         "testGetTeamNameSuccess",
			expectedCode: http.StatusOK,
			expectedBody: gen.GetTeamGet200JSONResponse{
				Members: []gen.TeamMember{
					{
						IsActive: true,
						Username: "testGetTeamNameSuccess",
						UserId:   "testGetTeamNameSuccess_1",
					},
				},
				TeamName: "testGetTeamNameSuccess",
			},
		},
		{
			path:         basePathTeam + "/get?team_name=",
			description:  "Get team_name NotFound",
			body:         "testGetTeamNameNotFound",
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

			resp, err := DoWebRequest(ctx, http.MethodGet, test.path+test.body.(string), nil)
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
