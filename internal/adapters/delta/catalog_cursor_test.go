package delta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Cursor values originate from an external service and must never turn into
// ambiguous HTTP request targets. Reject unsafe cursors before a second GET.
func TestPublicCatalogRejectsUnsafeCursorsBeforeNextRequest(t *testing.T) {
	cases := []struct {
		name, cursor string
	}{
		{"newline", `"next\npage"`},
		{"carriage return", `"next\rpage"`},
		{"tab", `"next\tpage"`},
		{"nul", `"next\u0000page"`},
		{"whitespace only", `"   "`},
		{"oversized", `"` + strings.Repeat("a", 257) + `"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				fmt.Fprintf(w, `{"success":true,"result":[{"id":1}],"meta":{"after":%s}}`, tc.cursor)
			}))
			defer srv.Close()
			client := srv.Client()
			client.Timeout = time.Second
			got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 {
				t.Fatalf("unsafe cursor must discard snapshot: %+v, %v", got, err)
			}
			if requests != 1 {
				t.Fatalf("unsafe cursor caused %d requests, want one", requests)
			}
		})
	}
}

func TestPublicCatalogSecondPageFailureDiscardsFirstPage(t *testing.T) {
	for _, status := range []int{http.StatusBadGateway, http.StatusTooManyRequests} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			requests := 0
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if requests == 1 {
					fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":"page2"}}`)
					return
				}
				if r.URL.Query().Get("after") != "page2" {
					t.Errorf("missing second-page cursor: %s", r.URL)
				}
				w.WriteHeader(status)
			}))
			defer srv.Close()
			client := srv.Client()
			client.Timeout = time.Second
			got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 || requests != 2 {
				t.Fatalf("partial page leaked: %+v, err=%v, requests=%d", got, err, requests)
			}
		})
	}
}

func TestPublicCatalogCanceledContextNeverPublishes(t *testing.T) {
	requests := 0
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		fmt.Fprint(w, `{"success":true,"result":[{"id":1}],"meta":{"after":null}}`)
	}))
	defer srv.Close()
	client := srv.Client()
	client.Timeout = time.Second
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(ctx)
	if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || requests != 0 {
		t.Fatalf("canceled context must not fetch or publish: %+v, err=%v, requests=%d", got, err, requests)
	}
}
