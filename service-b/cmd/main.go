package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/infra/web"
	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/otel"
	"github.com/luigolima/go-otel-distributed-tracing-challenge/service-b/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "service-b"
	}
	collectorAddr := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorAddr == "" {
		collectorAddr = "otel-collector:4317"
	}
	apiKey := os.Getenv("WEATHER_API_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	shutdown, err := otel.InitTracer(serviceName, collectorAddr)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	uc := usecase.NewWeatherUseCase()
	handler := web.NewWeatherHandler(uc, apiKey)

	// Instrument handler with OTEL
	otelHandler := otelhttp.NewHandler(handler, "weather-handler")

	mux := http.NewServeMux()
	mux.Handle("/", otelHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Service B starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down Service B...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
