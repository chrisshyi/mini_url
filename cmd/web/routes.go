package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {

	mux := pat.New()

	mux.Get("/:shortened_url", http.NotFoundHandler())
	mux.Post("/", http.NotFoundHandler())

	return mux
}
