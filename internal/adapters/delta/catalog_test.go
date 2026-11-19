package delta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPublicCatalogPagination(t *testing.T) {
	requests := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet || r.URL.Path != "/v2/products" || r.Header.Get("Accept") != "application/json" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
		}
		if requests == 1 {
			if r.URL.Query().Has("after") {
				t.Error("first page has cursor")
			}
			fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":"next cursor"}}`)
		} else {
			if r.URL.Query().Get("after") != "next cursor" {
				t.Errorf("cursor not encoded: %s", r.URL)
			}
			fmt.Fprint(w, `{"success":true,"result":[{"id":2}],"meta":{"after":null}}`)
		}
	}))
	defer server.Close()
	client := PublicCatalog{BaseURL: server.URL, Client: server.Client()}
	client.Client.Timeout = time.Second
	result, err := client.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 || result.Pages != 2 || len(result.Products) != 2 || result.AsOf.IsZero() {
		t.Fatalf("incomplete result: %+v requests=%d", result, requests)
	}
	var p map[string]int
	if err := json.Unmarshal(result.Products[1], &p); err != nil || p["id"] != 2 {
		t.Fatalf("unexpected raw product %s: %v", result.Products[1], err)
	}
}

func TestPublicCatalogFailClosed(t *testing.T) {
	cases := []struct {
		name        string
		response    string
		status      int
		maxPages    int
		maxProducts int
	}{
		{"server error", "{}", 500, 0, 0},
		{"API failed", `{"success":false,"result":[{"id":1}]}`, 200, 0, 0},
		{"missing success", `{"result":[{"id":1}]}`, 200, 0, 0},
		{"empty", `{"success":true,"result":[]}`, 200, 0, 0},
		{"invalid json", "{", 200, 0, 0},
		{"extra json", `{"success":true,"result":[{"id":1}]} {}`, 200, 0, 0},
		{"bad product", `{"success":true,"result":[42]}`, 200, 0, 0},
		{"bad cursor", `{"success":true,"result":[{"id":1}],"meta":{"after":{}}}`, 200, 0, 0},
		{"repeated cursor", `{"success":true,"result":[{"id":1}],"meta":{"after":"same"}}`, 200, 0, 0},
		{"page limit", `{"success":true,"result":[{"id":1}],"meta":{"after":"next"}}`, 200, 1, 0},
		{"product limit", `{"success":true,"result":[{"id":1},{"id":2}]}`, 200, 0, 1},
		{"oversized body", strings.Repeat("x", maxPageBytes+1), 200, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.response)
			}))
			defer srv.Close()
			httpClient := srv.Client()
			httpClient.Timeout = time.Second
			c := PublicCatalog{BaseURL: srv.URL, Client: httpClient, MaxPages: tc.maxPages, MaxProducts: tc.maxProducts}
			got, err := c.Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 {
				t.Fatalf("must discard partial results: %+v, %v", got, err)
			}
		})
	}
}

func TestPublicCatalogRejectsUnsafeConfiguration(t *testing.T) {
	for _, tc := range []PublicCatalog{
		{BaseURL: "http://example.com", Client: &http.Client{Timeout: time.Second}},
		{BaseURL: "https://example.com/v2/products", Client: &http.Client{Timeout: time.Second}},
		{BaseURL: "https://example.com", Client: &http.Client{}},
		{BaseURL: "https://user:secret@example.com", Client: &http.Client{Timeout: time.Second}},
		{BaseURL: "https://example.com", Client: &http.Client{Timeout: time.Second}, MaxPages: 101},
	} {
		_, err := tc.Fetch(context.Background())
		if !errors.Is(err, ErrCatalogIncomplete) {
			t.Fatalf("expected invalid config: %+v, %v", tc, err)
		}
	}
	_, err := (PublicCatalog{}).Fetch(nil)
	if !errors.Is(err, ErrCatalogIncomplete) {
		t.Fatal(err)
	}
}
