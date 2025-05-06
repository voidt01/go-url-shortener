package main

import (
	"flag"
	"log"
	"net/http"
)

const addr = "http://localhost:4000/"

var urlDatabase = NewStoreURL()

func main() {
	address := flag.String("address", ":4000", "HTTP Network Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", Shortening)
	mux.HandleFunc("GET /{shortCode}", Redirecting)

	log.Printf("starting a server on: %s", *address)
	err := http.ListenAndServe(*address, mux)
	log.Fatal(err)
}
