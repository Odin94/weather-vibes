package main

import (
	"fmt"
	"log"
	"github.com/Odin94/weather-vibes/src/apis"
)

func main() {
	location, err := apis.GetGeolocation()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("City: %s, %s\n", location.City, location.Country)
	fmt.Printf("Coordinates: Lat %f, Lon %f\n", location.Lat, location.Lon)

	weather, err := apis.GetWeather(location.Lat, location.Lon)
	if err != nil {
		log.Fatal(err)
	}

	// TODOdin: Add Rain, Snowfall, Showers, Precipitation
	temperature, unit := weather.CurrentWeather.Temperature2m, weather.CurrentUnits.Temperature2m

	// TODOdin: Add cool artsy bubble tea interface  (maybe visualize the weather with an animation for rainfall etc...?)
	fmt.Printf("\n\n-----------\n\n")
	fmt.Printf("Temperature: %.1f %s", temperature, unit)
	fmt.Printf("\n")
}
