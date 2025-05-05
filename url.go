package main

type urlStore struct {
	urls map[string]string
}

func (u *urlStore) Set(key, value string) {
	u.urls[key] = value
}

func (u *urlStore) Get(key string) (value string, ok bool) {
	value, ok = u.urls[key]
	return value, ok
}

func NewStoreURL() *urlStore {
	return &urlStore{
		urls: make(map[string]string),
	}
}

// Request and Response for creating shortenURL
type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}
