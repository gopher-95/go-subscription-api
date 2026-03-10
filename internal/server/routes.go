package server

import (
	"github.com/go-chi/chi"
	"github.com/gopher-95/go-subscription-api/internal/handlers"
)

func Router() *chi.Mux {
	var router *chi.Mux

	router.Get("get/{id}", handlers.GetHandler)

	return router
}
