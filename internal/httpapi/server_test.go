package httpapi

import (
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
