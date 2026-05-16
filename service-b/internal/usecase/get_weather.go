package usecase

import (
	"context"

	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/entity"
	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/infra/viacep"
	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/infra/weatherapi"
)

type WeatherOutputDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherUseCase struct{}

func NewWeatherUseCase() *WeatherUseCase {
	return &WeatherUseCase{}
}

func (u *WeatherUseCase) Execute(ctx context.Context, zipcode, apiKey string) (*WeatherOutputDTO, error) {
	err := entity.ValidateZipCode(zipcode)
	if err != nil {
		return nil, err
	}

	city, err := viacep.GetCityByZipCode(ctx, zipcode)
	if err != nil {
		if err.Error() == "zipcode not found" {
			return nil, entity.ErrZipCodeNotFound
		}
		return nil, err
	}

	tempC, err := weatherapi.GetTemperatureByCity(ctx, city, apiKey)
	if err != nil {
		return nil, err
	}

	weather := entity.NewWeather(tempC)

	return &WeatherOutputDTO{
		City:  city,
		TempC: weather.TempC,
		TempF: weather.TempF,
		TempK: weather.TempK,
	}, nil
}
