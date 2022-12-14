package service

import "ShortenedUrls/internal/storage"

type Services struct {
	UrlCutter UrlCutter
}

func NewServices(storage *storage.Storage) *Services {
	return &Services{&UrlCutterService{storage.UrlsRepository}}
}
