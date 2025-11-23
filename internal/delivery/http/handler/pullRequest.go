package handler

import (
	"context"
	"net/http"

	"avito/internal/cerr"
	"avito/internal/entity"
	"avito/internal/gen"
	"avito/internal/service"
)

type PullRequest struct {
	service service.PullRequest
}

func InitPullRequestHandler(service service.PullRequest) *PullRequest {
	return &PullRequest{
		service: service,
	}
}

func (r *PullRequest) PostPullRequestCreate(ctx context.Context, request gen.PostPullRequestCreateRequestObject) (gen.PostPullRequestCreateResponseObject, error) {
	PullRequestCreate := entity.PullRequestCreate{
		AuthorId:        request.Body.AuthorId,
		PullRequestId:   request.Body.PullRequestId,
		PullRequestName: request.Body.PullRequestName,
	}

	pullRequest, err := r.service.Create(ctx, &PullRequestCreate)
	if err != nil {
		code, message := cerr.HandleErrs(err)

		if code == http.StatusConflict {
			return gen.PostPullRequestCreate409JSONResponse(message), nil
		}

		if code == http.StatusNotFound {
			return gen.PostPullRequestCreate404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	return gen.PostPullRequestCreate201JSONResponse{
		Pr: &gen.PullRequest{
			AssignedReviewers: pullRequest.AssignedReviewers,
			AuthorId:          pullRequest.AuthorId,
			CreatedAt:         pullRequest.CreatedAt,
			MergedAt:          pullRequest.MergedAt,
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			Status:            gen.PullRequestStatus(pullRequest.Status),
		},
	}, nil
}

func (r *PullRequest) PostPullRequestMerge(ctx context.Context, request gen.PostPullRequestMergeRequestObject) (gen.PostPullRequestMergeResponseObject, error) {
	pullRequest, err := r.service.Merge(ctx, request.Body.PullRequestId)
	if err != nil {
		code, message := cerr.HandleErrs(err)
		if code == http.StatusNotFound {
			return gen.PostPullRequestMerge404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	return gen.PostPullRequestMerge200JSONResponse{
		Pr: &gen.PullRequest{
			AssignedReviewers: pullRequest.AssignedReviewers,
			AuthorId:          pullRequest.AuthorId,
			CreatedAt:         pullRequest.CreatedAt,
			MergedAt:          pullRequest.MergedAt,
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			Status:            gen.PullRequestStatus(pullRequest.Status),
		},
	}, nil
}

func (r *PullRequest) PostPullRequestReassign(ctx context.Context, request gen.PostPullRequestReassignRequestObject) (gen.PostPullRequestReassignResponseObject, error) {
	pullRequest, newReviewer, err := r.service.Reassign(ctx, request.Body.PullRequestId, request.Body.OldUserId)
	if err != nil {
		code, message := cerr.HandleErrs(err)
		if code == http.StatusNotFound {
			return gen.PostPullRequestReassign404JSONResponse(message), nil
		}

		if code == http.StatusConflict {
			return gen.PostPullRequestReassign409JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	return gen.PostPullRequestReassign200JSONResponse{
		Pr: gen.PullRequest{
			AssignedReviewers: pullRequest.AssignedReviewers,
			AuthorId:          pullRequest.AuthorId,
			CreatedAt:         pullRequest.CreatedAt,
			MergedAt:          pullRequest.MergedAt,
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			Status:            gen.PullRequestStatus(pullRequest.Status),
		},
		ReplacedBy: newReviewer,
	}, nil
}
