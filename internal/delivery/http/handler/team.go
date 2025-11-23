package handler

import (
	"context"
	"net/http"

	"avito/internal/cerr"
	"avito/internal/entity"
	"avito/internal/gen"
	"avito/internal/service"
)

type Team struct {
	service service.Team
}

func InitTeamHandler(service service.Team) *Team {
	return &Team{
		service: service,
	}
}

func (r *Team) PostTeamAdd(ctx context.Context, request gen.PostTeamAddRequestObject) (gen.PostTeamAddResponseObject, error) {
	createTeam := entity.Team{
		TeamName: request.Body.TeamName,
		Members:  make([]entity.TeamMember, len(request.Body.Members)),
	}
	for i, member := range request.Body.Members {
		createTeam.Members[i] = entity.TeamMember{
			IsActive: member.IsActive,
			UserId:   member.UserId,
			Username: member.Username,
		}
	}

	err := r.service.Create(ctx, &createTeam)
	if err != nil {
		code, message := cerr.HandleErrs(err)
		if code == http.StatusBadRequest {
			return gen.PostTeamAdd400JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	return gen.PostTeamAdd201JSONResponse{Team: request.Body}, nil
}

func (r *Team) GetTeamGet(ctx context.Context, request gen.GetTeamGetRequestObject) (gen.GetTeamGetResponseObject, error) {
	team, err := r.service.Get(ctx, request.Params.TeamName)
	if err != nil {
		code, message := cerr.HandleErrs(err)

		if code == http.StatusNotFound {
			return gen.GetTeamGet404JSONResponse(message), nil
		}

		return nil, cerr.ErrServerTime
	}

	genTeam := gen.Team{
		TeamName: team.TeamName,
		Members:  make([]gen.TeamMember, len(team.Members)),
	}
	for i, member := range team.Members {
		genTeam.Members[i] = gen.TeamMember{
			IsActive: member.IsActive,
			UserId:   member.UserId,
			Username: member.Username,
		}
	}

	return gen.GetTeamGet200JSONResponse(genTeam), nil
}
