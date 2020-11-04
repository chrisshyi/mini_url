package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/{shortURL}", app.redirectShortURL).Methods("GET")
	router.HandleFunc("/", app.createNewShortURL).Methods("POST")
	return router
}
