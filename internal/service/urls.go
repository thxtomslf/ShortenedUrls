package service

import (
	"ShortenedUrls/internal/domain"
	"ShortenedUrls/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UrlCutter interface {
	GetCutUrl(ctx *gin.Context)
	GetSourceURL(ctx *gin.Context)
}

type UrlCutterService struct {
	UrlsRepository storage.Urls
}

func (cutter *UrlCutterService) GetSourceURL(ctx *gin.Context) {
	shortUrl := ctx.Request.URL.Path

	urls := &domain.Urls{
		ShortUrl: strings.Trim(shortUrl, "/"),
	}

	gotUrls, err := cutter.UrlsRepository.GetSourceUrl(urls)
	if err != nil {
		fmt.Println("[DEBUG5] " + err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	// если не указывали протокол, надо добавить для корректного редиректа
	if !strings.Contains(gotUrls.SourceUrl, "://") {
		gotUrls.SourceUrl = "https://" + gotUrls.SourceUrl
	}

	ctx.Redirect(http.StatusPermanentRedirect, gotUrls.SourceUrl)
}

func (cutter *UrlCutterService) GetCutUrl(ctx *gin.Context) {

	urls := &domain.Urls{}
	err := json.NewDecoder(ctx.Request.Body).Decode(&urls)
	if err != nil {
		fmt.Println("[DEBUG1] " + err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	if urls.SourceUrl == "" {
		fmt.Println("[DEBUG2] incorrect request body")
		ctx.Status(http.StatusBadRequest)
		return
	}

	path := cutter.cutUrl(urls.SourceUrl)

	urls.ShortUrl = string(path)

	existenceState, gotUrls := cutter.checkSourceUrlPresence(urls)
	if existenceState {
		urls = gotUrls
	} else {
		err = cutter.UrlsRepository.PutNewUrlPair(urls)
		if err != nil {
			fmt.Println("[DEBUG3] database create field error: ", err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
	ctx.JSON(http.StatusOK, urls)
	if err != nil {
		fmt.Println("[DEBUG4] url sending error")
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (cutter *UrlCutterService) cutUrl(url string) []rune {
	alphabet := fillAlphabet()

	path := make([]rune, 10)
	hash := cutter.getHash([]rune(url), len(alphabet))

	for i := 0; i < 10; i++ {
		path[i] = alphabet[hash[i]]
	}

	return path
}

// разрешаем коллизии хешей методом линейного пробирования
func (cutter *UrlCutterService) getHash(url []rune, diapason int) []rune {
	var linearProbeFactor rune = 7

	key := make([]rune, 0, 10)

	var j rune
	for j = 0; j < 10; j++ {
		var hashFactor rune = 4
		for i := 0; i < 10; i++ {
			key = append(key, cutter.hashString(url, hashFactor, diapason)+j*linearProbeFactor)
			hashFactor++
			for hashFactor%63 == 0 { // требуется, чтобы длина диапазона и множитель хеширования были взаимно простыми
				hashFactor++ // в этот цикл программа практически не заходит
			}
		}

		existenceState, _ := cutter.checkSourceUrlPresence(&domain.Urls{ShortUrl: string(key)})

		if !existenceState {
			return key
		}
	}
	return nil
}

// проверяем наличие ссылки-источника в базе
func (cutter *UrlCutterService) checkSourceUrlPresence(urls *domain.Urls) (bool, *domain.Urls) {
	gotUrls, err := cutter.UrlsRepository.GetSourceUrl(urls)
	if err != nil {
		return false, nil
	}

	return true, gotUrls
}

// вычисляем хеш строки методом умножения с использванием схемы Горнера для ускорения вычислений
func (cutter *UrlCutterService) hashString(key []rune, factor rune, diapason int) rune {
	hashFactor := factor
	diapasonLength := rune(diapason)
	var hash rune = 0

	for _, value := range key {
		hash = (hash*hashFactor + value) % diapasonLength
	}

	return hash

}

func fillAlphabet() []rune {
	alphabet := make([]rune, 63)
	var index int = 0
	var i rune
	for i = 48; i < 58; i++ { // Добавляем коды цифр от 0 до 9
		alphabet[index] = i
		index++
	}
	for i = 65; i < 91; i++ { // Добавляем коды прописных латинских букв
		alphabet[index] = i
		index++
	}
	for i = 97; i < 123; i++ { // Добавляем коды строчных латинских букв
		alphabet[index] = i
		index++
	}
	alphabet[index] = 95 // Добавляем код "_"
	return alphabet
}
