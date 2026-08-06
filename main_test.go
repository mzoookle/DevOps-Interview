package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoundThree(test *testing.T) {
	got := roundThree(1.23456)
	want := 1.235
	if got != want {
		test.Fatalf("roundThree(1.23456) = %v, want %v", got, want)
	}
}

func TestCalcMean(test *testing.T) {
	got := calcMean([]float64{1, 2, 3, 4})
	want := 2.5
	if got != want {
		test.Fatalf("calcMean = %v, want %v", got, want)
	}
}

func TestCalcStddev(test *testing.T) {
	mean := calcMean([]float64{1, 2, 3, 4})
	got := calcStddev([]float64{1, 2, 3, 4}, mean)
	want := 1.118033988749895
	if got != want {
		test.Fatalf("calcStddev = %v, want %v", got, want)
	}
}

func TestParseBodySuccess(test *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/mean", bytes.NewBufferString(`[1,2,3]`))
	req.Header.Set("Content-Type", "application/json")

	numbers, err := parseBody(req)
	if err != nil {
		test.Fatalf("parseBody returned error: %v", err)
	}

	want := []float64{1, 2, 3}
	for i := range want {
		if numbers[i] != want[i] {
			test.Fatalf("parseBody numbers[%d] = %v, want %v", i, numbers[i], want[i])
		}
	}
}

func TestParseBodyBadJSON(test *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/mean", bytes.NewBufferString(`["a", "b"]`))
	req.Header.Set("Content-Type", "application/json")

	_, err := parseBody(req)
	if err == nil {
		test.Fatal("expected parseBody to return error for invalid JSON numbers")
	}
}

func TestMeanHandler(test *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/mean", bytes.NewBufferString(`[1,2,3]`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	meanHandler(rec, req)

	if rec.Code != http.StatusOK {
		test.Fatalf("meanHandler status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp ResultResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		test.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 2 {
		test.Fatalf("meanHandler result = %v, want 2", resp.Result)
	}
}
