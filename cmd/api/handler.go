package main

import (
	"errors"
	"net/http"

	"github.com/voidt01/go-url-shortener/internal/models"
)

type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}

func (a *App) Shortening(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	var clientErr *clientError

	// decoding JSON to Go obj
	err := a.decodeJSON(w, r, &req)
	if err != nil {
		if errors.As(err, &clientErr) {
			a.ErrorResponse(w, clientErr.msg, clientErr.status)
		} else {
			a.errorLog.Print(err.Error())
			a.ErrorResponse(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError)
		}
		return
	}

	// validation for the original url
	url, err := a.urlValidation(req.OriginalURL)
	if err != nil {
		if errors.As(err, &clientErr) {
			a.ErrorResponse(w, clientErr.msg, clientErr.status)
		} else {
			a.errorLog.Print(err.Error())
			a.ErrorResponse(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError)
		}
		return
	}

	// Generate Short URL & Store urls in Database
	shortCode, err_model := a.urls.Insert(url)
	if err_model != nil {
		a.errorLog.Print(err_model.Error())
		a.ErrorResponse(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError)
		return
	}

	// creating post response struct
	resp := &shortenResponse{
		OriginalURL: url,
		ShortenURL:  a.builderShortenURL(shortCode),
	}

	// encoding response struct (G0) to JSON
	err_encode := a.encodeJSON(w, &resp, http.StatusCreated)
	if err_encode != nil {
		a.errorLog.Print(err_encode.Error())
		a.ErrorResponse(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError)
	}

}

func (a *App) Redirecting(w http.ResponseWriter, r *http.Request) {
	short_code := r.URL.Path[1:]

	value, err := a.urls.Get(short_code)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			a.errorLog.Print(err.Error())
			http.Error(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}
