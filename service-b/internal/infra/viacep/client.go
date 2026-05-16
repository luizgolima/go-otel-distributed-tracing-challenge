package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

type ViaCepResponse struct {
	Localidade string      `json:"localidade"`
	Erro       interface{} `json:"erro"`
}

func GetCityByZipCode(ctx context.Context, zipcode string) (string, error) {
	tr := otel.Tracer("service-b")
	ctx, span := tr.Start(ctx, "get-city-by-zipcode-viacep")
	defer span.End()

	resp, err := http.Get("https://viacep.com.br/ws/" + zipcode + "/json/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep returned status %d", resp.StatusCode)
	}

	var viacep ViaCepResponse
	err = json.NewDecoder(resp.Body).Decode(&viacep)
	if err != nil {
		return "", err
	}

	if viacep.Erro != nil {
		if errBool, ok := viacep.Erro.(bool); ok && errBool {
			return "", fmt.Errorf("zipcode not found")
		}
		if errStr, ok := viacep.Erro.(string); ok && errStr == "true" {
			return "", fmt.Errorf("zipcode not found")
		}
	}

	if viacep.Localidade == "" {
		return "", fmt.Errorf("zipcode not found")
	}

	return viacep.Localidade, nil
}
