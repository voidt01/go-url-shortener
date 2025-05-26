package models

import (
	"database/sql"
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

type UrlModel struct {
	DB *sql.DB
}

func (u *UrlModel) Insert(original_url string) (short_code string, err error) {
	short_code = generateURL()

	stmt := `INSERT INTO urls(short_code, original_url)
	VALUES(?, ?)`

	_, err = u.DB.Exec(stmt, short_code, original_url)
	if err != nil {
		return "", err
	}

	return short_code, nil
}

func (u *UrlModel) Get(short_code string) (original_url string, err error) {

	err = u.DB.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", short_code).Scan(&original_url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		} else {
			return "", err
		}
	}

	return original_url, nil

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
