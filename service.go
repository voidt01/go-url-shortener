package main

import (
	"math/rand"
	"net/url"
	"time"
)

type urlService struct{
	store *urlStore
}

func NewUrlService(store *urlStore) *urlService{
	return &urlService{store: store}
}

func (us *urlService) ShortenUrl(original_url string) (string, string, error) {
	// url validation & sanitazion
	original_url, err := us.normalizeUrl(original_url)
	if err != nil {
		return "", "", err
	}

	// short code generation
	short_code := us.generateUrl()

	// insert url data to db
	err = us.store.SetUrl(original_url, short_code)
	if err != nil {
		return "", "", err
	}

	return original_url, short_code, nil
}

func (us *urlService) ResolveUrl(short_code string) (string, error) {
	url, err := us.store.GetUrl(short_code)
	if err != nil {
		return "", err
	}
	return url.OriginalUrl, nil
}

func (us *urlService) normalizeUrl(ori_url string) (string, error) {
	u, err := url.Parse(ori_url)
	if err != nil {
		return "", ErrUrlInvalid
	}

	if u.Scheme == "" {
		ori_url = "https://" + ori_url
		u, err = url.Parse(ori_url)
		if err != nil {
			return "", ErrUrlInvalid
		}
	}
	if !(u.Scheme == "https" || u.Scheme == "http") {
		return "", ErrUrlInvalid
	}

	if u.Host == "" {
		return "", ErrUrlInvalid
	}
	return ori_url, nil
}

func (us *urlService) generateUrl() string {
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