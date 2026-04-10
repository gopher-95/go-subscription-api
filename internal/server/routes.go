package server

import (
	"github.com/go-chi/chi"
	"github.com/gopher-95/go-subscription-api/internal/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(subscriptionHandler *handlers.SubscriptionHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL к swagger.json
	))

	router.Route("/api/v1/subscriptions", func(r chi.Router) {
		r.Post("/", subscriptionHandler.Create)
		r.Get("/{id}", subscriptionHandler.Get)
		r.Put("/{id}", subscriptionHandler.Update)
		r.Delete("/{id}", subscriptionHandler.Delete)
		r.Get("/total-cost", subscriptionHandler.GetTotalCost)
		r.Get("/", subscriptionHandler.GetAll)
	})

	return router
}
