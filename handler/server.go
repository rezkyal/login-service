package handler

import (
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
)

type Server struct {
	Repository repository.RepositoryInterface
	Usecase    usecase.UsecaseInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Usecase    usecase.UsecaseInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Usecase:    opts.Usecase,
	}
}
