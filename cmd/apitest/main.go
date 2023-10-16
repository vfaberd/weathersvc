package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const forecastAPIEndpoint = "https://api.open-meteo.com/v1/forecast?latitude=43.2567&longitude=76.9286&current=temperature_2m,relativehumidity_2m,precipitation_probability,visibility,windspeed_10m&timezone=GMT"

const TimeFormat = "2006-01-02T15:04"

type weather struct {
	Latitude                 float64   `json:"latitude"`
	Longitude                float64   `json:"longitude"`
	Timestamp                time.Time `json:"timestamp"`
	Temperature2M            float64   `json:"temperature_2m"`
	Relativehumidity2M       int       `json:"relativehumidity_2m"`
	PrecipitationProbability int       `json:"precipitation_probability"`
	Visibility               float64   `json:"visibility"`
	Windspeed10M             float64   `json:"windspeed_10m"`
}

type currentWeatherResp struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Time                     string  `json:"time"`
		Interval                 int     `json:"interval"`
		Temperature2M            float64 `json:"temperature_2m"`
		Relativehumidity2M       int     `json:"relativehumidity_2m"`
		PrecipitationProbability int     `json:"precipitation_probability"`
		Visibility               float64 `json:"visibility"`
		Windspeed10M             float64 `json:"windspeed_10m"`
	} `json:"current"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	resp, err := http.Get(forecastAPIEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var r currentWeatherResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}
	w, err := respToWeather(r)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", w)

	json.NewEncoder(os.Stdout).Encode(w)
	return nil
}

func respToWeather(resp currentWeatherResp) (weather, error) {
	t, err := time.Parse(TimeFormat, resp.Current.Time)
	if err != nil {
		return weather{}, fmt.Errorf("parse resp.Current.Time: %w", err)
	}
	return weather{
		Latitude:                 resp.Latitude,
		Longitude:                resp.Longitude,
		Timestamp:                t,
		Temperature2M:            resp.Current.Temperature2M,
		Relativehumidity2M:       resp.Current.Relativehumidity2M,
		PrecipitationProbability: resp.Current.PrecipitationProbability,
		Visibility:               resp.Current.Visibility,
		Windspeed10M:             resp.Current.Windspeed10M,
	}, nil
}
