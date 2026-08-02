package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

// Parse and validate the JSON body
func parseBody(r *http.Request) ([]float64, error) {
	var numbers []float64
	err := json.NewDecoder(r.Body).Decode(&numbers)
	if err != nil {
		return nil, err
	}
	return numbers, nil
}

// Start: Calculations
func roundThree(val float64) float64 {
	return math.Round(val*1000) / 1000
}

func calcMean(numbers []float64) float64 {
	var sum float64
	for _, n := range numbers {
		sum += n
	}
	return sum / float64(len(numbers))
}

func calcStddev(numbers []float64, mean float64) float64 {
	var sumSquares float64
	for _, n := range numbers {
		sumSquares += (n - mean) * (n - mean)
	}
	stddev := math.Sqrt(sumSquares / float64(len(numbers)))
	return stddev
}
// End: Calculations

// Start: Handlers
func meanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	numbers, err := parseBody(r)
	if err != nil || len(numbers) == 0 {
		http.Error(w, "Malformed request: strictly requires a JSON array of numbers", http.StatusBadRequest)
		return
	}

	rounded := roundThree(calcMean(numbers))
	
	// Format float outputs to drop trailing zeros (e.g., 3 instead of 3.000)
	w.Write([]byte(strconv.FormatFloat(rounded, 'f', -1, 64) + "\n"))
}

func stddevHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	numbers, err := parseBody(r)
	if err != nil || len(numbers) == 0 {
		http.Error(w, "Malformed request: strictly requires a JSON array of numbers", http.StatusBadRequest)
		return
	}

	// Calculate Population Standard Deviation
	rounded := roundThree(calcStddev(numbers, calcMean(numbers)))
	
	// Format float outputs to drop trailing zeros 
	w.Write([]byte(strconv.FormatFloat(rounded, 'f', -1, 64) + "\n"))
}
// End: Handlers

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Fallback 
	}

	http.HandleFunc("/mean", meanHandler)
	http.HandleFunc("/stddev", stddevHandler)

	log.Printf("Starting HTTP web server - port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}