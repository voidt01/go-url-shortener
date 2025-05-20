package main

import "net/http"

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", a.Shortening)
	mux.HandleFunc("GET /{shortCode}", a.Redirecting)

	return mux
}
	