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
	mux.Get("/comment/{id}/replies", app.GetReplies)
	mux.Post("/authenticate", app.authenticate)

	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	mux.Put("/signup", app.registerUser)

	mux.Route("/create", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Put("/thread", app.InsertThread)
		mux.Put("/comment", app.InsertComment)
		mux.Put("/reply", app.InsertReply)
	})

	return mux
}
