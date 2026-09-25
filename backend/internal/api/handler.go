// Package api exposes the calculator over HTTP as a JSON API.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"calculator-go/backend/internal/calculator"
)

// NewRouter returns the HTTP handler with all API routes registered.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/calculate", handleCalculate)
	mux.HandleFunc("GET /health", handleHealth)
	return mux
}

type calculateRequest struct {
	Operation string `json:"operation"`
	// Pointers distinguish a missing or null operand from a real 0.
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req calculateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be a JSON object with operation, a and b")
		return
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request body must contain a single JSON object")
		return
	}

	// Check the operation first, so the operand messages below can rely on it.
	op := calculator.Operation(req.Operation)
	switch {
	case !op.IsValid():
		writeError(w, http.StatusBadRequest, calculator.ErrUnknownOperation.Error())
		return
	case op.IsUnary() && req.A == nil:
		writeError(w, http.StatusBadRequest, "a is required")
		return
	case op.IsUnary() && req.B != nil:
		writeError(w, http.StatusBadRequest, string(op)+" takes only a")
		return
	case !op.IsUnary() && (req.A == nil || req.B == nil):
		writeError(w, http.StatusBadRequest, "a and b are required")
		return
	}

	var b float64
	if req.B != nil {
		b = *req.B
	}

	result, err := calculator.Calculate(op, *req.A, b)
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, calculateResponse{Result: result})
	case errors.Is(err, calculator.ErrDivisionByZero),
		errors.Is(err, calculator.ErrNotRealNumber),
		errors.Is(err, calculator.ErrResultOutOfRange):
		writeError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Printf("calculate: unexpected error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
