package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect to our API service
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// a heartbeat route, to ensure things are up
	mux.Use(middleware.Heartbeat("/ping"))

	// this route is just to ensure things work, and is never
	// used after that
	mux.Get("/", app.Broker)

	mux.Post("/", app.Broker)

	// grpc route
	mux.Post("/log-grpc", app.LogViaGRPC)

	// a route for everything
	mux.Post("/handle", app.HandleSubmission)

	return mux
}
