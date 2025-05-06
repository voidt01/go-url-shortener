package main

import (
	"sync"
)

type urlStore struct {
	urls map[string]string
	mu sync.Mutex
}

func NewStoreURL() *urlStore {
	return &urlStore{
		urls: make(map[string]string),
	}
}

func (u *urlStore) Set(key, value string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.urls[key] = value
}

func (u *urlStore) Get(key string) (value string, ok bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	value, ok = u.urls[key]
	return value, ok
}

// Request and Response for creating shortenURL
type shortenRequest struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"shorten_url"`
}
