package http

import (
	"avito/internal/delivery/http/handler"
	"avito/internal/gen"
	"avito/internal/postgres"
	PRRepo "avito/internal/repo/pullRequest"
	statRepo "avito/internal/repo/stat"
	teamRepo "avito/internal/repo/team"
	userRepo "avito/internal/repo/user"
	PRServ "avito/internal/service/pullRequest"
	statServ "avito/internal/service/stat"
	teamServ "avito/internal/service/team"
	userServ "avito/internal/service/user"
)

func InitServer(db *postgres.Pg) gen.ServerInterface {
	repoUser := userRepo.InitUserRepo(db)
	servUser := userServ.InitUserServ(repoUser)
	handlerUser := handler.InitUserHandler(servUser)

	repoTeam := teamRepo.InitTeamRepo(db)
	servTeam := teamServ.InitTeamServ(repoTeam)
	handlerTeam := handler.InitTeamHandler(servTeam)

	repoPR := PRRepo.InitPullRequestRepo(db)
	servPR := PRServ.InitPullRequestServ(repoPR)
	handlerPR := handler.InitPullRequestHandler(servPR)

	repoStat := statRepo.InitStatRepo(db)
	servStat := statServ.InitStatServ(repoStat)
	handlerStat := handler.InitStatHandler(servStat)

	server := handler.NewServer(handlerUser, handlerPR, handlerTeam, handlerStat)

	strictHandler := gen.NewStrictHandler(server, nil)

	return strictHandler
}
