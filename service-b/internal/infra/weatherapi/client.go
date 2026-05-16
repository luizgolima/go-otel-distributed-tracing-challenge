package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
)

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetTemperatureByCity(ctx context.Context, city string, apiKey string) (float64, error) {
	tr := otel.Tracer("service-b")
	ctx, span := tr.Start(ctx, "get-temperature-by-city-weatherapi")
	defer span.End()

	cityEscaped := url.QueryEscape(city)
	resp, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + cityEscaped + "&aqi=no")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi returned status %d", resp.StatusCode)
	}

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return 0, err
	}

	return weather.Current.TempC, nil
}
