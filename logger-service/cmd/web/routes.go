package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Mount("/", app.webRouter())
	mux.Mount("/api", app.apiRouter())

	return mux
}

// apiRouter is for api routes (no session load)
func (app *Config) apiRouter() http.Handler {
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

	mux.Post("/log", app.WriteLog)
	return mux
}

// webRouter is for web routes
func (app *Config) webRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.SessionLoad)

	mux.Get("/", app.LoginPage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.LoginPagePost)
	mux.Get("/logout", app.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Get("/dashboard", app.Dashboard)
		mux.Get("/log-entry/{id}", app.DisplayOne)
		mux.Get("/delete-all", app.DeleteAll)
		mux.Get("/update/{id}", app.UpdateTimeStamp)
	})

	return mux
}
