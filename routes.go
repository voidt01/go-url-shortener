package main

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", app.URLHandler.Shortening)
	mux.HandleFunc("GET /{shortCode}", app.URLHandler.Redirecting)

	return app.logRequest(secureHeaders(mux))
}
