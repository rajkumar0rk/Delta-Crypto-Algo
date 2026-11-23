package httpapi

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEndpointsAreReadOnly(t *testing.T) {
	for _, tc := range []struct {
		path  string
		want  int
		state string
	}{
		{"/health/live", http.StatusOK, "alive"},
		{"/health/ready", http.StatusServiceUnavailable, "foundation_only"},
		{"/orders", http.StatusNotFound, ""},
		{"/api/v1/orders", http.StatusNotFound, ""},
	} {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != tc.want {
				t.Fatalf("status %d != %d", rec.Code, tc.want)
			}
			if tc.state != "" && !strings.Contains(rec.Body.String(), tc.state) {
				t.Fatalf("body %s", rec.Body.String())
			}
			if rec.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("cache header missing")
			}
		})
	}
}

func TestRequestIDsAreServerGeneratedAndUnique(t *testing.T) {
	handler := NewHandler()
	seen := make(map[string]bool)
	for i := 0; i < 128; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
		req.Header.Set("X-Request-ID", "attacker-controlled")
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: unexpected status %d", i, rec.Code)
		}
		id := rec.Header().Get("X-Request-ID")
		if len(id) != 32 {
			t.Fatalf("request %d: expected 128-bit hex ID, got %q", i, id)
		}
		if _, err := hex.DecodeString(id); err != nil {
			t.Fatalf("request %d: invalid hex ID: %v", i, err)
		}
		if seen[id] || id == "attacker-controlled" {
			t.Fatalf("request %d: repeated or untrusted ID %q", i, id)
		}
		seen[id] = true
	}
}

func TestRequestIDOnDeniedRoutes(t *testing.T) {
	for _, tc := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/health/ready", http.StatusServiceUnavailable},
		{http.MethodPost, "/health/live", http.StatusMethodNotAllowed},
		{http.MethodPost, "/orders", http.StatusNotFound},
	} {
		rec := httptest.NewRecorder()
		NewHandler().ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.status {
			t.Fatalf("%s %s: status %d, want %d", tc.method, tc.path, rec.Code, tc.status)
		}
		if len(rec.Header().Get("X-Request-ID")) != 32 {
			t.Fatalf("%s %s: missing request ID", tc.method, tc.path)
		}
	}
}

func TestSecurityHeadersOnEveryResponse(t *testing.T) {
	handler := NewHandler()
	for _, tc := range []struct {
		method string
		path   string
		want   int
	}{
		{http.MethodGet, "/health/live", http.StatusOK},
		{http.MethodGet, "/health/ready", http.StatusServiceUnavailable},
		{http.MethodPost, "/health/live", http.StatusMethodNotAllowed},
		{http.MethodPost, "/orders", http.StatusNotFound},
		{http.MethodOptions, "/api/v1/orders", http.StatusNotFound},
	} {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Origin", "https://untrusted.example")
			handler.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status %d, want %d", rec.Code, tc.want)
			}
			for name, want := range map[string]string{
				"Cache-Control":                "no-store",
				"Content-Security-Policy":      "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'",
				"Cross-Origin-Resource-Policy": "same-origin",
				"X-Content-Type-Options":       "nosniff",
				"X-Frame-Options":              "DENY",
				"Referrer-Policy":              "no-referrer",
			} {
				if got := rec.Header().Get(name); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
				t.Errorf("unexpected CORS grant %q", got)
			}
			if len(rec.Header().Get("X-Request-ID")) != 32 {
				t.Error("missing generated request ID")
			}
		})
	}
}

func TestHealthContractAndUnsafeMethodsFailClosed(t *testing.T) {
	handler := NewHandler()
	for _, tc := range []struct {
		method string
		path   string
		want   int
		state  string
	}{
		{http.MethodGet, "/health/live", http.StatusOK, "alive"},
		{http.MethodGet, "/health/ready", http.StatusServiceUnavailable, "foundation_only"},
		{http.MethodHead, "/health/live", http.StatusOK, ""},
		{http.MethodHead, "/health/ready", http.StatusServiceUnavailable, ""},
		{http.MethodPost, "/health/ready", http.StatusMethodNotAllowed, ""},
		{http.MethodPut, "/health/live", http.StatusMethodNotAllowed, ""},
		{http.MethodDelete, "/health/live", http.StatusMethodNotAllowed, ""},
		{http.MethodPatch, "/health/ready", http.StatusMethodNotAllowed, ""},
		{http.MethodPost, "/api/v1/orders", http.StatusNotFound, ""},
		{http.MethodDelete, "/api/v1/orders/123", http.StatusNotFound, ""},
		{http.MethodGet, "/health/live/extra", http.StatusNotFound, ""},
	} {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Origin", "https://untrusted.example")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status %d, want %d", rec.Code, tc.want)
			}
			if rec.Header().Get("Cache-Control") != "no-store" {
				t.Error("response may be cached")
			}
			if rec.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Error("unexpected cross-origin grant")
			}
			if tc.state != "" {
				var got health
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("invalid health JSON: %v", err)
				}
				if got.State != tc.state || got.TradeEnabled {
					t.Fatalf("unsafe health response: %+v", got)
				}
				if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
					t.Fatalf("health content type %q", ct)
				}
			}
			if tc.method == http.MethodHead && rec.Body.Len() != 0 {
				t.Fatalf("HEAD response unexpectedly contains body: %q", rec.Body.String())
			}
		})
	}
}
