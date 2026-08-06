package main

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"os"
)

const maxBodySize = 1 << 20 // 1 MiB

type ResultResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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

	writeJSON(writer, http.StatusOK, ResultResponse{Result: roundThree(calcMean(numbers))})
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
	stddev := calcStddev(numbers, mean)
	writeJSON(writer, http.StatusOK, ResultResponse{Result: roundThree(stddev)})
}

// End: Handlers

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Fallback
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mean", meanHandler)
	mux.HandleFunc("/stddev", stddevHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5e9, // 5s
		WriteTimeout: 5e9, // 5s
	}

	log.Printf("Starting HTTP web server - port %s...\n", port)
	log.Fatal(server.ListenAndServe())

}
