package kpp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

func getHTML(url string) ([]byte, error) {
	var body []byte
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
		var ch = toUtf(char)
		fmt.Fprintf(buffer, "%c", ch)
	}
	doc := buffer.Bytes()
	return doc, nil
}

// GetRating - получение рейтингов
func GetRating(name string, year int64) (KP, error) {
	var kp KP
	yearStr := fmt.Sprintf("%d", year)
	url, err := urlEncoded(apiURL + strings.Replace(name, " ", "+", -1) + "+" + yearStr + "/view/movie/")
	if err != nil {
		return kp, err
	}
	body, err := getHTML(url)
	if err != nil {
		return kp, err
	}
	findStr := regexp.QuoteMeta(name + ", " + yearStr)
	reHref := regexp.MustCompile(`<a href="(.*?)">` + findStr + `<\/a>`)
	reK := regexp.MustCompile(`<b>рейтинг фильма:</b>.*?<i>(.*?)</i>`)
	reI := regexp.MustCompile(`<b>рейтинг IMDB:</b>.*?<i>(.*?)</i>`)
	if reHref.Match(body) == true {
		findHref := reHref.FindSubmatch(body)
		href := string(findHref[1])
		body, err = getHTML(href)
		if err != nil {
			return kp, err
		}
		if reK.Match(body) == true {
			kindK := reK.FindSubmatch(body)
			kp.Kinopoisk, _ = strconv.ParseFloat(string(kindK[1]), 64)
		}
		if reI.Match(body) == true {
			kindI := reI.FindSubmatch(body)
			kp.IMDb, _ = strconv.ParseFloat(string(kindI[1]), 64)
		}
	}
	return kp, nil
}
