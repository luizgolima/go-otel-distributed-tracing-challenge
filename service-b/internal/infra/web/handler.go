package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/entity"
	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/usecase"
)

type WeatherHandler struct {
	UseCase *usecase.WeatherUseCase
	ApiKey  string
}

func NewWeatherHandler(uc *usecase.WeatherUseCase, apiKey string) *WeatherHandler {
	return &WeatherHandler{UseCase: uc, ApiKey: apiKey}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract CEP from URL /88888888
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	zipcode := parts[len(parts)-1]

	output, err := h.UseCase.Execute(r.Context(), zipcode, h.ApiKey)
	if err != nil {
		if err == entity.ErrInvalidZipCode {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if err == entity.ErrZipCodeNotFound {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
