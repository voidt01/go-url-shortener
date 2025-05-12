package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/voidt01/go-url-shortener/internal/models"
)

type Config struct {
	port    string
	baseURL string
}

type App struct {
	configApp *Config
	urls      *models.UrlStore
	errorLog  *log.Logger
	infoLog   *log.Logger
}

func main() {
	// app's config
	var config Config
	flag.StringVar(&config.port, "addr", ":4000", "HTTP Network Address Port")
	flag.StringVar(&config.baseURL, "base-url", "http://localhost", "Base URL for URL Shortener Service")
	flag.Parse()

	// app's logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := App{
		configApp: &config,
		urls:      models.NewStoreURL(),
		errorLog:  errorLog,
		infoLog:   infoLog,
	}

	srv := &http.Server{
		Addr:     config.port,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("starting a server on: %s", config.port)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
