package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func (a *App) Shortening(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	urlRequest := &shortenRequest{}
	urlResponse := &shortenResponse{}

	// Decode POST Request body (JSON) to Go
	err_decode := json.NewDecoder(r.Body).Decode(urlRequest)

	if err_decode != nil {
		log.Printf("Error decoding request: %v", err_decode)
		http.Error(w, err_decode.Error(), http.StatusBadRequest)
		return
	}

	// Generate Short URL & Store urls in Database
	key := generateURL()
	value := urlRequest.OriginalURL
	a.urls.Set(key, value)

	// creating post response struct
	urlResponse.OriginalURL = value
	urlResponse.ShortenURL = a.configApp.baseURL + a.configApp.port + "/" + key

	// encoding response struct (G0) to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err_encode := json.NewEncoder(w).Encode(urlResponse)

	if err_encode != nil {
		log.Printf("Error encoding response: %v", err_encode)
	}

}

func (a *App) Redirecting(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	value, ok := a.urls.Get(key)
	if !ok {
		log.Printf("There is no value for key: %s", key)
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}

// Request and Response for creating shortenURL
type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}

// Url Shortener Logic
func generateURL() string {
	url := generator(6)
	return url
}

func generator(digitString int) string {
	var randomizer = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var alphaNumeric = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	tempURL := make([]rune, digitString)
	for i := range tempURL {
		tempURL[i] = alphaNumeric[randomizer.Intn(len(alphaNumeric))]
	}

	return string(tempURL)
}
