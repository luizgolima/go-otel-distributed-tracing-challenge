package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type CepRequest struct {
	Cep string `json:"cep"`
}

type CepHandler struct {
	ServiceBURL string
}

func NewCepHandler(serviceBURL string) *CepHandler {
	return &CepHandler{ServiceBURL: serviceBURL}
}

func (h *CepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CepRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Validation
	matched, _ := regexp.MatchString(`^[0-9]{8}$`, req.Cep)
	if !matched {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Forward to Service B
	ctx := r.Context()
	tr := otel.Tracer("service-a")
	_, span := tr.Start(ctx, "forward-to-service-b")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	
	reqB, err := http.NewRequestWithContext(ctx, "GET", h.ServiceBURL+"/"+req.Cep, nil)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(reqB)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
