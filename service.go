package main

import (
	"math/rand"
	"time"
)

type urlService struct{
	store *urlStore
}

func (us *urlService) ShortenUrl() {}

func (us *urlService) ResolveUrl(short_code string) (string, error) {
	original_url, err := us.store.GetUrl(short_code)
	if err != nil {
		return "", err
	}
	return original_url, nil
}

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