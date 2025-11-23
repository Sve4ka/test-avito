package stat

import (
	"context"

	"avito/internal/entity"
	"avito/internal/log"
	"avito/internal/repo"
	"avito/internal/service"
)

type Serv struct {
	Repo repo.Stat
}

func InitStatServ(repo repo.Stat) service.Stat {
	return Serv{Repo: repo}
}

func (s Serv) User(ctx context.Context, userID string) (*entity.UserStat, error) {
	user, err := s.Repo.User(ctx, userID)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	return user, nil
}

func (s Serv) Team(ctx context.Context, teamName string) (*entity.TeamStat, error) {
	team, err := s.Repo.Team(ctx, teamName)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	return team, nil
}
