// Package api facilitates url requests
package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go_surf/models"
)

func FetchSpitcastSpots(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("FetchSpitcastSpot(): No response from request: ", err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d", response.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}
	return data, nil
}

func FetchSpitcastForecast(spots []models.SurfSpot) ([][]byte, error) {
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()
	locations := make([][]byte, 0)

	for i := range spots {
		url := fmt.Sprintf(SpitcastForecastURL, spots[i].SpotID, year, month, day)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("FetchSpitcastForecast response error: ", err)
			return nil, err
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("FetchSpitcastForecast response.Body error: ", err)
		}

		locations = append(locations, data)
	}
	return locations, nil
}

func FetchWeatherPoint(long float64, lat float64) ([]byte, error) {
	url := fmt.Sprintf(NWSWeatherURL, long, lat)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: api_cpient FetchWeatherPoint() response err: ", err)
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR: In api_client FetchWeatherPoint() response.Body error: ", err)
		return nil, err
	}
	return data, nil
}

func FetchWeatherForecast(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: in api_client FetchWeatherForecast() response err: ", err)
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR: In api_client FetchWeatherForecast() response.Body err: ", err)
		return nil, err
	}
	return data, nil
}

func FetchHourlyWeatherForecast(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: In api_client FetchHourlyWeatherForecast() response err: ", err)
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR: In api_client FetchHourlyWeatherForecast() reslonse.Body error: ", err)
		return nil, err
	}

	return data, nil
}

func FetchWeatherGridForecast(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: In api_client FetchWeatherGridForecast() response err: ", err)
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR: In api_client FetchWeatherGridForecast() reslonse.Body error: ", err)
		return nil, err
	}

	return data, nil
}
