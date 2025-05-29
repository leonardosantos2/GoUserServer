package routes

import (
	"github.com/go-chi/chi/v5"

	"flippOneApi/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)

	router.Get("/users/{id}", app.UserHandler.GetUser)

	return router
}
