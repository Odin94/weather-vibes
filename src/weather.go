package weathervibes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	weatherTodayUrl = "https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code,rain,showers,snowfall,uv_index,wind_speed_10m"
	weatherWeekUrl  = "https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&daily=weather_code,temperature_2m_max,temperature_2m_min,rain_sum,showers_sum,snowfall_sum,uv_index_max,wind_speed_10m_max,precipitation_probability_mean"
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

type WeeklyWeatherResponse struct {
	Lat                  float64      `json:"latitude"`
	Lon                  float64      `json:"longitude"`
	GenerationTimeMs     float64      `json:"generationtime_ms"`
	UtcOffsetSeconds     int          `json:"utc_offset_seconds"`
	Timezone             string       `json:"timezone"`
	TimezoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float64      `json:"elevation"`
	DailyUnits           DailyUnits   `json:"daily_units"`
	Daily                DailyWeather `json:"daily"`
}

type BaseWeatherUnits struct {
	Time         string `json:"time"`
	WeatherCode  string `json:"weather_code"`
	RainSum      string `json:"rain_sum"`
	ShowersSum   string `json:"showers_sum"`
	SnowfallSum  string `json:"snowfall_sum"`
	UvIndexMax   string `json:"uv_index_max"`
	WindSpeed10mMax string `json:"wind_speed_10m_max"`
}

type CurrentUnits struct {
	Time               string `json:"time"`
	Interval           string `json:"interval"`
	Temperature2m      string `json:"temperature_2m"`
	RelativeHumidity2m string `json:"relative_humidity_2m"`
	WeatherCode        string `json:"weather_code"`
	Rain               string `json:"rain"`
	Showers            string `json:"showers"`
	Snowfall           string `json:"snowfall"`
	UvIndex            string `json:"uv_index"`
	WindSpeed10m       string `json:"wind_speed_10m"`
}

type DailyUnits struct {
	BaseWeatherUnits
	Temperature2mMax            string `json:"temperature_2m_max"`
	Temperature2mMin            string `json:"temperature_2m_min"`
	PrecipitationProbabilityMean string `json:"precipitation_probability_mean"`
}

type BaseWeatherData struct {
	Time         string  `json:"time"`
	WeatherCode  int     `json:"weather_code"`
	Rain         float64 `json:"rain"`
	Showers      float64 `json:"showers"`
	Snowfall     float64 `json:"snowfall"`
	UvIndex      float64 `json:"uv_index"`
	WindSpeed10m float64 `json:"wind_speed_10m"`
}

type CurrentWeather struct {
	BaseWeatherData
	Interval           int     `json:"interval"`
	Temperature2m      float64 `json:"temperature_2m"`
	RelativeHumidity2m int     `json:"relative_humidity_2m"`
}

// Daily weather data (arrays for multiple days)
type DailyWeather struct {
	Time                        []string  `json:"time"`
	WeatherCode                 []int     `json:"weather_code"`
	Temperature2mMax            []float64 `json:"temperature_2m_max"`
	Temperature2mMin            []float64 `json:"temperature_2m_min"`
	RainSum                     []float64 `json:"rain_sum"`
	ShowersSum                  []float64 `json:"showers_sum"`
	SnowfallSum                 []float64 `json:"snowfall_sum"`
	UvIndexMax                  []float64 `json:"uv_index_max"`
	WindSpeed10mMax             []float64 `json:"wind_speed_10m_max"`
	PrecipitationProbabilityMean []int     `json:"precipitation_probability_mean"`
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

func GetWeeklyWeather(lat float64, lon float64) (WeeklyWeatherResponse, error) {
	var weeklyResponse WeeklyWeatherResponse

	resp, err := http.Get(fmt.Sprintf(weatherWeekUrl, lat, lon))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weeklyResponse, fmt.Errorf("failed to read weekly weather response body: %w", err)
	}

	if resp.StatusCode > 299 {
		return weeklyResponse, fmt.Errorf("weekly weather response failed with status code: %d and body: %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &weeklyResponse)
	if err != nil {
		return weeklyResponse, fmt.Errorf("failed to unmarshal weekly weather json: %w", err)
	}

	return weeklyResponse, nil
}

func GetWeatherDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code >= 1 && code <= 3:
		return "Partly cloudy"
	case code == 45 || code == 48:
		return "Fog"
	case code >= 51 && code <= 55:
		return "Drizzle"
	case code == 56 || code == 57:
		return "Freezing Drizzle"
	case code >= 61 && code <= 65:
		return "Rain"
	case code == 66 || code == 67:
		return "Freezing Rain"
	case code >= 71 && code <= 75:
		return "Snow fall"
	case code == 77:
		return "Snow grains"
	case code >= 80 && code <= 82:
		return "Rain showers"
	case code == 85 || code == 86:
		return "Snow showers"
	case code == 95:
		return "Thunderstorm"
	case code == 96 || code == 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown weather condition"
	}
}

func GetWeatherEmoji(description string) string {
	switch description {
	case "Clear sky":
		return "☀️"
	case "Partly cloudy":
		return "⛅"
	case "Fog":
		return "🌫️"
	case "Drizzle":
		return "🌦️"
	case "Freezing Drizzle":
		return "🌨️"
	case "Rain":
		return "🌧️"
	case "Freezing Rain":
		return "🌨️"
	case "Snow fall":
		return "❄️"
	case "Snow grains":
		return "🌨️"
	case "Rain showers":
		return "🌦️"
	case "Snow showers":
		return "🌨️"
	case "Thunderstorm":
		return "⛈️"
	case "Thunderstorm with hail":
		return "⛈️"
	default:
		return "🌤️"
	}
}
