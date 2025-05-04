package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", Shorten)

	log.Print("starting a server on:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

const addr = "http://localhost:4000/"

type urlStore struct {
	urls map[string]string
}

type postUrlRequest struct {
	OriginalURL string `json:"original_url"`
}

type postUrlResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	urlDatabase := &urlStore{
		urls: make(map[string]string),
	}
	urlRequest := &postUrlRequest{}
	urlResponse := &postUrlResponse{}

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
	urlDatabase.urls[key] = value

	// creating post response struct
	urlResponse.OriginalURL = value
	urlResponse.ShortenURL = addr + key

	// encoding response struct (G0) to JSON
	w.Header().Set("Content-Type", "application/json")
	err_encode := json.NewEncoder(w).Encode(urlResponse)

	if err_encode != nil {
		log.Printf("Error encoding response: %v", err_encode)
	}

	fmt.Fprintf(w, "%+v\n", urlDatabase)
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
