package weathervibes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	weatherTodayUrl = "https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code"
	weatherWeekUrl  = "https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&daily=weather_code,temperature_2m_max,temperature_2m_min&timezone=auto"
)

type WeatherResponse struct {
	Lat                  float64        `json:"latitude"`
	Lon                  float64        `json:"longitude"`
	GenerationTimeMs     float64        `json:"generationtime_ms"`
	UtcOffsetSeconds     int            `json:"utc_offset_seconds"`
	Timezone             string         `json:"timezone"`
	TimezoneAbbreviation string         `json:"timezone_abbreviation"`
	Elevation            float64        `json:"elevation"`
	CurrentUnits         CurrentUnits   `json:"current_units"`
	CurrentWeather       CurrentWeather `json:"current"`
}

type CurrentUnits struct {
	Time               string `json:"time"` // datetime
	Interval           string `json:"interval"`
	Temperature2m      string `json:"temperature_2m"`
	RelativeHumidity2m string `json:"relative_humidity_2m"`
	WeatherCode        string `json:"weather_code"`
}

type CurrentWeather struct {
	Time               string  `json:"time"` // datetime format, eg. "iso8601"
	Interval           int     `json:"interval"`
	Temperature2m      float64 `json:"temperature_2m"`
	RelativeHumidity2m int     `json:"relative_humidity_2m"`
	WeatherCode        int     `json:"weather_code"`
}

func GetWeather(lat float64, lon float64) (WeatherResponse, error) {
	var weatherResponse WeatherResponse

	resp, err := http.Get(fmt.Sprintf(weatherTodayUrl, lat, lon))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weatherResponse, fmt.Errorf("failed to read weather response body: %w", err)
	}

	if resp.StatusCode > 299 {
		return weatherResponse, fmt.Errorf("weather response failed with status code: %d and body: %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return weatherResponse, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return weatherResponse, nil
}
