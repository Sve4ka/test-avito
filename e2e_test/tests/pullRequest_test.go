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

// TestCreate test/pullRequest/create
func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			teamForTest: &gen.Team{
				TeamName: "TestCreatePRSuccess_2assign",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_2assign_1",
						Username: "TestCreatePRSuccess_2assign",
					},
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_2assign_2",
						Username: "TestCreatePRSuccess_2assign",
					},
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_2assign_3",
						Username: "TestCreatePRSuccess_2assign",
					},
				},
			},
			path:        basePathPR + "/create",
			description: "Create success 2 assign",
			body: gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRSuccess_2assign_1",
				PullRequestId:   "TestCreatePRSuccess_2assign",
				PullRequestName: "TestCreatePRSuccess_2assign",
			},
			expectedCode: http.StatusCreated,
			expectedBody: gen.PostPullRequestCreate201JSONResponse{
				Pr: &gen.PullRequest{
					PullRequestId:   "TestCreatePRSuccess_2assign",
					PullRequestName: "TestCreatePRSuccess_2assign",
					AuthorId:        "TestCreatePRSuccess_2assign_1",
					Status:          gen.PullRequestStatusOPEN,
					AssignedReviewers: []string{
						"TestCreatePRSuccess_2assign_2",
						"TestCreatePRSuccess_2assign_3",
					},
				},
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestCreatePRSuccess_1assign_2users",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_1assign_2users_1",
						Username: "TestCreatePRSuccess_1assign_2users",
					},
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_1assign_2users_2",
						Username: "TestCreatePRSuccess_1assign_2users",
					},
				},
			},
			path:        basePathPR + "/create",
			description: "Create success 1 assign when 2 users in team",
			body: gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRSuccess_1assign_2users_1",
				PullRequestId:   "TestCreatePRSuccess_1assign_2users",
				PullRequestName: "TestCreatePRSuccess_1assign_2users",
			},
			expectedCode: http.StatusCreated,
			expectedBody: gen.PostPullRequestCreate201JSONResponse{
				Pr: &gen.PullRequest{
					PullRequestId:   "TestCreatePRSuccess_1assign_2users",
					PullRequestName: "TestCreatePRSuccess_1assign_2users",
					AuthorId:        "TestCreatePRSuccess_1assign_2users_1",
					Status:          gen.PullRequestStatusOPEN,
					AssignedReviewers: []string{
						"TestCreatePRSuccess_1assign_2users_2",
					},
				},
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestCreatePRSuccess_1assign_3users",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_1assign_3users_1",
						Username: "TestCreatePRSuccess_1assign_3users",
					},
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_1assign_3users_2",
						Username: "TestCreatePRSuccess_1assign_3users",
					},
					{
						IsActive: false,
						UserId:   "TestCreatePRSuccess_1assign_3users_3",
						Username: "TestCreatePRSuccess_1assign_3users",
					},
				},
			},
			path:        basePathPR + "/create",
			description: "Create success 1 assign when 3 users in team one is not active",
			body: gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRSuccess_1assign_3users_1",
				PullRequestId:   "TestCreatePRSuccess_1assign_3users",
				PullRequestName: "TestCreatePRSuccess_1assign_3users",
			},
			expectedCode: http.StatusCreated,
			expectedBody: gen.PostPullRequestCreate201JSONResponse{
				Pr: &gen.PullRequest{
					PullRequestId:   "TestCreatePRSuccess_1assign_3users",
					PullRequestName: "TestCreatePRSuccess_1assign_3users",
					AuthorId:        "TestCreatePRSuccess_1assign_3users_1",
					Status:          gen.PullRequestStatusOPEN,
					AssignedReviewers: []string{
						"TestCreatePRSuccess_1assign_3users_2",
					},
				},
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestCreatePRSuccess_0assign",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestCreatePRSuccess_0assign_1",
						Username: "TestCreatePRSuccess_0assign",
					},
				},
			},
			path:        basePathPR + "/create",
			description: "Create success 0 assign ",
			body: gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRSuccess_0assign_1",
				PullRequestId:   "TestCreatePRSuccess_0assign",
				PullRequestName: "TestCreatePRSuccess_0assign",
			},
			expectedCode: http.StatusCreated,
			expectedBody: gen.PostPullRequestCreate201JSONResponse{
				Pr: &gen.PullRequest{
					PullRequestId:     "TestCreatePRSuccess_0assign",
					PullRequestName:   "TestCreatePRSuccess_0assign",
					AuthorId:          "TestCreatePRSuccess_0assign_1",
					Status:            gen.PullRequestStatusOPEN,
					AssignedReviewers: []string{},
				},
			},
		}, {
			path:        basePathPR + "/create",
			description: "Create NotFound data",
			body: gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRNotFound",
				PullRequestId:   "TestCreatePRNotFound",
				PullRequestName: "TestCreatePRNotFound",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: gen.PostPullRequestCreate404JSONResponse{
				Error: GetError(cerr.NOT_FOUND).Error,
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestCreatePRExist",
				Members: []gen.TeamMember{
					{
						IsActive: false,
						UserId:   "TestCreatePRExist",
						Username: "TestCreatePRExist",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRExist",
				PullRequestId:   "TestCreatePRExist",
				PullRequestName: "TestCreatePRExist",
			},
			path:        basePathPR + "/create",
			description: "Create when PR exist",
			body: gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestCreatePRExist",
				PullRequestId:   "TestCreatePRExist",
				PullRequestName: "TestCreatePRExist",
			},
			expectedCode: http.StatusConflict,
			expectedBody: gen.PostPullRequestCreate409JSONResponse{
				Error: GetError(cerr.PR_EXISTS).Error,
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
			} else if resp.StatusCode == http.StatusCreated {
				var respBody gen.PostPullRequestCreate201JSONResponse
				expBody := test.expectedBody.(gen.PostPullRequestCreate201JSONResponse)
				if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				} else {
					if respBody.Pr.AuthorId != expBody.Pr.AuthorId &&
						respBody.Pr.PullRequestId != expBody.Pr.PullRequestId &&
						respBody.Pr.PullRequestName != expBody.Pr.PullRequestName &&
						respBody.Pr.Status != expBody.Pr.Status &&
						respBody.Pr.CreatedAt == nil {
						t.Fatalf("want %v got %v", expBody, respBody)
					}
					assert.ElementsMatch(t, respBody.Pr.AssignedReviewers, expBody.Pr.AssignedReviewers)
				}
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

// TestMerge test/pullRequest/merge
func TestMerge(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			teamForTest: &gen.Team{
				TeamName: "TestMergePRSuccess",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestMergePRSuccess",
						Username: "TestMergePRSuccess",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestMergePRSuccess",
				PullRequestId:   "TestMergePRSuccess",
				PullRequestName: "TestMergePRSuccess",
			},
			path:        basePathPR + "/merge",
			description: "Merge PR Success",
			body: gen.PostPullRequestMergeJSONRequestBody{
				PullRequestId: "TestMergePRSuccess",
			},
			expectedCode: http.StatusOK,
			expectedBody: gen.PostPullRequestMerge200JSONResponse{
				Pr: &gen.PullRequest{
					PullRequestId:     "TestMergePRSuccess",
					PullRequestName:   "TestMergePRSuccess",
					AuthorId:          "TestMergePRSuccess",
					Status:            gen.PullRequestStatusOPEN,
					AssignedReviewers: []string{},
				},
			},
		}, {
			path:        basePathPR + "/merge",
			description: "Merge PR NotFound",
			body: gen.PostPullRequestMergeJSONRequestBody{
				PullRequestId: "TestMergePRNotFound",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: gen.PostPullRequestMerge404JSONResponse{
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
			} else if resp.StatusCode == http.StatusOK {
				var respBody gen.PostPullRequestMerge200JSONResponse
				expBody := test.expectedBody.(gen.PostPullRequestMerge200JSONResponse)
				if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				} else {
					if respBody.Pr.AuthorId != expBody.Pr.AuthorId &&
						respBody.Pr.PullRequestId != expBody.Pr.PullRequestId &&
						respBody.Pr.PullRequestName != expBody.Pr.PullRequestName &&
						respBody.Pr.Status != expBody.Pr.Status &&
						respBody.Pr.CreatedAt == nil && respBody.Pr.MergedAt == nil {
						t.Fatalf("want %v got %v", expBody, respBody)
					}
					assert.ElementsMatch(t, respBody.Pr.AssignedReviewers, expBody.Pr.AssignedReviewers)
				}
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

// TestReassign test/pullRequest/reassign
func TestReassign(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tests := []TestData{
		{
			teamForTest: &gen.Team{
				TeamName: "TestReassignPRSuccess",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestReassignPRSuccess_1",
						Username: "TestReassignPRSuccess",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRSuccess_2",
						Username: "TestReassignPRSuccess",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRSuccess_3",
						Username: "TestReassignPRSuccess",
					},
					{
						IsActive: false,
						UserId:   "TestReassignPRSuccess_4",
						Username: "TestReassignPRSuccess",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestReassignPRSuccess_1",
				PullRequestId:   "TestReassignPRSuccess",
				PullRequestName: "TestReassignPRSuccess",
			},
			setActiveForTest: &gen.PostUsersSetIsActiveJSONBody{
				UserId:   "TestReassignPRSuccess_4",
				IsActive: true,
			},
			path:        basePathPR + "/reassign",
			description: "Reassign PR Success",
			body: gen.PostPullRequestReassignJSONBody{
				PullRequestId: "TestReassignPRSuccess",
				OldUserId:     "TestReassignPRSuccess_2",
			},
			expectedCode: http.StatusOK,
			expectedBody: gen.PostPullRequestReassign200JSONResponse{
				Pr: gen.PullRequest{
					PullRequestId:   "TestReassignPRSuccess",
					PullRequestName: "TestReassignPRSuccess",
					AuthorId:        "TestReassignPRSuccess_1",
					Status:          gen.PullRequestStatusMERGED,
					AssignedReviewers: []string{
						"TestReassignPRSuccess_3",
						"TestReassignPRSuccess_4",
					},
				},
				ReplacedBy: "TestReassignPRSuccess_4",
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: " TestReassignPRNotFoundUser",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestReassignPRNotFoundUser_1",
						Username: "TestReassignPRNotFoundUser",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNotFoundUser_2",
						Username: "TestReassignPRNotFoundUser",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNotFoundUser_3",
						Username: "TestReassignPRNotFoundUser",
					},
					{
						IsActive: false,
						UserId:   "TestReassignPRNotFoundUser_4",
						Username: "TestReassignPRNotFoundUser",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestReassignPRNotFoundUser_1",
				PullRequestId:   "TestReassignPRNotFoundUser",
				PullRequestName: "TestReassignPRNotFoundUser",
			},
			setActiveForTest: &gen.PostUsersSetIsActiveJSONBody{
				UserId:   "TestReassignPRNotFoundUser_4",
				IsActive: true,
			},
			path:        basePathPR + "/reassign",
			description: "Reassign PR Not Found User",
			body: gen.PostPullRequestReassignJSONBody{
				PullRequestId: "TestReassignPRNotFoundUser",
				OldUserId:     "TestReassignPRNotFoundUser_5",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: gen.PostPullRequestReassign404JSONResponse{
				Error: GetError(cerr.NOT_FOUND).Error,
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestReassignPRNotFoundPR",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestReassignPRNotFoundPR_1",
						Username: "TestReassignPRNotFoundPR",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNotFoundPR_2",
						Username: "TestReassignPRNotFoundPR",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNotFoundPR_3",
						Username: "TestReassignPRNotFoundPR",
					},
					{
						IsActive: false,
						UserId:   "TestReassignPRNotFoundPR_4",
						Username: "TestReassignPRNotFoundPR",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestReassignPRNotFoundPR_1",
				PullRequestId:   "TestReassignPRNotFoundPR",
				PullRequestName: "TestReassignPRNotFoundPR",
			},
			setActiveForTest: &gen.PostUsersSetIsActiveJSONBody{
				UserId:   "TestReassignPRNotFoundPR_4",
				IsActive: true,
			},
			path:        basePathPR + "/reassign",
			description: "Reassign PR Not Found PR",
			body: gen.PostPullRequestReassignJSONBody{
				PullRequestId: "TestReassignPRNotFoundPR_2",
				OldUserId:     "TestReassignPRNotFoundPR_2",
			},
			expectedCode: http.StatusNotFound,
			expectedBody: gen.PostPullRequestReassign404JSONResponse{
				Error: GetError(cerr.NOT_FOUND).Error,
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestReassignPRMergedPR",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestReassignPRMergedPR_1",
						Username: "TestReassignPRMergedPR",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRMergedPR_2",
						Username: "TestReassignPRMergedPR",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRMergedPR_3",
						Username: "TestReassignPRMergedPR",
					},
					{
						IsActive: false,
						UserId:   "TestReassignPRMergedPR_4",
						Username: "TestReassignPRMergedPR",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestReassignPRMergedPR_1",
				PullRequestId:   "TestReassignPRMergedPR",
				PullRequestName: "TestReassignPRMergedPR",
			},
			mergeForTest: &gen.PostPullRequestMergeJSONBody{
				PullRequestId: "TestReassignPRMergedPR",
			},
			path:        basePathPR + "/reassign",
			description: "Reassign PR Merged PR",
			body: gen.PostPullRequestReassignJSONBody{
				PullRequestId: "TestReassignPRMergedPR",
				OldUserId:     "TestReassignPRMergedPR_2",
			},
			expectedCode: http.StatusConflict,
			expectedBody: gen.PostPullRequestReassign409JSONResponse{
				Error: GetError(cerr.PR_MERGED).Error,
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestReassignPRNotAssigned",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestReassignPRNotAssigned_1",
						Username: "TestReassignPRNotAssigned",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNotAssigned_2",
						Username: "TestReassignPRNotAssigned",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNotAssigned_3",
						Username: "TestReassignPRNotAssigned",
					},
					{
						IsActive: false,
						UserId:   "TestReassignPRNotAssigned_4",
						Username: "TestReassignPRNotAssigned",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestReassignPRNotAssigned_1",
				PullRequestId:   "TestReassignPRNotAssigned",
				PullRequestName: "TestReassignPRNotAssigned",
			},
			setActiveForTest: &gen.PostUsersSetIsActiveJSONBody{
				UserId:   "TestReassignPRNotAssigned_4",
				IsActive: true,
			},
			path:        basePathPR + "/reassign",
			description: "Reassign PR NotAssigned user",
			body: gen.PostPullRequestReassignJSONBody{
				PullRequestId: "TestReassignPRNotAssigned",
				OldUserId:     "TestReassignPRNotAssigned_4",
			},
			expectedCode: http.StatusConflict,
			expectedBody: gen.PostPullRequestReassign409JSONResponse{
				Error: GetError(cerr.NOT_ASSIGNED).Error,
			},
		}, {
			teamForTest: &gen.Team{
				TeamName: "TestReassignPRNoCandidate",
				Members: []gen.TeamMember{
					{
						IsActive: true,
						UserId:   "TestReassignPRNoCandidate_1",
						Username: "TestReassignPRNoCandidate",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNoCandidate_2",
						Username: "TestReassignPRNoCandidate",
					},
					{
						IsActive: true,
						UserId:   "TestReassignPRNoCandidate_3",
						Username: "TestReassignPRNoCandidate",
					},
					{
						IsActive: false,
						UserId:   "TestReassignPRNoCandidate_4",
						Username: "TestReassignPRNoCandidate",
					},
				},
			},
			prForTest: &gen.PostPullRequestCreateJSONBody{
				AuthorId:        "TestReassignPRNoCandidate_1",
				PullRequestId:   "TestReassignPRNoCandidate",
				PullRequestName: "TestReassignPRNoCandidate",
			},
			path:        basePathPR + "/reassign",
			description: "Reassign PR NoCandidate user",
			body: gen.PostPullRequestReassignJSONBody{
				PullRequestId: "TestReassignPRNoCandidate",
				OldUserId:     "TestReassignPRNoCandidate_2",
			},
			expectedCode: http.StatusConflict,
			expectedBody: gen.PostPullRequestReassign409JSONResponse{
				Error: GetError(cerr.NO_CANDIDATE).Error,
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

			if test.setActiveForTest != nil {
				err := SetIsActiveForTest(test.setActiveForTest)
				if err != nil {
					t.Fatalf("error: %s", err)
				}
			}

			if test.mergeForTest != nil {
				err := MergePRForTest(test.mergeForTest)
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
			} else if resp.StatusCode == http.StatusOK {
				var respBody gen.PostPullRequestReassign200JSONResponse
				expBody := test.expectedBody.(gen.PostPullRequestReassign200JSONResponse)
				if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				} else {
					if respBody.Pr.AuthorId != expBody.Pr.AuthorId &&
						respBody.Pr.PullRequestId != expBody.Pr.PullRequestId &&
						respBody.Pr.PullRequestName != expBody.Pr.PullRequestName &&
						respBody.Pr.Status != expBody.Pr.Status &&
						respBody.Pr.CreatedAt == nil && respBody.ReplacedBy != expBody.ReplacedBy {
						t.Fatalf("want %v got %v", expBody, respBody)
					}
					assert.ElementsMatch(t, respBody.Pr.AssignedReviewers, expBody.Pr.AssignedReviewers)
				}
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
