package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type URLData struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url,omitempty"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	urlStore := &URLData{}
	err := json.NewDecoder(r.Body).Decode(urlStore)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%+v\n", urlStore)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", shorten)

	log.Print("starting a server on:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
