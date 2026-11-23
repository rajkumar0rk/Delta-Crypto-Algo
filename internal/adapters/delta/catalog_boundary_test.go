package delta_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/adapters/delta"
	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/domain/catalog"
)

// A successful public REST fetch is not a verified venue contract mapping.
// Synthetic records deliberately omit critical contract fields; they must not
// become tradeable by merely copying transport JSON into the domain model.
func TestPublicCatalogRequiresVerifiedMappingBeforeTradeability(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v2/products" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"result":[{"id":1001,"symbol":"SYNTH-CRYPTO-PERP"},{"id":1002,"symbol":"SYNTH-RWA-OPTION"}],"meta":{"after":null}}`))
	}))
	defer server.Close()

	raw, err := (delta.PublicCatalog{BaseURL: server.URL, Client: server.Client()}).Fetch(context.Background())
	if err != nil {
		t.Fatalf("public fetch: %v", err)
	}
	if raw.Pages != 1 || len(raw.Products) != 2 {
		t.Fatalf("unexpected raw snapshot: pages=%d products=%d", raw.Pages, len(raw.Products))
	}
	products := make([]catalog.Product, 0, len(raw.Products))
	for _, record := range raw.Products {
		var p catalog.Product
		if err := json.Unmarshal(record, &p); err != nil {
			t.Fatalf("decode synthetic record: %v", err)
		}
		if err := p.Validate(); err == nil {
			t.Fatalf("unverified product %d accepted without contract metadata", p.ID)
		}
		products = append(products, p)
	}
	snapshot := catalog.Snapshot{Source: "synthetic-unverified", AsOf: raw.AsOf, Complete: true, Products: products}
	if err := snapshot.Tradeable(1001, raw.AsOf.Add(time.Second), time.Minute); !errors.Is(err, catalog.ErrNotTradeable) {
		t.Fatalf("unverified product should fail closed: %v", err)
	}
	if err := snapshot.Tradeable(1002, raw.AsOf.Add(time.Second), time.Minute); !errors.Is(err, catalog.ErrNotTradeable) {
		t.Fatalf("unverified option should fail closed: %v", err)
	}
}

// A transport failure on a later page must discard the first page completely;
// no partial product set may be accidentally published.
func TestPublicCatalogSecondPageFailureIsAtomic(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("after") == "" {
			_, _ = w.Write([]byte(`{"success":true,"result":[{"id":1}],"meta":{"after":"next"}}`))
			return
		}
		http.Error(w, "upstream unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	raw, err := (delta.PublicCatalog{BaseURL: server.URL, Client: server.Client()}).Fetch(context.Background())
	if !errors.Is(err, delta.ErrCatalogIncomplete) {
		t.Fatalf("expected incomplete error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected two requests, got %d", calls)
	}
	if len(raw.Products) != 0 || raw.Pages != 0 || !raw.AsOf.IsZero() {
		t.Fatalf("partial catalog escaped failure: %+v", raw)
	}
}
