package core

import (
	"context"
	"time"
)

const TimestampFormat = "2006-01-02T15:04"

type Weather struct {
	Location                 Location  `json:"location"`
	Timestamp                time.Time `json:"timestamp"`
	Temperature2M            float64   `json:"temperature_2m"`
	Relativehumidity2M       int       `json:"relativehumidity_2m"`
	PrecipitationProbability int       `json:"precipitation_probability"`
	Visibility               float64   `json:"visibility"`
	Windspeed10M             float64   `json:"windspeed_10m"`
}

type WeatherStorage interface {
	Save(ctx context.Context, weather Weather) error
	GetLatest(ctx context.Context, location string) (Weather, error)
	GetPeriod(ctx context.Context, location string, from, to time.Time) ([]Weather, error)
}

type WeatherFetcher interface {
	GetCurrent(ctx context.Context, loc Location) (Weather, error)
}