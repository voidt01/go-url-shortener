package main

import (
	"errors"
	// "fmt"
	"net/http"

	"github.com/voidt01/go-url-shortener/internal/models"
)

type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

// type shortenResponse struct {
// 	OriginalURL string `json:"original_url"`
// 	ShortenURL  string `json:"shorten_url"`
// }

func (a *App) Shortening(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	var clientErr *clientError

	err := a.decodeJSON(w, r, &req)
	if err != nil {
		if errors.As(err, &clientErr){
			http.Error(w, clientErr.msg, clientErr.status)
		} else {
			a.errorLog.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = a.isValid(req.OriginalURL)
	if err != nil {
		if errors.As(err, &clientErr){
			http.Error(w, clientErr.msg, clientErr.status)
		} else {
			a.errorLog.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

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
			a.clientError(w, http.StatusNotFound)
		} else {
			a.serveError(w, err)
		}
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}