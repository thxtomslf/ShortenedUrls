package app

import (
	"ShortenedUrls/internal/delivery/http"
	"ShortenedUrls/internal/server"
	"ShortenedUrls/internal/service"
	"ShortenedUrls/internal/storage"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Run() {
	db, err := storage.Open()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("[DEBUG] db connection failed")
	}
	internalStorage := storage.NewStorage(db)

	services := service.NewServices(internalStorage)

	handler := http.NewHandler(services.UrlCutter)

	server := server.NewServer(handler.Init())

	server.Start()

}
