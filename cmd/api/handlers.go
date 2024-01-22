package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Hello world from  ", app.Domain)
}

func (app *application) AllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := app.DB.AllThreads()
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := json.Marshal(threads)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
