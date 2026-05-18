package meteo

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

type bouyObservation struct {
	BuoyID                int
	RecordedAt            time.Time
	WindDirectionDegT     *float64
	WindSpeedMetersPerSec *float64
	WindGustMetersPerSec  *float64
	WaveHeightM           *float64
	DominantWavePeriodSec *float64
	AvgWavePeriodSec      *float64
	MeanWaveDirectionDegT *float64
	AirTempDegC           *float64
	WaterTempDegC         *float64
	InsertedAt            time.Time
}

// getData wraps get() for the client.RTBouyService type.
func (s *RTBouyService) getData(ctx context.Context, bouyId string) ([]byte, error) {
	return s.get(ctx, bouyId)
}

// parseBouyObservation takes raw data as a byte slice that
// is returned by a get() method and parses the data into a
// BouyObservation struct. It returns a pointer to a bouyObservation struct and an error.
func (s *RTBouyService) parseBuoyObservation(data []byte, bouyId string) (*bouyObservation, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var obs *bouyObservation

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)

		// Format time.
		timeLayout := "2006 01 02 15 04"
		timestamp, err := time.Parse(timeLayout, strings.Join(fields[:5], " "))
		if err != nil {
			return &bouyObservation{}, err
		}

		// bouyId convert to string
		id, err := strconv.Atoi(bouyId)
		if err != nil {
			return &bouyObservation{}, err
		}

		// safely parse datatypes
		windDirection, _ := parseDataFloat(fields[5])
		windSpeed, _ := parseDataFloat(fields[6])
		windGust, _ := parseDataFloat(fields[7])
		waveHeightM, _ := parseDataFloat(fields[8])
		dominantWavePeriod, _ := parseDataFloat(fields[9])
		avgWavePeriod, _ := parseDataFloat(fields[10])
		meanWaveDirection, _ := parseDataFloat(fields[11])
		airTemperature, _ := parseDataFloat(fields[13])
		waterTemperature, _ := parseDataFloat(fields[14])

		// build BouyObservation
		obs = &bouyObservation{
			BuoyID:                id,
			RecordedAt:            timestamp,
			WindDirectionDegT:     windDirection,
			WindSpeedMetersPerSec: windSpeed,
			WindGustMetersPerSec:  windGust,
			WaveHeightM:           waveHeightM,
			DominantWavePeriodSec: dominantWavePeriod,
			AvgWavePeriodSec:      avgWavePeriod,
			MeanWaveDirectionDegT: meanWaveDirection,
			AirTempDegC:           airTemperature,
			WaterTempDegC:         waterTemperature,
			InsertedAt:            time.Now().UTC(),
		}
		break
	}
	return obs, nil
}

// GetObservation takes context.Context and a string.
// It returns parsed JSON of the bouy observation in the format
// acceptable to the databse.
func (s *RTBouyService) GetObservation(ctx context.Context, bouyId string) (*bouyObservation, error) {
	data, err := s.getData(ctx, bouyId)
	if err != nil {
		return &bouyObservation{}, err
	}
	obs, err := s.parseBuoyObservation(data, bouyId)
	if err != nil {
		return &bouyObservation{}, err
	}
	return obs, nil
}

type weatherObservation struct {
	Properties properties `json:"properties"`
}

type properties struct {
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

// getData wraps get() for the client.RTWeatherService type.
// Returns a slice of raw data and an error.
func (s *RTWeatherService) getData(ctx context.Context, stationId string) ([]byte, error) {
	return s.get(ctx, stationId)
}

// parseWeatherObservation takes raw data as a byte slice that
// is returned by a get() method and parses the json into a
// WeatherObservation struct. It returns a WeatherObservation struct and an error.
func parseWeatherObservation(data []byte) (weatherObservation, error) {
	var obs weatherObservation
	if err := json.Unmarshal(data, &obs); err != nil {
		return weatherObservation{}, err
	}
	return obs, nil
}

// GetObservation takes context.Context and a string.
// It returns parsed JSON of the weather observation in the format
// acceptable to the databse.
func (s *RTWeatherService) GetObservation(ctx context.Context, stationId string) (weatherObservation, error) {
	data, err := s.getData(ctx, stationId)
	if err != nil {
		return weatherObservation{}, err
	}
	obs, err := parseWeatherObservation(data)
	if err != nil {
		return weatherObservation{}, err
	}
	return obs, err
}

// UTILITY FUNCTIONS
// parseDataFloat returns pointer. Enables data point to return multiple states.
// Allows for a number, a missing value ("MM") or nil which are all
// valid values for the dataset.
func parseDataFloat(value string) (*float64, error) {
	if value == "MM" {
		return nil, nil
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
