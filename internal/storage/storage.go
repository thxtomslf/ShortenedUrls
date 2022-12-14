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

func Open() (*sql.DB, error) {
	settings := "user=postgres password=popisau11 dbname=shorturlsservicedatabase port=5432 host=localhost sslmode=disable"
	db, err := sql.Open("postgres", settings)
	return db, err
}

func Close(db *sql.DB) error {
	err := db.Close()
	return err
}
