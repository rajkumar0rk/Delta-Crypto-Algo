package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFoundationEndpointsFailClosed(t *testing.T) {
	handler := NewHandler()
	cases := []struct {
		name, method, path string
		status int
		state string
	}{
		{"liveness", http.MethodGet, "/health/live", http.StatusOK, "alive"},
		{"readiness", http.MethodGet, "/health/ready", http.StatusServiceUnavailable, "foundation_only"},
		{"post readiness", http.MethodPost, "/health/ready", http.StatusMethodNotAllowed, ""},
		{"post liveness", http.MethodPost, "/health/live", http.StatusMethodNotAllowed, ""},
		{"orders get", http.MethodGet, "/orders", http.StatusNotFound, ""},
		{"orders post", http.MethodPost, "/orders", http.StatusNotFound, ""},
		{"api orders post", http.MethodPost, "/api/v1/orders", http.StatusNotFound, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d", rec.Code, tc.status)
			}
			for key, expected := range map[string]string{
				"Cache-Control": "no-store",
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options": "DENY",
				"Referrer-Policy": "no-referrer",
			} {
				if got := rec.Header().Get(key); got != expected {
					t.Errorf("%s = %q, want %q", key, got, expected)
				}
			}
			if tc.state == "" {
				return
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("content type = %q", got)
			}
			var body struct {
				State string `json:"state"`
				TradeEnabled bool `json:"trade_enabled"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode health: %v", err)
			}
			if body.State != tc.state || body.TradeEnabled {
				t.Fatalf("unsafe health response: %+v", body)
			}
		})
	}
}
