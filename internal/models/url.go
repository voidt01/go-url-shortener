package models

import (
	"sync"
)

type UrlStore struct {
	urls map[string]string
	mu sync.Mutex
}

func NewStoreURL() *UrlStore {
	return &UrlStore{
		urls: make(map[string]string),
	}
}

func (u *UrlStore) Set(key, value string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.urls[key] = value
}

func (u *UrlStore) Get(key string) (value string, ok bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	value, ok = u.urls[key]
	return value, ok
}


