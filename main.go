package main

import (
    "encoding/json"
    "errors"
    "log"
    "math"
    "net/http"
    "os"
    "strconv"
)

const maxBodySize = 1 << 20 // 1 MiB

// Parse and validate the JSON body
func parseBody(req *http.Request) ([]float64, error) {
    var raw []interface{}
    if err := json.NewDecoder(req.Body).Decode(&raw); err != nil {
        return nil, errors.New("Malformed request: request body must be a JSON array")
    }

    var numbers []float64
    for _, item := range raw {
        switch temp := item.(type) {
        case float64:
            numbers = append(numbers, temp)
        case string:
            if strNum, err := strconv.ParseFloat(temp, 64); err == nil {
                numbers = append(numbers, strNum)
            }
        }
    }

    if len(numbers) == 0 {
        return nil, errors.New("Malformed request: request body must contain at least one numeric value")
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
    return math.Sqrt(sumSquares / float64(len(numbers)))
}
// End: Calculations

// Start: Write Helpers
func writeError(writer http.ResponseWriter, status int, message string) {
    http.Error(writer, message, status)
}

func writeFloat(writer http.ResponseWriter, value float64) {
    writer.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64) + "\n"))
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

    writeFloat(writer, roundThree(calcMean(numbers)))
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

    writeFloat(writer, roundThree(calcStddev(numbers, calcMean(numbers))))
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