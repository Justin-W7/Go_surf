package processing

import (
	"testing"
)

func TestParseSpotWeather_ExtractsHourlyURL(t *testing.T) {

	// Fake NWS response (trimmed but realistic structure)
	rawJSON := []byte(`
{
  "properties": {
    "forecastHourly": "https://api.weather.gov/gridpoints/SGX/35,58/forecast/hourly"
  }
}
`)

	result, err := ParseSpotWeather(rawJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Properties.ForecastHourly == "" {
		t.Fatal("forecastHourly is empty (this caused your production bug)")
	}

	expected := "https://api.weather.gov/gridpoints/SGX/35,58/forecast/hourly"

	if result.Properties.ForecastHourly != expected {
		t.Fatalf(
			"unexpected forecastHourly.\nexpected: %s\ngot: %s",
			expected,
			result.Properties.ForecastHourly,
		)
	}
}

func TestParseSpotWeather_MissingHourlyURL(t *testing.T) {

	rawJSON := []byte(`
{
  "properties": {}
}
`)

	result, err := ParseSpotWeather(rawJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Properties.ForecastHourly != "" {
		t.Fatal("expected empty ForecastHourly for missing field")
	}
}
