package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewURL represents a new URL to be shortened
type NewURL struct {
	URL string `json:"URL"`
}

// ShortenedURL represents a shortened URL
type ShortenedURL struct {
	URL string `json:"URL"`
}

func (app *application) redirectShortURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shortURL := vars["shortURL"]
	app.logInfo(fmt.Sprintf("Redirecting short URL %s...", shortURL))
	ID, err := shortURLToID(shortURL)
	if err != nil {
		app.logErr(err.Error())
		if errors.Is(err, ErrInvalidShortURL) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	miniURL, err := app.miniURLModel.GetByID(ID)
	if err != nil {
		app.logErr(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if miniURL == nil {
		http.Error(w, "short URL does not exist", http.StatusNotFound)
		return
	}
	app.logInfo(fmt.Sprintf("Redirecting short URL %s to %s", shortURL, miniURL.URL))
	http.Redirect(w, r, miniURL.URL, http.StatusSeeOther)
}

func (app *application) createNewShortURL(w http.ResponseWriter, r *http.Request) {
	var newURL NewURL
	err := decodeJSONBody(w, r, &newURL)

	if err != nil {
		app.logErr(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	miniURL, err := app.miniURLModel.GetByURL(newURL.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.logInfo(fmt.Sprintf("No rows found for URL %s", newURL.URL))
			app.logInfo(fmt.Sprintf("Creating new short URL for URL %s", newURL.URL))
		} else {
			app.logErr(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if miniURL != nil {
		http.Error(w, "URL already exists", http.StatusBadRequest)
		return
	}
	newMiniURLID, err := app.miniURLModel.Insert(newURL.URL)
	if err != nil {
		app.logErr(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	shortenedURL, err := IDToShortURL(newMiniURLID)
	if err != nil {
		app.logErr(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortenedURL{
		URL: shortenedURL,
	})
	return
}
