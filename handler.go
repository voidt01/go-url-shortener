package main

import (
	"errors"
	"net/http"
)

type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}

type urlHandler struct {
	service *urlService
}

func (uh *urlHandler) Shortening(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	// decoding JSON to Go obj
	err := decodeJSON(w, r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validation for the original url
	url, err := urlValidation(req.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Generate Short URL & Store urls in Database
	shortCode, err_model := app.urls.Insert(url)
	if err_model != nil {
		http.Error(w, err_model.Error(), http.StatusBadRequest)
		return
	}

	// creating post response struct
	resp := &shortenResponse{
		OriginalURL: url,
		ShortenURL:  app.builderShortenURL(shortCode),
	}

	// encoding response struct (G0) to JSON
	err_encode := encodeJSON(w, &resp, http.StatusCreated)
	if err_encode != nil {
		a.errorLog.Print(err_encode.Error())
		a.ErrorResponseJSON(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError)
	}

}

func (uh *urlHandler) Redirecting(w http.ResponseWriter, r *http.Request) {
	short_code := r.URL.Path[1:]

	original_url, err := uh.service.ResolveUrl(short_code)
	if err != nil {
		switch {
		case errors.Is(err, ErrUrlNotFound):
			http.Error(w, "This link doesn't exist", http.StatusNotFound) 
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, original_url, http.StatusFound)
}
