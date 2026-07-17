package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	APIKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadAPIConfig(filename string) (*apiConfigData, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config apiConfigData
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}


func query(city string) (*weatherData, error) {
	config, err := loadAPIConfig(".apiConfig")
	if err != nil {
		return nil, err
	}

	apiKey := config.APIKey
	url := "http://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data weatherData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {


	http.HandleFunc("/hello",hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
     city := strings.SplitN(r.URL.Path,"/",3)[2]
	 data, err := query(city)
	 if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		 return
	 }
	 w.Header().Set("Content-Type", "application/json")
	 json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}
