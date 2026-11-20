package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// API responses, including 404s and fail-closed readiness, must not become
// executable or embeddable content when consumed by an untrusted browser.
func TestAPIResponseSecurityHeaders(t *testing.T) {
	for _, tc := range []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"liveness", http.MethodGet, "/health/live", http.StatusOK},
		{"unready", http.MethodGet, "/health/ready", http.StatusServiceUnavailable},
		{"missing", http.MethodGet, "/missing", http.StatusNotFound},
		{"unsafe-method", http.MethodPost, "/health/live", http.StatusMethodNotAllowed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			NewHandler().ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.wantStatus {
				t.Fatalf("status=%d want=%d", rec.Code, tc.wantStatus)
			}
			for header, want := range map[string]string{
				"Cache-Control":                "no-store"
				"X-Content-Type-Options":       "nosniff"
				"X-Frame-Options":              "DENY"
				"Referrer-Policy":              "no-referrer"
				"Content-Security-Policy":      "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'"
				"Cross-Origin-Resource-Policy": "same-origin"
			} {
				if got := rec.Header().Get(header); got != want {
					t.Errorf("%s=%q want=%q", header, got, want)
				}
			}
		})
	}
}
