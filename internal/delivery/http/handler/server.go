package handler

import (
	"avito/internal/gen"
)

var _ gen.StrictServerInterface = (*Server)(nil)

type Server struct {
	*Team
	*PullRequest
	*User
	*Stat
}

func NewServer(
	userHandler *User,
	prHandler *PullRequest,
	teamHandler *Team,
	statHandler *Stat,
) *Server {
	return &Server{
		User:        userHandler,
		PullRequest: prHandler,
		Team:        teamHandler,
		Stat:        statHandler,
	}
}
