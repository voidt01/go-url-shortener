package main

import (
	"database/sql"
	"errors"
)

type UrlModel struct {
	DB *sql.DB
}

func (u *UrlModel) SetUrl(original_url, short_code string) error {

	stmt := `INSERT INTO urls(short_code, original_url)
	VALUES(?, ?)`

	_, err := u.DB.Exec(stmt, short_code, original_url)
	if err != nil {
		return err
	}

	return nil
}

func (u *UrlModel) GetUrl(short_code string) (url *Url, err error) {

	err = u.DB.QueryRow("SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = ?", short_code).Scan(url)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUrlNotFound
		default:
			return nil, err
		}
	}

	return url, nil

}
