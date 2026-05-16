# 🚀 Full Cycle Challenge: Distributed Tracing with OTEL & Zipkin

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![OpenTelemetry](https://img.shields.io/badge/OTEL-00ADD8?style=flat&logo=opentelemetry)](https://opentelemetry.io/)
[![Zipkin](https://img.shields.io/badge/Zipkin-FF6600?style=flat&logo=zipkin)](https://zipkin.io/)
[![Status](https://img.shields.io/badge/Status-Completed-success?style=flat)]()

This project implements a distributed system in Go composed of two microservices (**Service A** and **Service B**) that cooperate to query the weather of a city based on a Brazilian ZIP code (CEP). It features **Distributed Tracing** using **OpenTelemetry (OTEL)** and **Zipkin**.

## 🧠 Architecture & Logic

The system follows a distributed orchestration pattern:

1.  **Service A (Input)**:
    - Receives a POST request with a CEP.
    - Validates if the CEP is a string with exactly 8 digits.
    - If valid, forwards the request to Service B.
2.  **Service B (Orchestration)**:
    - Receives the CEP from Service A.
    - Consults **ViaCEP** to identify the city name (Manual Span).
    - Consults **WeatherAPI** to get the current temperature (Manual Span).
    - Returns temperatures in Celsius, Fahrenheit, and Kelvin.
3.  **Observability**:
    - Both services are instrumented with OpenTelemetry.
    - An **OTEL Collector** receives the spans and exports them to **Zipkin**.

---

## 📁 Project Structure

```text
.
├── service-a/           # Input Service
│   ├── cmd/             # Entry point
│   ├── internal/        # Logic & OTEL setup
│   └── Dockerfile       
├── service-b/           # Orchestration Service
│   ├── cmd/             # Entry point
│   ├── internal/        # Logic, Use Cases & External Clients
│   └── Dockerfile       
├── docker-compose.yaml  # Infrastructure orchestration
├── otel-collector-config.yaml # OTEL Collector configuration
└── README.md            # Documentation
```

---

## ⚙️ Configuration

The only mandatory configuration is the **WeatherAPI** key for **Service B**. You can get one for free at [weatherapi.com](https://www.weatherapi.com/).

### How to configure:
1. Open `docker-compose.yaml`.
2. Locate the `service-b` section.
3. Replace `your_api_key_here` with your actual API key.

| Variable | Description | Default/Status |
|----------|-------------|----------------|
| `WEATHER_API_KEY` | API Key for WeatherAPI | **Mandatory** |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP Collector Address | Pre-configured (`otel-collector:4317`) |

---

## 🚀 How to Run

### 1. Start the System
```bash
docker compose up --build
```

This will start:
- **Service A** on `http://localhost:8080`
- **Service B** on `http://localhost:8081`
- **OTEL Collector** on `4317`
- **Zipkin** on `http://localhost:9411`

---

## 📊 API Usage & Testing

You can use the following `curl` commands to test the system:

### 1. Success Scenario (Valid CEP)
**Request:**
```bash
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{"cep": "01153000"}'
```
**Expected Response (200 OK):**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

### 2. Invalid Format (Validation Error)
**Request:**
```bash
curl -i --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{"cep": "123"}'
```
**Expected Response (422 Unprocessable Entity):**
```text
invalid zipcode
```

### 3. Zip Code Not Found
**Request:**
```bash
curl -i --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{"cep": "99999999"}'
```
**Expected Response (404 Not Found):**
```text
can not find zipcode
```

---

## 🔍 Visualizing Traces in Zipkin

## 🛠️ Technologies
- **Go** (Golang)
- **OpenTelemetry** (SDK & OTLP)
- **Zipkin**
- **OTEL Collector**
- **Docker & Docker Compose**
- **ViaCEP API**
- **WeatherAPI**
