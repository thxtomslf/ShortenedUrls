package storage

import (
	"ShortenedUrls/internal/domain"
	"database/sql"
)

type UrlsRepository struct {
	db *sql.DB
}

func (repository UrlsRepository) GetSourceUrl(urls *domain.Urls) (*domain.Urls, error) {
	data, err := repository.db.Query("SELECT source_url, short_url FROM public.urls WHERE short_url = $1", urls.ShortUrl)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	data.Next()
	givenUrls := domain.Urls{}
	err = data.Scan(&givenUrls.SourceUrl, &givenUrls.ShortUrl)
	if err != nil {
		return nil, err
	}

	return &givenUrls, nil
}

func (repository UrlsRepository) PutNewUrlPair(urls *domain.Urls) error {
	shortUrl := urls.ShortUrl
	sourceUrl := urls.SourceUrl

	data, err := repository.db.Query("insert into public.urls (short_url, source_url) values ($1, $2)", shortUrl, sourceUrl)
	if err != nil {
		return err
	}
	defer data.Close()

	return err
}
