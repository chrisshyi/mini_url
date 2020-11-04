package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) redirectShortURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shortURL := vars["shortURL"]
	fmt.Fprintf(w, "Short url: %s\n", shortURL)
	return
}

func (app *application) createNewShortURL(w http.ResponseWriter, r *http.Request) {
	return
}
