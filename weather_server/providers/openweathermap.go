package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
)

type OpenWeatherMap struct {
	WeatherKey string
}

const (
	openWeatherMapURL = "http://api.openweathermap.org/data/2.5/weather"
)

type OpenWeatherMapProvider interface {
	Search(string) (WeatherInfo, error)
}

type openWeatherMapResult struct {
	Name string `json:"name"`
	Main struct {
		KelvinTemp    float64 `json:"temp"`
		KelvinTempMin float64 `json:"temp_min"`
		KelvinTempMax float64 `json:"temp_max"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

type WeatherInfo struct {
	Temperature    float64
	TemperatureMin float64
	TemperatureMax float64
	Description    string
	Country        string
	Found          bool
}

func (openweathermap OpenWeatherMap) Search(query string) (WeatherInfo, error) {

	queryURL := fmt.Sprintf("%s?q=%s&appid=%s", openWeatherMapURL, url.QueryEscape(query), openweathermap.WeatherKey)

	body, err := httpClient.get(queryURL)
	if err != nil {
		return WeatherInfo{}, err
	}

	if err != nil {
		log.Fatalln(err)
	}

	var result openWeatherMapResult
	json.Unmarshal(body, &result)

	fmt.Println(string(body))

	return result.asWeatherInfo(), nil
}

func (r openWeatherMapResult) asWeatherInfo() WeatherInfo {
	if r.found() {
		return WeatherInfo{
			Temperature:    r.toCelcius(r.Main.KelvinTemp),
			TemperatureMax: r.toCelcius(r.Main.KelvinTempMax),
			TemperatureMin: r.toCelcius(r.Main.KelvinTempMin),
			Description:    r.description(),
			Country:        r.Sys.Country,
			Found:          true,
		}
	}
	return WeatherInfo{}
}

func (r openWeatherMapResult) toCelcius(kelvin float64) float64 {
	return math.Floor(kelvin-273.15) + 0.5
}

func (r openWeatherMapResult) description() string {
	return r.Weather[0].Description
}

func (r openWeatherMapResult) found() bool {
	return len(r.Weather) > 0
}
