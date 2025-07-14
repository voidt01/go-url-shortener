package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}

type UrlHandler struct {
	service *urlService
}

func NewUrlHandler(service *urlService) *UrlHandler {
	return &UrlHandler{service: service}
}

func (uh *UrlHandler) Shortening(w http.ResponseWriter, r *http.Request) {
	var req *shortenRequest = new(shortenRequest)

	// decoding JSON to Go obj
	msg, HTTPStatus, err := readJSON(w, r, req)
	if err != nil {
		writeJSON(w, msg, HTTPStatus, "error")
		return
	}

	// url service
	url, shortCode, err := uh.service.ShortenUrl(req.OriginalURL)
	if err != nil {
		switch {
		case errors.Is(err, ErrUrlInvalid):
			writeJSON(w, "invalid url type", http.StatusUnprocessableEntity, "error")
		default:
			writeJSON(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError, "error")
		}
	}

	// creating post response struct
	resp := &shortenResponse{
		OriginalURL: url,
		ShortenURL:  shortCode,
	}

	// encoding response struct (G0) to JSON
	err_encode := writeJSON(w, &resp, http.StatusCreated, "success")
	if err_encode != nil {
		writeJSON(w, "The server encountered a problem and couldn't process your request", http.StatusInternalServerError, "error")
	}

}

func (uh *UrlHandler) Redirecting(w http.ResponseWriter, r *http.Request) {
	short_code := r.URL.Path[1:]

	original_url, err := uh.service.ResolveUrl(short_code)
	if err != nil {
		switch {
		case errors.Is(err, ErrUrlNotFound):
			http.Error(w, "This link doesn't exist", http.StatusNotFound)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, original_url, http.StatusFound)
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) (string, int, error) {
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
			return msg, http.StatusBadRequest, err_decode

		case errors.Is(err_decode, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return msg, http.StatusBadRequest, err_decode

		case errors.As(err_decode, &unmarshalErr):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at the position %d)", unmarshalErr.Field, unmarshalErr.Offset)
			return msg, http.StatusBadRequest, err_decode

		case strings.HasPrefix(err_decode.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err_decode.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return msg, http.StatusBadRequest, err_decode

		case errors.Is(err_decode, io.EOF):
			msg := "Request body must not be empty"
			return msg, http.StatusBadRequest, err_decode

		case errors.As(err_decode, &reqSizeErr):
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", reqSizeErr.Limit)
			return msg, http.StatusBadRequest, err_decode
		}
	}

	err := dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contains a single JSON object"
		return msg, http.StatusBadRequest, err_decode
	}

	msg := "Internal server error"
	return msg, http.StatusInternalServerError, err_decode
}

func writeJSON(w http.ResponseWriter, data any, statusCode int, status string) error {
	js, err := json.Marshal(map[string]any{
		"status": status,
		"data":   data,
	})

	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}
