package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *Application) serve() error {
	srv := &http.Server{
		Addr:         app.config.port,
		ErrorLog:     app.errorLog,
		Handler:      app.Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigChannel
			app.infoLog.Printf("received OS signal: %s", sig.String())
		}
	}()

	app.infoLog.Printf("starting a server on %s", app.config.port)
	return srv.ListenAndServe()
}

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", app.URLHandler.Shortening)
	mux.HandleFunc("GET /{shortCode}", app.URLHandler.Redirecting)

	return app.logRequest(app.secureHeaders(mux))
}

func (app *Application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
