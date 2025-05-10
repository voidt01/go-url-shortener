package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/voidt01/go-url-shortener/internal/models"
)

type Config struct {
	port    string
	baseURL string
}

type App struct {
	configApp *Config
	urls      *models.UrlStore
}

func main() {
	var config Config

	flag.StringVar(&config.port, "addr", ":4000", "HTTP Network Address Port")
	flag.StringVar(&config.baseURL, "base-url", "http://localhost", "Base URL for URL Shortener Service")
	flag.Parse()

	app := App{
		configApp: &config,
		urls:      models.NewStoreURL(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", app.Shortening)
	mux.HandleFunc("GET /{shortCode}", app.Redirecting)

	log.Printf("starting a server on: %s", config.port)
	err := http.ListenAndServe(config.port, mux)
	log.Fatal(err)
}
