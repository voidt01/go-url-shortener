package url

import (
	"database/sql"
	"errors"
	"time"
)

type URL struct {
	Id          int
	ShortCode   string
	OriginalUrl string
	CreatedAt   time.Time
}

type urlStore struct {
	db *sql.DB
}

func NewUrlStore(db *sql.DB) *urlStore {
	return &urlStore{db: db}
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

func (u *urlStore) GetUrl(short_code string) (*URL, error) {
	var url *URL = new(URL)

	err := u.db.QueryRow("SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = ?", short_code).Scan(&url.Id, &url.OriginalUrl, &url.ShortCode, &url.CreatedAt)
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
