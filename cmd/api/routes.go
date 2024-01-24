package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/threads", app.AllThreads)
	mux.Get("/thread/{id}", app.GetThread)
	mux.Get("/thread/{id}/comments", app.GetComments)
	mux.Post("/authenticate", app.authenticate)

	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	mux.Route("/create", func(mux chi.Router) {
		mux.Use(app.authRequired)
	})

	return mux
}
