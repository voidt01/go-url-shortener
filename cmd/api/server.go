package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *Application) serve() error {
	srv := &http.Server{
		Addr:         app.config.server.port,
		ErrorLog:     app.errorLog,
		Handler:      app.Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	shutdownErrChan := make(chan error)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan

		app.infoLog.Printf("Shutting down the server: %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErrChan <- srv.Shutdown(ctx)
	}()

	app.infoLog.Printf("starting a server on %s", app.config.server.port)

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	app.infoLog.Printf("Stopped server on %s", app.config.server.port)

	return nil
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
