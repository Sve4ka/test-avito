package handler

import (
	"context"
	"net/http"

	"avito/internal/cerr"
	"avito/internal/gen"
	"avito/internal/service"
)

type User struct {
	service service.User
}

func InitUserHandler(service service.User) *User {
	return &User{
		service: service,
	}
}

func (r *User) GetUsersGetReview(ctx context.Context, request gen.GetUsersGetReviewRequestObject) (gen.GetUsersGetReviewResponseObject, error) {
	PullRequests, err := r.service.GetReview(ctx, request.Params.UserId)
	if err != nil {
		code, message := cerr.HandleErrs(err)
		if code == http.StatusNotFound {
			return gen.GetUsersGetReview404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	genPullRequests := make([]gen.PullRequestShort, len(PullRequests))

	for i, pr := range PullRequests {
		genPullRequests[i] = gen.PullRequestShort{
			AuthorId:        pr.AuthorId,
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			Status:          gen.PullRequestShortStatus(pr.Status),
		}
	}

	return gen.GetUsersGetReview200JSONResponse{
		UserId:       request.Params.UserId,
		PullRequests: genPullRequests,
	}, nil
}

func (r *User) PostUsersSetIsActive(ctx context.Context, request gen.PostUsersSetIsActiveRequestObject) (gen.PostUsersSetIsActiveResponseObject, error) {
	user, err := r.service.SetIsActive(ctx, request.Body.UserId, request.Body.IsActive)
	if err != nil {
		code, message := cerr.HandleErrs(err)
		if code == http.StatusNotFound {
			return gen.PostUsersSetIsActive404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	return gen.PostUsersSetIsActive200JSONResponse{
		User: &gen.User{
			UserId:   user.UserId,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}, nil
}
