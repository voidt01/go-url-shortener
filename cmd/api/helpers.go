package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

func (a *App) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

type clientError struct {
	status int 
	msg string
}

func (ce *clientError) Error() string {
	return ce.msg
}

func (a *App) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// checking if the Content-Type: application/json	
	ct := r.Header.Get("Content-Type")
	if ct != ""{
		if !strings.HasPrefix(ct, "application/json"){
			msg := "Content-Type must be application/json"
			return &clientError{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	// limit request body to 1 MB	
	const maxRequestSize = 1024 * 1024
	limitedRead := http.MaxBytesReader(w, r.Body, maxRequestSize)

	// Read and Decode POST Request body (JSON) to Go
	dec := json.NewDecoder(limitedRead)
	dec.DisallowUnknownFields()

	err_decode := dec.Decode(dst)
	if err_decode != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError
		var reqSizeErr *http.MaxBytesError

		switch {
		case errors.As(err_decode, &syntaxErr):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxErr.Offset)
			return &clientError{status: http.StatusBadRequest, msg: msg}
		
		case errors.Is(err_decode, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &clientError{status: http.StatusBadRequest, msg: msg}
		
		case errors.As(err_decode, &unmarshalErr):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at the position %d)", unmarshalErr.Field, unmarshalErr.Offset)
			return &clientError{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err_decode.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err_decode.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &clientError{status: http.StatusUnsupportedMediaType, msg: msg}

		case errors.Is(err_decode, io.EOF):
			msg := "Request body must not be empty"
			return &clientError{status: http.StatusBadRequest, msg: msg}

		case errors.As(err_decode, &reqSizeErr):
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", reqSizeErr.Limit)
			return &clientError{status: http.StatusRequestEntityTooLarge, msg: msg}
		
		default:
			return err_decode
		}
	}

	err := dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contains a single JSON object"
		return &clientError{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func (a *App) isValid(ori_url string) error {
	u, err := url.Parse(ori_url)
	if !(err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "") {
		msg := "URL must be valid"
		return &clientError	{status: http.StatusBadRequest, msg: msg}
	}
	return nil
}