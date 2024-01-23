package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Webforum is up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := app.DB.AllThreads()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, threads)
}

func (app *application) GetThread(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	threadID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	thread, err := app.DB.SingleThread(threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, thread)
}

func (app *application) GetComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	threadID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	thread, err := app.DB.GetCommentsByThreadID(threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, thread)
}
