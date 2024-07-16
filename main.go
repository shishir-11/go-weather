package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	WeatherAPIKey string `json:"WeatherAPIKey"`
}

type weatherData struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Weather struct {
		Temp float64 `json:"temp_c"`
		Wind float64 `json:"wind_mph"`
	} `json:"current"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiConfig.WeatherAPIKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			city := strings.Split(path, "/")[2]
			// fmt.Println(city)
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go]n"))
}
