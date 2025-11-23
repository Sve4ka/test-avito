package pullRequest

import (
	"context"

	"avito/internal/entity"
	"avito/internal/log"
	"avito/internal/repo"
	"avito/internal/service"
)

type Serv struct {
	Repo repo.PullRequest
}

func InitPullRequestServ(repo repo.PullRequest) service.PullRequest {
	return Serv{Repo: repo}
}

func (s Serv) Create(ctx context.Context, pullRequestCreate *entity.PullRequestCreate) (*entity.PullRequest, error) {
	pullRequest, err := s.Repo.Create(ctx, pullRequestCreate)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	return pullRequest, nil
}

func (s Serv) Merge(ctx context.Context, pullRequestID string) (*entity.PullRequest, error) {
	pullRequest, err := s.Repo.Merge(ctx, pullRequestID)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	return pullRequest, nil
}

func (s Serv) Reassign(ctx context.Context, pullRequestID string, oldUserID string) (*entity.PullRequest, string, error) {
	pullRequest, newReviewer, err := s.Repo.Reassign(ctx, pullRequestID, oldUserID)
	if err != nil {
		log.Log.Error(err)

		return nil, "", err
	}

	return pullRequest, newReviewer, nil
}
