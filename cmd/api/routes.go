package main

import "net/http"

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	shortenHandler := http.HandlerFunc(a.Shortening)
	redirectHandler := http.HandlerFunc(a.Redirecting)

	mux.Handle("POST /shorten", a.enforceJSONhandler(shortenHandler))
	mux.Handle("GET /{shortCode}", redirectHandler)

	return a.logRequest(secureHeaders(mux))
}
