package server

import (
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, hanlder http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: hanlder,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("сервер запущен на порту %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
