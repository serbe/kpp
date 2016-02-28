package kpp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// GetRating - получение рейтингов
func GetRating(name string, year int64) (KP, error) {
	url := strings.Replace(string, " ", "+", -1)
	resp, err := http.Get(apiURL + "/configuration?api_key=" + tmdb.apiKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Status Code %d received from TMDb", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(body, &config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tmdb.config = config
}

// GetByName get data from themoviedb by name
func (tmdb *TMDB) GetByName(movieName string, year string) (tmdbResult, error) {
	tmdb.getConfig()
	time.Sleep(1 * time.Second)
	var response = &tmdbResponse{}
	if year != "" {
		year = "&year=" + year
	}
	queryString := apiURL + "/search/movie?api_key=" + tmdb.apiKey + "&language=ru&query=" + url.QueryEscape(movieName) + year
	resp, err := http.Get(queryString)
	if err != nil {
		return tmdbResult{}, err
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.Header)
		return tmdbResult{}, fmt.Errorf("Status Code %d received from TMDb", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tmdbResult{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return tmdbResult{}, err
	}
	if len(response.Results) == 0 {
		return tmdbResult{}, err
	}
	response.Results[0].Poster_base_url = tmdb.config.Images.Base_url
	return response.Results[0], err
}
