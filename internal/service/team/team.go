package team

import (
	"context"

	"avito/internal/cerr"
	"avito/internal/entity"
	"avito/internal/log"
	"avito/internal/repo"
	"avito/internal/service"
)

type Serv struct {
	Repo repo.Team
}

func InitTeamServ(repo repo.Team) service.Team {
	return Serv{Repo: repo}
}

func (s Serv) Create(ctx context.Context, team *entity.Team) error {
	isFreeName, err := s.Repo.CheckTeamName(ctx, team.TeamName)
	if err != nil {
		log.Log.Error(err)

		return err
	}

	if !isFreeName {
		err = cerr.CustomError{Err: err, ErrType: cerr.TEAM_EXISTS}
		log.Log.Error(err)

		return err
	}

	err = s.Repo.Create(ctx, team)
	if err != nil {
		log.Log.Error(err)

		return err
	}

	return nil
}

func (s Serv) Get(ctx context.Context, teamName string) (*entity.Team, error) {
	team, err := s.Repo.Get(ctx, teamName)
	if err != nil {
		log.Log.Error(err)

		return nil, err
	}

	if len(team.Members) == 0 {
		err = cerr.CustomError{Err: err, ErrType: cerr.NOT_FOUND}

		log.Log.Error(err)

		return nil, err
	}

	return team, nil
}
