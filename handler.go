package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func Shortening(w http.ResponseWriter, r *http.Request) {
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
	urlDatabase.Set(key, value)

	// creating post response struct
	urlResponse.OriginalURL = value
	urlResponse.ShortenURL = addr + key

	// encoding response struct (G0) to JSON
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err_encode := json.NewEncoder(w).Encode(urlResponse)

	if err_encode != nil {
		log.Printf("Error encoding response: %v", err_encode)
	}

}

func Redirecting(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	value, ok := urlDatabase.Get(key)
	if !ok {
		log.Printf("There is no value for key: %s", key)
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}

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
