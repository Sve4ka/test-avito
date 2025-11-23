package service

import (
	"context"

	"avito/internal/entity"
)

type Team interface {
	Create(ctx context.Context, team *entity.Team) error
	Get(ctx context.Context, teamName string) (*entity.Team, error)
}

type User interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*entity.User, error)
	GetReview(ctx context.Context, userID string) ([]entity.PullRequestShort, error)
}

type PullRequest interface {
	Create(ctx context.Context, PullRequestCreate *entity.PullRequestCreate) (*entity.PullRequest, error)
	Merge(ctx context.Context, PullRequestID string) (*entity.PullRequest, error)
	Reassign(ctx context.Context, PullRequestID string, oldUserID string) (*entity.PullRequest, string, error)
}

type Stat interface {
	User(ctx context.Context, userID string) (*entity.UserStat, error)
	Team(ctx context.Context, teamName string) (*entity.TeamStat, error)
}
