package httpapi

import (
	"net/http"
	"encoding/hex"
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
