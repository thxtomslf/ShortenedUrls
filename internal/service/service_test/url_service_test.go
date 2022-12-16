package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type FirstTestGroupRequest struct {
	name               string
	method             string
	shortUrl           string
	expectedSourceUrl  string
	expectedStatusCode int
	testFunction       func(req FirstTestGroupRequest, t *testing.T, resp *http.Response)
}

func TestUrlCutterService(t *testing.T) {
	// создаем клиента
	client := &http.Client{}

	// список корректных запросов для заполнения данных хранилища
	fillingRequestsTable := []struct {
		method      string
		requestBody string
		url         string
		source      string
	}{
		{
			"POST",
			"{\"source_url\": \"google.com\"}",
			"http://0.0.0.0:8080/service/cut",
			"google.com",
		},
		{
			// намеренно сделано два одинаковых запросов для демонстрации корректности работы
			// второго метода сервиса CutterService
			"POST",
			"{\"source_url\": \"google.com\"}",
			"http://0.0.0.0:8080/service/cut",
			"google.com",
		},
		{
			"POST",
			"{\"source_url\": \"music.yandex.ru\"}",
			"http://0.0.0.0:8080/service/cut",
			"music.yandex.ru",
		},
		{
			"POST",
			"{\"source_url\": \"https://www.ozon.ru/\"}",
			"http://0.0.0.0:8080/service/cut",
			"https://www.ozon.ru/",
		},
		{
			"POST",
			"{\"source_url\": \"https://stackoverflow.com/questions/43041331/docker-forever-in-docker-is-starting-at-windows-task\"}",
			"http://0.0.0.0:8080/service/cut",
			"https://stackoverflow.com/questions/43041331/docker-forever-in-docker-is-starting-at-windows-task",
		},
		{
			"POST",
			"{\"source_url\": \"http://www.mathtask.ru/0033-two-dimensional-random-variables.php#:~:text=%D0%94%D0%B2%D1%83%D0%BC%D0%B5%D1%80%D0%BD%D0%BE%D0%B9%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D0%B9%D0%BD%D0%BE%D0%B9%20%D0%B2%D0%B5%D0%BB%D0%B8%D1%87%D0%B8%D0%BD%D0%BE%D0%B9%20%D0%BD%D0%B0%D0%B7%D1%8B%D0%B2%D0%B0%D0%B5%D1%82%D1%81%D1%8F%20%D1%84%D1%83%D0%BD%D0%BA%D1%86%D0%B8%D1%8F,%D1%85%20%D0%B8%20y%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D0%B9%D0%BD%D1%8B%D1%85%20%D0%B7%D0%BD%D0%B0%D1%87%D0%B5%D0%BD%D0%B8%D0%B9.&text=X%20%D0%B8%20Y%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D0%B9%D0%BD%D1%8B%D0%B5%20%D0%B2%D0%B5%D0%BB%D0%B8%D1%87%D0%B8%D0%BD%D1%8B,%D0%BA%D0%B0%D0%BA%20%D0%B4%D0%B8%D1%81%D0%BA%D1%80%D0%B5%D1%82%D0%BD%D1%8B%D0%BC%D0%B8%2C%20%D1%82%D0%B0%D0%BA%20%D0%B8%20%D0%BD%D0%B5%D0%BF%D1%80%D0%B5%D1%80%D1%8B%D0%B2%D0%BD%D1%8B%D0%BC%D0%B8.\"}",
			"http://0.0.0.0:8080/service/cut",
			"http://www.mathtask.ru/0033-two-dimensional-random-variables.php#:~:text=%D0%94%D0%B2%D1%83%D0%BC%D0%B5%D1%80%D0%BD%D0%BE%D0%B9%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D0%B9%D0%BD%D0%BE%D0%B9%20%D0%B2%D0%B5%D0%BB%D0%B8%D1%87%D0%B8%D0%BD%D0%BE%D0%B9%20%D0%BD%D0%B0%D0%B7%D1%8B%D0%B2%D0%B0%D0%B5%D1%82%D1%81%D1%8F%20%D1%84%D1%83%D0%BD%D0%BA%D1%86%D0%B8%D1%8F,%D1%85%20%D0%B8%20y%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D0%B9%D0%BD%D1%8B%D1%85%20%D0%B7%D0%BD%D0%B0%D1%87%D0%B5%D0%BD%D0%B8%D0%B9.&text=X%20%D0%B8%20Y%20%D1%81%D0%BB%D1%83%D1%87%D0%B0%D0%B9%D0%BD%D1%8B%D0%B5%20%D0%B2%D0%B5%D0%BB%D0%B8%D1%87%D0%B8%D0%BD%D1%8B,%D0%BA%D0%B0%D0%BA%20%D0%B4%D0%B8%D1%81%D0%BA%D1%80%D0%B5%D1%82%D0%BD%D1%8B%D0%BC%D0%B8%2C%20%D1%82%D0%B0%D0%BA%20%D0%B8%20%D0%BD%D0%B5%D0%BF%D1%80%D0%B5%D1%80%D1%8B%D0%B2%D0%BD%D1%8B%D0%BC%D0%B8.",
		},
	}

	receivedUrlsList := make([]map[string]string, len(fillingRequestsTable))
	// делаем запросы, заполняя список сокращенных урлов
	for ind, fillingRequest := range fillingRequestsTable {
		req, err := http.NewRequest(fillingRequest.method, fillingRequest.url,
			bytes.NewBufferString(fillingRequest.requestBody))
		if err != nil {
			t.Fatalf("[DEBUG] %s", err.Error())
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[DEBUG] %s", err.Error())
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("[DEBUG] %s", err.Error())
		}
		parsedRespBody := make(map[string]string)
		err = json.Unmarshal(respBody, &parsedRespBody)
		if err != nil {
			t.Fatalf("[DEBUG] %s", err.Error())
		}

		// сохраняем ответ сервера
		receivedUrlsList[ind] = parsedRespBody
		// заменяем полученный урл источника на тот, который передвавали в запросе для дальнейшей проверки
		receivedUrlsList[ind]["source_url"] = fillingRequest.source
	}
	// список запросов, тестирующий обработки GetSourceUrl
	getSrcTestTable := []FirstTestGroupRequest{
		{
			"GetTest1",
			"GET",
			"http://0.0.0.0:8080/" + receivedUrlsList[0]["short_url"],
			receivedUrlsList[0]["source_url"],
			http.StatusPermanentRedirect,
			GetTest,
		},
		{
			"PostTest",
			"POST",
			"http://0.0.0.0:8080/" + receivedUrlsList[1]["short_url"],
			"",
			http.StatusNotFound,
			PostTest,
		},
		{
			"GetTest2",
			"GET",
			"http://0.0.0.0:8080/" + receivedUrlsList[1]["short_url"],
			receivedUrlsList[1]["source_url"],
			http.StatusPermanentRedirect,
			GetTest,
		},
		{
			"PutTest",
			"PUT",
			"http://0.0.0.0:8080/" + receivedUrlsList[1]["short_url"],
			"",
			http.StatusNotFound,
			PutTest,
		},
		{
			"GetTest3",
			"GET",
			"http://0.0.0.0:8080/" + receivedUrlsList[3]["short_url"],
			receivedUrlsList[3]["source_url"],
			http.StatusPermanentRedirect,
			GetTest,
		},
		{
			"NotExistingShortLinkTest",
			"GET",
			"http://0.0.0.0:8080/" + "I_AM_NOT_EXCITING",
			"",
			http.StatusBadRequest,
			NotExistingShortLink,
		},
	}

	for _, testCase := range getSrcTestTable {
		req, err := http.NewRequest(testCase.method, testCase.shortUrl, nil)
		if err != nil {
			t.Fatalf("[DEBUG] %s", err.Error())
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[DEBUG] %s", err.Error())
		}

		testRequest := func(t2 *testing.T) { testCase.testFunction(testCase, t2, resp) }

		t.Run(testCase.name, testRequest)

	}

}

func GetTest(testCase FirstTestGroupRequest, t *testing.T, resp *http.Response) {
	host := strings.Trim(resp.Request.URL.Host, "www.")
	if strings.Contains(testCase.expectedSourceUrl, host) {
		assert.Equal(t, 1, 1)
	} else {
		assert.Equal(t, 0, 1)
	}
	assert.Equal(t, testCase.expectedStatusCode/100, resp.Request.Response.StatusCode/100)

	fmt.Println()
}

func PostTest(testCase FirstTestGroupRequest, t *testing.T, resp *http.Response) {
	assert.Equal(t, resp.StatusCode, testCase.expectedStatusCode)
	fmt.Println()
}

func PutTest(testCase FirstTestGroupRequest, t *testing.T, resp *http.Response) {
	assert.Equal(t, resp.StatusCode, testCase.expectedStatusCode)
	fmt.Println()
}

func NotExistingShortLink(testCase FirstTestGroupRequest, t *testing.T, resp *http.Response) {
	assert.Equal(t, resp.StatusCode, testCase.expectedStatusCode)
	fmt.Println()
}
