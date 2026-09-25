package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serve(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, req)
	return rec
}

func TestCalculate(t *testing.T) {
	const (
		invalidBody    = `{"error":"request body must be a JSON object with operation, a and b"}`
		multipleValues = `{"error":"request body must contain a single JSON object"}`
		missingOperand = `{"error":"a and b are required"}`
		sqrtOperands   = `{"error":"sqrt takes only a"}`
		unknownOp      = `{"error":"unknown operation"}`
	)

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		// success: one per operation, to check the JSON names map to the domain
		{name: "add", body: `{"operation":"add","a":2,"b":3}`, wantStatus: http.StatusOK, wantBody: `{"result":5}`},
		{name: "subtract", body: `{"operation":"subtract","a":10,"b":4}`, wantStatus: http.StatusOK, wantBody: `{"result":6}`},
		{name: "multiply", body: `{"operation":"multiply","a":3,"b":4}`, wantStatus: http.StatusOK, wantBody: `{"result":12}`},
		{name: "divide", body: `{"operation":"divide","a":10,"b":4}`, wantStatus: http.StatusOK, wantBody: `{"result":2.5}`},
		{name: "power", body: `{"operation":"power","a":2,"b":10}`, wantStatus: http.StatusOK, wantBody: `{"result":1024}`},
		{name: "sqrt takes only a", body: `{"operation":"sqrt","a":16}`, wantStatus: http.StatusOK, wantBody: `{"result":4}`},
		{name: "percentage", body: `{"operation":"percentage","a":10,"b":200}`, wantStatus: http.StatusOK, wantBody: `{"result":20}`},
		{name: "percentage requires b", body: `{"operation":"percentage","a":10}`, wantStatus: http.StatusBadRequest, wantBody: missingOperand},
		{name: "zero is a valid operand", body: `{"operation":"add","a":0,"b":0}`, wantStatus: http.StatusOK, wantBody: `{"result":0}`},
		{name: "trailing whitespace is allowed", body: "{\"operation\":\"add\",\"a\":2,\"b\":3}\n", wantStatus: http.StatusOK, wantBody: `{"result":5}`},

		// malformed request body
		{name: "empty body", body: ``, wantStatus: http.StatusBadRequest, wantBody: invalidBody},
		{name: "malformed JSON", body: `{"operation":`, wantStatus: http.StatusBadRequest, wantBody: invalidBody},
		{name: "operand of wrong type", body: `{"operation":"add","a":"2","b":3}`, wantStatus: http.StatusBadRequest, wantBody: invalidBody},
		{name: "unknown field", body: `{"operation":"add","a":2,"b":3,"c":4}`, wantStatus: http.StatusBadRequest, wantBody: invalidBody},
		{name: "number outside float64 range", body: `{"operation":"add","a":1e999,"b":1}`, wantStatus: http.StatusBadRequest, wantBody: invalidBody},
		{name: "second JSON object", body: `{"operation":"add","a":2,"b":3}{}`, wantStatus: http.StatusBadRequest, wantBody: multipleValues},
		{name: "trailing garbage", body: `{"operation":"add","a":2,"b":3} x`, wantStatus: http.StatusBadRequest, wantBody: multipleValues},

		// missing or invalid fields
		{name: "missing b", body: `{"operation":"add","a":2}`, wantStatus: http.StatusBadRequest, wantBody: missingOperand},
		{name: "null a", body: `{"operation":"add","a":null,"b":3}`, wantStatus: http.StatusBadRequest, wantBody: missingOperand},
		{name: "sqrt with b", body: `{"operation":"sqrt","a":16,"b":2}`, wantStatus: http.StatusBadRequest, wantBody: sqrtOperands},
		{name: "sqrt without a", body: `{"operation":"sqrt"}`, wantStatus: http.StatusBadRequest, wantBody: sqrtOperands},
		{name: "missing operation", body: `{"a":2,"b":3}`, wantStatus: http.StatusBadRequest, wantBody: unknownOp},
		{name: "unknown operation", body: `{"operation":"modulo","a":2,"b":3}`, wantStatus: http.StatusBadRequest, wantBody: unknownOp},

		// valid request, but the math has no answer
		{name: "division by zero", body: `{"operation":"divide","a":1,"b":0}`, wantStatus: http.StatusUnprocessableEntity, wantBody: `{"error":"cannot divide by zero"}`},
		{name: "result out of range", body: `{"operation":"multiply","a":1e308,"b":10}`, wantStatus: http.StatusUnprocessableEntity, wantBody: `{"error":"result is out of range"}`},
		{name: "sqrt of negative", body: `{"operation":"sqrt","a":-4}`, wantStatus: http.StatusUnprocessableEntity, wantBody: `{"error":"result is not a real number"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(t, http.MethodPost, "/api/v1/calculate", tt.body)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json")
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
		})
	}
}

func TestCalculateRejectsWrongMethod(t *testing.T) {
	rec := serve(t, http.MethodGet, "/api/v1/calculate", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHealth(t *testing.T) {
	rec := serve(t, http.MethodGet, "/health", "")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got, want := strings.TrimSpace(rec.Body.String()), `{"status":"ok"}`; got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}
