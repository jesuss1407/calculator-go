package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewHandlerServesFrontendAndAPI(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("<h1>Calculator</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := newHandler(staticDir)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{name: "frontend index", method: http.MethodGet, path: "/", wantStatus: http.StatusOK, wantBody: "<h1>Calculator</h1>"},
		{name: "calculate", method: http.MethodPost, path: "/api/v1/calculate", body: `{"operation":"add","a":2,"b":3}`, wantStatus: http.StatusOK, wantBody: `{"result":5}`},
		{name: "health", method: http.MethodGet, path: "/health", wantStatus: http.StatusOK, wantBody: `{"status":"ok"}`},
		{name: "API method rules still apply", method: http.MethodGet, path: "/api/v1/calculate", wantStatus: http.StatusMethodNotAllowed},
		{name: "unknown API path is not a static file", method: http.MethodGet, path: "/api/unknown", wantStatus: http.StatusNotFound},
		{name: "missing static file", method: http.MethodGet, path: "/missing.js", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body)))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := strings.TrimSpace(rec.Body.String()); tt.wantBody != "" && got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
		})
	}
}

func TestNewHandlerWithoutStaticDirServesOnlyAPI(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler("").ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
