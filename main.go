package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"
)

const maxBodySize = 1 << 20 // 1 MiB

// Start: Types
type HealthResponse struct {
	Status string `json:"status"`
}

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

type ResultResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// End: Types

// Intercept the WriteHeader method to capture the status code
func (writer *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	writer.statusCode = statusCode
	writer.ResponseWriter.WriteHeader(statusCode)
}

// Logging wrapper middleware to log incoming requests and their execution time
func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()

			// Get Status Code from ResponseWriterInterceptor
			interceptor := &ResponseWriterInterceptor{
				ResponseWriter: writer,
				statusCode:     http.StatusOK,
			}

			// Continue with request
			next.ServeHTTP(interceptor, request)

			// Log execution details
			extraParams := []slog.Attr{
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Int("status", interceptor.statusCode),
				slog.Duration("duration", time.Since(start)),
				slog.String("ip", request.RemoteAddr),
			}

			msg := "Request failed with unexpected status code"
			level := slog.LevelError
			switch interceptor.statusCode {
			case http.StatusOK:
				msg = "Completed HTTP Request"
				level = slog.LevelInfo
			case http.StatusBadRequest:
				msg = "Request failed with bad request"
				level = slog.LevelWarn
			case http.StatusMethodNotAllowed:
				msg = "Request failed with method not allowed"
				level = slog.LevelWarn
			}

			slog.LogAttrs(request.Context(), level, msg,
				extraParams...,
			)
		})
	}
}

// Parse and validate the JSON body
func parseBody(req *http.Request) ([]float64, error) {
	defer req.Body.Close()

	var numbers []float64

	if err := json.NewDecoder(req.Body).Decode(&numbers); err != nil {
		return nil, errors.New("Request body must be a JSON array of numbers")
	}

	if len(numbers) == 0 {
		return nil, errors.New("Request array must contain at least one Numeric Value")
	}

	return numbers, nil
}

// Start: Calculations
func roundThree(val float64) float64 {
	return math.Round(val*1000) / 1000
}

func calcMean(numbers []float64) float64 {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func calcStddev(numbers []float64, mean float64) float64 {
	var sumSquares float64
	for _, n := range numbers {
		sumSquares += (n - mean) * (n - mean)
	}
	return math.Sqrt(sumSquares / float64(len(numbers)))
}

// End: Calculations

// Start: Write Helpers
func writeJSON(writer http.ResponseWriter, status int, payload interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(payload)
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, ErrorResponse{Error: message})
}

// End: Write Helpers

// Start: Handlers
func healthHandler(writer http.ResponseWriter, req *http.Request) {
	writeJSON(writer, http.StatusOK, HealthResponse{Status: "OK"})
}

func meanHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		writeError(writer, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}

	req.Body = http.MaxBytesReader(writer, req.Body, maxBodySize)

	numbers, err := parseBody(req)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	mean := roundThree(calcMean(numbers))

	slog.Debug("Calculating mean",
		slog.Any("numbers", numbers),
		slog.Float64("Result", mean),
	)

	writeJSON(writer, http.StatusOK, ResultResponse{Result: mean})
}

func stddevHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		writeError(writer, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}

	req.Body = http.MaxBytesReader(writer, req.Body, maxBodySize)

	numbers, err := parseBody(req)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	mean := calcMean(numbers)
	stddev := roundThree(calcStddev(numbers, mean))

	slog.Debug("Calculating standard deviation",
		slog.Any("numbers", numbers),
		slog.Float64("mean", mean),
		slog.Float64("Result", stddev),
	)

	writeJSON(writer, http.StatusOK, ResultResponse{Result: stddev})
}

// End: Handlers

func main() {
	logLevel := new(slog.LevelVar)

	switch os.Getenv("LOG_LEVEL") {
	case "ERROR":
		logLevel.Set(slog.LevelError)
	case "WARN":
		logLevel.Set(slog.LevelWarn)
	case "DEBUG":
		logLevel.Set(slog.LevelDebug)
	default:
		logLevel.Set(slog.LevelInfo)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(log)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Fallback
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/mean", meanHandler)
	mux.HandleFunc("/stddev", stddevHandler)

	wrappedMux := Logger(log)(mux)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      wrappedMux,
		ReadTimeout:  5e9, // 5s
		WriteTimeout: 5e9, // 5s
	}

	slog.Info("Running Math Service with logging", slog.String("level", logLevel.Level().String()))
	slog.Info("Starting HTTP web server", slog.String("port", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed to start", slog.Any("error", err))
	}

}
