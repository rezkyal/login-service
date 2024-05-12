package handler

import (
	"github.com/SawitProRecruitment/UserService/usecase"
)

type Server struct {
	Usecase usecase.UsecaseInterface
}

type NewServerOptions struct {
	Usecase usecase.UsecaseInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Usecase: opts.Usecase,
	}
}
