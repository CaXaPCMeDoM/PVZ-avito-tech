package entity

import "errors"

type City string

const (
	CityMoscow City = "Москва"
	CitySpb    City = "Санкт-Петербург"
	CityKazan  City = "Казань"
)

var validCityMap = map[City]struct{}{
	CityMoscow: {},
	CitySpb:    {},
	CityKazan:  {},
}

func (r City) IsValidCity() bool {
	_, exists := validCityMap[r]
	return exists
}

func (r City) ValidateCity() error {
	if !r.IsValidCity() {
		return errors.New("invalid city")
	}
	return nil
}
