package open_meteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"weathersvc/internal/core"
)


type OpenMeteoAPIClient struct {
	forecastEndpoint string
	timeout          time.Duration
}

func NewOpenMeteoAPIClient(forecastEndpoint string, timeout time.Duration) *OpenMeteoAPIClient {
	return &OpenMeteoAPIClient{
		forecastEndpoint: forecastEndpoint,
		timeout:          timeout,
	}
}

type getCurrentResp struct {
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

func (r *getCurrentResp) toWeather() (core.Weather, error) {
	t, err := time.Parse(core.TimestampFormat, r.Current.Time)
	if err != nil {
		return core.Weather{}, fmt.Errorf("parse resp.Current.Time: %w", err)
	}
	loc, err := core.LocationFromCoordinates(r.Latitude, r.Longitude)
	if err != nil {
		return core.Weather{}, err
	}
	return core.Weather{
		Location:                 loc,
		Timestamp:                t,
		Temperature2M:            r.Current.Temperature2M,
		Relativehumidity2M:       r.Current.Relativehumidity2M,
		PrecipitationProbability: r.Current.PrecipitationProbability,
		Visibility:               r.Current.Visibility,
		Windspeed10M:             r.Current.Windspeed10M,
	}, nil
}

// "https://api.open-meteo.com/v1/forecast?latitude=43.2567&longitude=76.9286&current=temperature_2m,relativehumidity_2m,precipitation_probability,visibility,windspeed_10m&timezone=GMT"
// "https://api.open-meteo.com/v1/forecast?current=temperature_2m&current=relativehumidity_2m&current=precipitation_probability&current=visibility&current=windspeed_10m&latitude=43.26&longitude=76.93&timezone=GMT"
// https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Crelativehumidity_2m%2Cprecipitation_probability%2Cvisibility%2Cwindspeed_10m&latitude=43.26&longitude=76.93&timezone=GMT
func (c *OpenMeteoAPIClient) GetCurrent(ctx context.Context, loc core.Location) (core.Weather, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	endpointURL, err := url.Parse(c.forecastEndpoint)
	if err != nil {
		return core.Weather{}, fmt.Errorf("parse forecast endpoint: %w", err)
	}
	query := endpointURL.Query()
	query.Add("current", "temperature_2m,relativehumidity_2m,precipitation_probability,visibility,windspeed_10m")
	query.Add("timezone", "GMT")
	query.Add("latitude", strconv.FormatFloat(loc.Latitude, 'f', 2, 32))
	query.Add("longitude", strconv.FormatFloat(loc.Longitude, 'f', 2, 32))
	endpointURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL.String(), nil)
	if err != nil {
		return core.Weather{}, fmt.Errorf("creating http request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return core.Weather{}, fmt.Errorf("sending http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return core.Weather{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var r getCurrentResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return core.Weather{}, fmt.Errorf("decoding response: %w", err)
	}
	weather, err := r.toWeather()
	if err != nil {
		return core.Weather{}, err
	}
	return weather, nil
}
