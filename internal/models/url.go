package models

import (
	"database/sql"
	"math/rand"
	"time"
)

type url struct {
	ID          int
	ShortCode   string
	OriginalUrl string
	CreatedAt   time.Time
}

type UrlModel struct {
	DB *sql.DB
}

func (u *UrlModel) Insert(original_url string) string {
	return ""
}

func (u *UrlModel) Get(short_code string) string {
	return ""
}

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
