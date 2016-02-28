package kpp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiURL = "http://m.kinopoisk.ru/search/"
)

// KP values:
// Kinopoisk - рейтинг кинопоиска
// IMDB      - рейтинг IMDb
type KP struct {
	Kinopoisk float64
	IMDb      float64
}

func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func getHTML(name string, year int64) ([]byte, error) {
	var body []byte
	url, err := urlEncoded(apiURL + strings.Replace(name, " ", "+", -1) + "+" + fmt.Sprintf("%d", year) + "/")
	if err != nil {
		return body, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	buffer := bytes.NewBufferString("")
	for _, char := range body {
		var ch = Utf(char)
		fmt.Fprintf(buffer, "%c", ch)
	}
	doc := buffer.Bytes()
	return doc, nil
}

// GetRating - получение рейтингов
func GetRating(name string, year int64) (KP, error) {
	var kp KP
	body, err := getHTML(name, year)
	if err != nil {
		return kp, err
	}
	fmt.Println(string(body))
	return kp, nil
}
