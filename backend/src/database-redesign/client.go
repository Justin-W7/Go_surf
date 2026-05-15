package meteo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// ndbcBouyDataUrl to access real time bouy data from NOAA.
	rtNDBCBouyDataURL = "https://www.ndbc.noaa.gov/data/realtime2/%s.txt"
	// rtWeatherUrl to access real-time weather data from weather.gov
	rtWeatherURL = "https://api.weather.gov/stations/%s/observations/latest"
)

type Client struct {
	httpClient *http.Client

	RTBouy    *RTBouyService
	RTWeather *RTWeatherService
}

type service struct {
	client  *Client
	baseURL string
}

type RTBouyService struct {
	*service
}

type RTWeatherService struct {
	*service
}

// NewClient returns a new API client.
func NewClient() *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	c.RTBouy = &RTBouyService{
		service: &service{
			client:  c,
			baseURL: rtNDBCBouyDataURL,
		},
	}
	c.RTWeather = &RTWeatherService{
		service: &service{
			client:  c,
			baseURL: rtWeatherURL,
		},
	}
	return c
}

// get takes context and an id (either a stationId or a bouyId- this may be expanded upon later).
// Sends an httml request to the designated baseURL.
// get returns the raw data of the html request in a slice of bytes - type []byte.
func (s *service) get(ctx context.Context, id string) ([]byte, error) {
	url := fmt.Sprintf(s.baseURL, id)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// getData wraps get() for the client.RTBouyService type.
func (s *RTBouyService) getData(ctx context.Context, bouyId string) ([]byte, error) {
	return s.get(ctx, bouyId)
}

type WeatherObservation struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	Timestamp            string       `josn:"timestamp"`
	Temperature          Value        `json:"temperature"`
	WindSpeed            Value        `json:"windSpeed"`
	WindDirection        Value        `json:"windDirection"`
	PrecipitaionLastHour Value        `json:"precipitationLastHour"`
	CloudLayers          []CloudLayer `json:"cloudLaters"`
}

type Value struct {
	Value *float64 `json:"value"`
}

type CloudLayer struct {
	Amount string `json:"amount"`
}

// parseWeatherObservation takes raw data as a byte slice that
// is returned by a get() method and parses the json into a
// WeatherObservation struct. It returns a WeatherObservation struct and an error.
func parseWeatherObservation(data []byte) (WeatherObservation, error) {
	var observation WeatherObservation
	if err := json.Unmarshal(data, &observation); err != nil {
		return WeatherObservation{}, err
	}
	return observation, nil
}

// getData wraps get() for the client.RTWeatherService type.
// Returns a slice of raw data and an error.
func (s *RTWeatherService) getData(ctx context.Context, stationId string) ([]byte, error) {
	return s.get(ctx, stationId)
}

// GetObservation takes context.Context and a string.
// It returns parsed JSON of the weather observation in the format
// acceptable to the databse.
func (s *RTWeatherService) GetObservation(ctx context.Context, stationId string) (WeatherObservation, error) {
	data, err := s.getData(ctx, stationId)
	if err != nil {
		return WeatherObservation{}, err
	}
	observation, err := parseWeatherObservation(data)
	if err != nil {
		return WeatherObservation{}, err
	}
	return observation, err
}
