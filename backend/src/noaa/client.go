package hydromet

import (
	"context"
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

func (s *RTBouyService) getData(ctx context.Context, bouyId string) ([]byte, error) {
	return s.get(ctx, bouyId)
}

func (s *RTWeatherService) getData(ctx context.Context, stationId string) ([]byte, error) {
	return s.get(ctx, stationId)
}
