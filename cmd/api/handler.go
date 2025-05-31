package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/voidt01/go-url-shortener/internal/models"
)

// Request and Response for creating shortenURL
type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}

func (a *App) Shortening(w http.ResponseWriter, r *http.Request) {
	urlRequest := &shortenRequest{}
	// urlResponse := &shortenResponse{}

	// Decode POST Request body (JSON) to Go
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() 

	err_decode := dec.Decode(urlRequest)
	if err_decode != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	// url nil value check
	if urlRequest.OriginalURL == "" {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	// url validation
	if !isValid(urlRequest.OriginalURL) {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%+v\n", urlRequest)

	// Generate Short URL & Store urls in Database
	// shortCode, err_model := a.urls.Insert(urlRequest.OriginalURL)
	// if err_model != nil {
	// 	a.serveError(w, err_model)
	// }

	// // creating post response struct
	// urlResponse.OriginalURL = urlRequest.OriginalURL
	// urlResponse.ShortenURL = a.configApp.Server.baseURL + a.configApp.Server.port + "/" + shortCode

	// // encoding response struct (G0) to JSON
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// err_encode := json.NewEncoder(w).Encode(urlResponse)

	// if err_encode != nil {
	// 	a.errorLog.Printf("Error encoding response: %v", err_encode)
	// }

}

func (a *App) Redirecting(w http.ResponseWriter, r *http.Request) {
	short_code := r.URL.Path[1:]

	value, err := a.urls.Get(short_code)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serveError(w, err)
		}
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}

func isValid(ori_url string) bool {
	u, err := url.Parse(ori_url)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http")
}
