package core

import (
	"errors"
	"math"
)

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var KnownLocations = []Location{
	{Name: "Almaty", Latitude: 43.2567, Longitude: 76.9286},
	{Name: "Moscow", Latitude: 55.7522, Longitude: 37.6156},
	{Name: "SaintP", Latitude: 59.9386, Longitude: 30.3141},
}

func LocationFromCoordinates(lat, long float64) (Location, error) {
	roughlySame := func(f1, f2 float64) bool {
		return math.Abs(f1 - f2) < 0.1
	}
	for _, loc := range KnownLocations {
		// Should be smart math instead
		if roughlySame(lat, loc.Latitude) && roughlySame(long, loc.Longitude) {
			return loc, nil
		}
	}
	return Location{}, errors.New("unknown coordinates")
}
