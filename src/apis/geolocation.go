package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const geolocationUrl = "http://ip-api.com/json/"

type GeoResponse struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

func GetGeolocation() (GeoResponse, error) {
	var geoResponse GeoResponse

	resp, err := http.Get(geolocationUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return geoResponse, fmt.Errorf("failed to read geo response body: %w", err)
	}

	if resp.StatusCode > 299 {
		return geoResponse, fmt.Errorf("geo response failed with status code: %d and body: %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &geoResponse)
	if err != nil {
		return geoResponse, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return geoResponse, nil
}
