package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/leonardosantos2/GoUserServer/internal/app"
	"github.com/leonardosantos2/GoUserServer/internal/middlewares"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(corsMiddleware.Handler)

	router.Get("/health", app.HealthCheck)

	router.With(middleware.EnsureValidToken()).Get("/user", app.UserHandler.GetUser)
	router.With(middleware.EnsureValidToken()).Get("/user/roles", app.UserHandler.GetUserRolesAndAppMetadata)

	return router
}
