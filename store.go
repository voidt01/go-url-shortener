package main

import (
	"database/sql"
	"errors"
)

type urlStore struct {
	db *sql.DB
}

func (u *urlStore) SetUrl(original_url, short_code string) error {

	stmt := `INSERT INTO urls(short_code, original_url)
	VALUES(?, ?)`

	_, err := u.db.Exec(stmt, short_code, original_url)
	if err != nil {
		return err
	}

	return nil
}

func (u *urlStore) GetUrl(short_code string) (original_url string, err error) {

	err = u.db.QueryRow("SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = ?", short_code).Scan(&original_url)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrUrlNotFound
		default:
			return "", err
		}
	}

	return original_url, nil

}
