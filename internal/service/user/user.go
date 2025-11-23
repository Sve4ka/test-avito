package user

import (
	"context"

	"avito/internal/entity"
	"avito/internal/log"
	"avito/internal/repo"
	"avito/internal/service"
)

type Serv struct {
	Repo repo.User
}

func InitUserServ(repo repo.User) service.User {
	return Serv{Repo: repo}
}

func (s Serv) SetIsActive(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	user, err := s.Repo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	return user, nil
}

func (s Serv) GetReview(ctx context.Context, userID string) ([]entity.PullRequestShort, error) {
	reviews, err := s.Repo.GetReview(ctx, userID)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	return reviews, nil
}
