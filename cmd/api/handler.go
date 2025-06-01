package main

import (
	"encoding/json"
	"errors"
	// "fmt"
	"io"
	"net/http"
	"net/url"

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
	urlRequest := shortenRequestHelper(a, w, r.Body)
	if urlRequest == nil {
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
			a.notFound(w)
		} else {
			a.serveError(w, err)
		}
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}


// HELPER FUNCTIONS 
func shortenRequestHelper(a *App, w http.ResponseWriter, rb io.ReadCloser) *shortenRequest {
	urlRequest := &shortenRequest{}
	var reqSizeErr *http.MaxBytesError

	// limit request body to 1 MB	
	const maxRequestSize = 1024 * 1024
	limitedRead := http.MaxBytesReader(w, rb, maxRequestSize)

	// Read and Decode POST Request body (JSON) to Go
	dec := json.NewDecoder(limitedRead)
	dec.DisallowUnknownFields()

	err_decode := dec.Decode(urlRequest)
	if err_decode != nil {
		if errors.As(err_decode, &reqSizeErr) {
			a.EntityTooLarge(w)
		}

		a.clientError(w, http.StatusBadRequest)
		return nil
	}

	// url nil value check
	if urlRequest.OriginalURL == "" {
		a.clientError(w, http.StatusBadRequest)
		return nil
	}
	// url validation
	if !isValid(urlRequest.OriginalURL) {
		a.unprocessableEntity(w)
		return nil
	}

	return urlRequest
}

func isValid(ori_url string) bool {
	u, err := url.Parse(ori_url)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != ""
}
