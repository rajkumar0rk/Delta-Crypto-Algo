package delta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestPublicCatalogPagination(t *testing.T) {
	requests := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
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
		{"valid JSON plus oversized whitespace", `{"success":true,"result":[{"id":1}]}` + strings.Repeat(" ", maxPageBytes), 200, 0, 0},
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

func TestPublicCatalogRejectsRedirectsWithoutFollowing(t *testing.T) {
	redirected := 0
	target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirected++
		fmt.Fprint(w, `{"success":true,"result":[{"id":1}]}`)
	}))
	defer target.Close()
	origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/v2/products", http.StatusTemporaryRedirect)
	}))
	defer origin.Close()
	client := origin.Client()
	client.Timeout = time.Second
	// The caller's client must remain unchanged.
	got, err := (PublicCatalog{BaseURL: origin.URL, Client: client}).Fetch(context.Background())
	if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 {
		t.Fatalf("redirect must fail closed: %+v %v", got, err)
	}
	if redirected != 0 {
		t.Fatalf("redirect followed to unexpected origin: %d", redirected)
	}
	if client.CheckRedirect != nil {
		t.Fatal("caller HTTP client mutated")
	}
}

func TestPublicCatalogRejectsUnexpectedHTTPRepresentation(t *testing.T) {
	cases := []struct {
		name        string
		contentType string
		encoding    string
	}{
		{name: "html content type", contentType: "text/html"},
		{name: "json suffix is not JSON", contentType: "application/problem+json"},
		{name: "gzip encoding", contentType: "application/json", encoding: "gzip"},
		{name: "br encoding", contentType: "application/json", encoding: "br"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Accept-Encoding") != "identity" {
					t.Errorf("unexpected encoding negotiation: %q", r.Header.Get("Accept-Encoding"))
				}
				w.Header().Set("Content-Type", tc.contentType)
				if tc.encoding != "" {
					w.Header().Set("Content-Encoding", tc.encoding)
				}
				fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":null}}`)
			}))
			defer srv.Close()
			client := srv.Client()
			client.Timeout = time.Second
			got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 {
				t.Fatalf("must reject representation without publishing products: %+v, %v", got, err)
			}
		})
	}
}

func TestPublicCatalogAcceptsDeclaredJSONMediaType(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Encoding", "identity")
		fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":null}}`)
	}))
	defer srv.Close()
	client := srv.Client()
	client.Timeout = time.Second
	got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
	if err != nil || len(got.Products) != 1 {
		t.Fatalf("valid JSON response rejected: %+v, %v", got, err)
	}
}

func TestPublicCatalogRejectsAmbiguousRepresentationHeaders(t *testing.T) {
	cases := []struct {
		name         string
		contentTypes []string
		encodings    []string
	}{
		{name: "duplicate content type", contentTypes: []string{"application/json", "text/html"}},
		{name: "duplicate content encoding", contentTypes: []string{"application/json"}, encodings: []string{"identity", "gzip"}},
		{name: "comma joined content type", contentTypes: []string{"application/json, text/html"}},
		{name: "non utf8 charset", contentTypes: []string{"application/json; charset=iso-8859-1"}},
		{name: "invalid media parameter", contentTypes: []string{"application/json; charset="}},
		{name: "unknown media parameter", contentTypes: []string{"application/json; version=1"}},
		{name: "empty declared content type", contentTypes: []string{""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, value := range tc.contentTypes {
					w.Header().Add("Content-Type", value)
				}
				for _, value := range tc.encodings {
					w.Header().Add("Content-Encoding", value)
				}
				fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":null}}`)
			}))
			defer srv.Close()
			client := srv.Client()
			client.Timeout = time.Second
			got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 {
				t.Fatalf("ambiguous representation accepted: %+v %v", got, err)
			}
		})
	}
}

func TestPublicCatalogRejectsSecondPageRepresentationWithoutPartialPublication(t *testing.T) {
	requests := 0
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":"next"}}`)
			return
		}
		if r.URL.Query().Get("after") != "next" {
			t.Errorf("unexpected second-page cursor: %s", r.URL)
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `{"success":true,"result":[{"id":2}],"meta":{"after":null}}`)
	}))
	defer srv.Close()
	client := srv.Client()
	client.Timeout = time.Second
	got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
	if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 || !got.AsOf.IsZero() {
		t.Fatalf("second-page representation must discard complete snapshot: %+v, %v", got, err)
	}
	if requests != 2 {
		t.Fatalf("expected two pages, got %d", requests)
	}
}

func TestPublicCatalogSnapshotTimestampAfterFinalPage(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		close(started)
		<-release
		fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":null}}`)
	}))
	defer srv.Close()
	httpClient := srv.Client()
	httpClient.Timeout = 5 * time.Second
	type outcome struct {
		snapshot RawSnapshot
		err      error
	}
	done := make(chan outcome, 1)
	go func() {
		snapshot, err := (PublicCatalog{BaseURL: srv.URL, Client: httpClient}).Fetch(context.Background())
		done <- outcome{snapshot, err}
	}()
	<-started
	completedAfter := time.Now().UTC()
	close(release)
	got := <-done
	if got.err != nil || len(got.snapshot.Products) != 1 || got.snapshot.AsOf.Before(completedAfter) {
		t.Fatalf("snapshot timestamp predates completed fetch: %+v %v", got.snapshot, got.err)
	}
}

func TestPublicCatalogCanceledAtFinalPageFailsClosed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cancel()
		fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":null}}`)
	}))
	defer srv.Close()
	httpClient := srv.Client()
	httpClient.Timeout = time.Second
	got, err := (PublicCatalog{BaseURL: srv.URL, Client: httpClient}).Fetch(ctx)
	if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 {
		t.Fatalf("canceled fetch must discard snapshot: %+v %v", got, err)
	}
}

// countingTransport detects forbidden network I/O even if the HTTP stack's
// configured transport does not honor a canceled context.
type countingTransport struct{ calls atomic.Int32 }

func (c *countingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	c.calls.Add(1)
	return nil, errors.New("network must not be reached")
}

func TestPublicCatalogCanceledContextNeverReachesTransport(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	transport := &countingTransport{}
	client := &http.Client{Timeout: time.Second, Transport: transport}
	got, err := (PublicCatalog{BaseURL: "https://example.invalid", Client: client}).Fetch(ctx)
	if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 || !got.AsOf.IsZero() {
		t.Fatalf("canceled request did not fail closed: %+v %v", got, err)
	}
	if n := transport.calls.Load(); n != 0 {
		t.Fatalf("canceled request reached transport %d times", n)
	}
}

func TestPublicCatalogRejectsOpaqueHTTPSOriginBeforeNetwork(t *testing.T) {
	transport := &countingTransport{}
	client := &http.Client{Timeout: time.Second, Transport: transport}
	got, err := (PublicCatalog{BaseURL: "https:opaque-origin", Client: client}).Fetch(context.Background())
	if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 {
		t.Fatalf("opaque origin accepted: %+v %v", got, err)
	}
	if n := transport.calls.Load(); n != 0 {
		t.Fatalf("invalid origin reached transport %d times", n)
	}
}
// A missing or empty cursor cannot prove that every venue page was fetched.
// Only an explicit JSON null terminates a complete catalog.
func TestPublicCatalogRequiresExplicitTerminalCursor(t *testing.T) {
	cases := []struct {
		name string
		meta string
	}{
		{name: "empty string", meta: `"after":""`},
		{name: "whitespace string", meta: `"after":"   "`},
		{name: "missing after", meta: `"other":null`},
		{name: "empty metadata", meta: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"success":true,"result":[{"id":1}],"meta":{%s}}`, tc.meta)
			}))
			defer srv.Close()
			client := srv.Client()
			client.Timeout = time.Second
			got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 || !got.AsOf.IsZero() {
				t.Fatalf("unproven terminal cursor published data: %+v, %v", got, err)
			}
		})
	}
}
