package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestFoundationDeniesTradingRoutes guards against accidentally exposing
// write-capable endpoints before independent risk, auth and OMS gates exist.
func TestFoundationDeniesTradingRoutes(t *testing.T) {
	handler := NewHandler()
	for _, method := range []string{
		http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete,
	} {
		for _, path := range []string{
			"/orders", "/orders/123", "/api/v1/orders",
			"/v1/orders", "/v1/strategies/123/live",
			"/v1/strategies/123/activate", "/v1/positions/close",
		} {
			t.Run(method+" "+path, func(t *testing.T) {
				req := httptest.NewRequest(method, path, strings.NewReader("{}"))
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				if rec.Code != http.StatusNotFound {
					t.Fatalf("unexpected trading route status: %d", rec.Code)
				}
				if rec.Header().Get("Cache-Control") != "no-store" ||
					rec.Header().Get("X-Content-Type-Options") != "nosniff" ||
					rec.Header().Get("X-Frame-Options") != "DENY" ||
					rec.Header().Get("Referrer-Policy") != "no-referrer" {
					t.Fatalf("missing security headers: %v", rec.Header())
				}
			})
		}
	}
}

func TestFoundationReadinessCannotSignalTradePermission(t *testing.T) {
	handler := NewHandler()
	for _, path := range []string{"/health/live", "/health/ready"} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
			want := http.StatusOK
			if path == "/health/ready" {
				want = http.StatusServiceUnavailable
			}
			if rec.Code != want {
				t.Fatalf("status %d, want %d", rec.Code, want)
			}
			var body struct {
				TradeEnabled *bool `json:"trade_enabled"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.TradeEnabled == nil || *body.TradeEnabled {
				t.Fatalf("trade permission must be explicitly false: %s", rec.Body.String())
			}
		})
	}
}
