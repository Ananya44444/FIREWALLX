package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"acf-demo/sidecar/pipeline"
)

type evaluateRequest struct {
	Input string `json:"input"`
}

type evaluateResponse struct {
	Decision string   `json:"decision"`
	Score    float64  `json:"score"`
	Signals  []string `json:"signals"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	patternsPath, err := resolvePatternsPath()
	if err != nil {
		log.Fatalf("failed to resolve patterns file: %v", err)
	}

	scanner, err := pipeline.NewScanner(patternsPath)
	if err != nil {
		log.Fatalf("failed to load patterns: %v", err)
	}

	http.HandleFunc("/evaluate", evaluateHandler(scanner))

	log.Println("ACF sidecar listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func evaluateHandler(scanner *pipeline.Scanner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
			return
		}

		var req evaluateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON payload"})
			return
		}

		if err := pipeline.ValidateInput(req.Input); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		normalized := pipeline.NormalizeIterative(req.Input, 3)
		signals := scanner.Scan(normalized)
		score := pipeline.AggregateScore(signals)
		decision := pipeline.Decide(score)

		resp := evaluateResponse{
			Decision: string(decision),
			Score:    score,
			Signals:  signals,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func resolvePatternsPath() (string, error) {
	candidates := []string{
		"sidecar/patterns.json",
		"patterns.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", errors.New("patterns.json not found")
}
