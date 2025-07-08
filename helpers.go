package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
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
			return &clientError{status: http.StatusBadRequest, msg: msg}

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

func encodeJSON(w http.ResponseWriter, data any, status int) error {
	js, err := json.Marshal(map[string]any{
		"status": "success",
		"data":   data,
	})

	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func urlValidation(ori_url string) (string, error) {
	u, err := url.Parse(ori_url)
	if err != nil {
		msg := "URL must be valid"
		return "", &clientError{status: http.StatusBadRequest, msg: msg}
	}

	if u.Scheme == "" {
		ori_url = "https://" + ori_url
		u, err = url.Parse(ori_url)
		if err != nil {
			msg := "URL must be valid"
			return "", &clientError{status: http.StatusBadRequest, msg: msg}
		}
	}
	if !(u.Scheme == "https" || u.Scheme == "http") {
		msg := "URL must use http or https protocol"
		return "", &clientError{status: http.StatusBadRequest, msg: msg}
	}

	if u.Host == "" {
		msg := "URL must have a valid host"
		return "", &clientError{status: http.StatusBadRequest, msg: msg}
	}
	return ori_url, nil
}

func (a *App) builderShortenURL(shortCode string) string {
	url := fmt.Sprintf("%s%s/%s", a.config.Server.baseURL, a.config.Server.port, shortCode)
	return url
}
