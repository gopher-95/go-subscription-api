package server

import (
	"github.com/go-chi/chi"
	"github.com/gopher-95/go-subscription-api/internal/handlers"
)

func Router(subscriptionHandler *handlers.SubscriptionHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api/v1/subscriptions", func(r chi.Router) {
		r.Post("/", subscriptionHandler.Create)
		r.Get("/{id}", subscriptionHandler.Get)
		r.Put("/{id}", subscriptionHandler.Update)
		r.Delete("/{id}", subscriptionHandler.Delete)
		r.Get("/", subscriptionHandler.GetAll)
	})

	return router
}
