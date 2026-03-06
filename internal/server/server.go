package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	router *chi.Mux
	port   string
}

func NewServer(r *chi.Mux, port string) *Server {
	return &Server{
		router: r,
		port:   port,
	}
}

func (s *Server) Run() error {
	log.Printf("Сервер запущен на порту: %s", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
