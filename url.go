package main

import (
	"errors"
	"math/rand"
	"time"
)

type Url struct {
	ID          int
	ShortCode   string
	OriginalUrl string
	CreatedAt   time.Time
}

var (
	ErrUrlNotFound = errors.New("models: No matching url found")
	ErrUrlInvalid  = errors.New("models: Url format is invalid")
)

// Url Shortener Logic
func generateURL() string {
	url := generator(6)
	return url
}

func generator(digitString int) string {
	var randomizer = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var alphaNumeric = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	tempURL := make([]rune, digitString)
	for i := range tempURL {
		tempURL[i] = alphaNumeric[randomizer.Intn(len(alphaNumeric))]
	}

	return string(tempURL)
}
