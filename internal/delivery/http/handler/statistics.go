package handler

import (
	"context"
	"net/http"

	"avito/internal/cerr"
	"avito/internal/gen"
	"avito/internal/service"
)

type Stat struct {
	service service.Stat
}

func InitStatHandler(service service.Stat) *Stat {
	return &Stat{
		service: service,
	}
}

func (s Stat) GetStatisticsTeam(ctx context.Context, request gen.GetStatisticsTeamRequestObject) (gen.GetStatisticsTeamResponseObject, error) {
	team, err := s.service.Team(ctx, request.Params.TeamName)
	if err != nil {
		code, message := cerr.HandleErrs(err)

		if code == http.StatusNotFound {
			return gen.GetStatisticsTeam404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	var genUsers []gen.UserStat
	for _, user := range team.UsersStat {
		genUsers = append(genUsers, gen.UserStat{
			AvgDuration: user.AvgDuration,
			CountPr:     user.CountPr,
			IsActive:    user.IsActive,
			UserId:      user.UserId,
		})
	}

	return gen.GetStatisticsTeam200JSONResponse{
		TeamName:    team.TeamName,
		UsersStat:   genUsers,
		AvgDuration: team.AvgDuration,
	}, nil
}

func (s Stat) GetStatisticsUser(ctx context.Context, request gen.GetStatisticsUserRequestObject) (gen.GetStatisticsUserResponseObject, error) {
	user, err := s.service.User(ctx, request.Params.UserId)
	if err != nil {
		code, message := cerr.HandleErrs(err)
		if code == http.StatusNotFound {
			return gen.GetStatisticsUser404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	return gen.GetStatisticsUser200JSONResponse{
		AvgDuration: user.AvgDuration,
		CountPr:     user.CountPr,
		IsActive:    user.IsActive,
		UserId:      user.UserId,
	}, nil
}
