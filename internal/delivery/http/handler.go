package http

import (
	"ShortenedUrls/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	urlCutterService service.UrlCutter
}

func NewHandler(urlCutService service.UrlCutter) *Handler {
	return &Handler{
		urlCutService,
	}
}

func (handler *Handler) Init() *gin.Engine {
	router := gin.New()

	// формат запроса: {"short_url": "some_hash"}
	// код ответа при успехе - 302
	router.GET("/:url", handler.urlCutterService.GetSourceURL)
	// формат запроса: {"source_url": "yourUrl.com"}
	// формат ответа при коде 200: {"short_url": "some_hash", "source_url": "someLink.com"}
	router.POST("/service/cut", handler.urlCutterService.GetCutUrl)

	return router
}
