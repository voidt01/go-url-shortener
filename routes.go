package main

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	shortenHandler := http.HandlerFunc(app.Shortening)
	redirectHandler := http.HandlerFunc(app.Redirecting)

	mux.Handle("POST /shorten", app.enforceJSONhandler(shortenHandler))
	mux.Handle("GET /{shortCode}", redirectHandler)

	return app.logRequest(secureHeaders(mux))
}
