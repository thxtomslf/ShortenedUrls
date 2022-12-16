package storage

import (
	"ShortenedUrls/internal/domain"
	"database/sql"
)

type Urls interface {
	GetSourceUrl(shortUrl *domain.Urls) (*domain.Urls, error)
	PutNewUrlPair(urls *domain.Urls) error
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UrlsRepository: &UrlsRepository{db},
	}
}

type Storage struct {
	UrlsRepository Urls
}

func Open(driver string) (*sql.DB, error) {
	settings := "user=postgres password=password dbname=url_service port=5432 host=db sslmode=disable"
	db, err := sql.Open(driver, settings)
	return db, err
}

func Close(db *sql.DB) error {
	err := db.Close()
	return err
}
