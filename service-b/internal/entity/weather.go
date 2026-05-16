package entity

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidZipCode = errors.New("invalid zipcode")
	ErrZipCodeNotFound = errors.New("can not find zipcode")
)

type Weather struct {
	TempC float64
	TempF float64
	TempK float64
}

func NewWeather(tempC float64) *Weather {
	return &Weather{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}
}

func ValidateZipCode(zipcode string) error {
	matched, _ := regexp.MatchString(`^[0-9]{8}$`, zipcode)
	if !matched {
		return ErrInvalidZipCode
	}
	return nil
}
