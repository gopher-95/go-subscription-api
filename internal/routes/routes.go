package routes

import (
	"github.com/go-chi/chi"
	"github.com/gopher-95/go-subscription-api/internal/handlers"
)

func Route(r *chi.Mux, h *handlers.SubscriptionHandler) {
	r.Get("/get/{id}", h.GetByID)
}
