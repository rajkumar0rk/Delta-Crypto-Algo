package delta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/domain/catalog"
)

func TestCatalogCachePublishesTermsAndProductsTogether(t *testing.T) {
	var cache CatalogCache
	if _, _, err := cache.Read(); err == nil {
		t.Fatal("empty cache available")
	}
	var response atomic.Value
	response.Store(productFixture)
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"result":[` + response.Load().(string) + `],"meta":{"after":null}}`))
	}))
	defer srv.Close()
	client := srv.Client()
	client.Timeout = time.Second
	adapter := PublicCatalog{BaseURL: srv.URL, Client: client}
	lots := fixtureLots(time.Now().Add(-time.Second))
	if _, err := cache.Refresh(context.Background(), adapter, lots, time.Minute); err != nil {
		t.Fatal(err)
	}
	first, terms, err := cache.Read()
	if err != nil || first.Products[0].TickSize != "0.05" || terms[42].MakerCommission != "0.0002" {
		t.Fatalf("first publication: %+v %+v %v", first, terms, err)
	}
	first.Products[0].TickSize = "changed"
	delete(terms, 42)
	unchanged, unchangedTerms, err := cache.Read()
	if err != nil || unchanged.Products[0].TickSize != "0.05" || unchangedTerms[42].MakerCommission != "0.0002" {
		t.Fatal("caller mutated cache")
	}
	response.Store(strings.Replace(strings.Replace(productFixture, "0.05", "0.10", 1), "0.0002", "0.0003", 1))
	changes, err := cache.Refresh(context.Background(), adapter, lots, time.Minute)
	if err != nil || len(changes) != 1 || changes[0].Kind != catalog.ContractChanged {
		t.Fatalf("contract transition: %+v %v", changes, err)
	}
	second, secondTerms, err := cache.Read()
	if err != nil || second.Products[0].TickSize != "0.10" || secondTerms[42].MakerCommission != "0.0003" || !second.AsOf.After(unchanged.AsOf) {
		t.Fatalf("second publication: %+v %+v %v", second, secondTerms, err)
	}
	response.Store(strings.Replace(productFixture, "0.05", "bad", 1))
	if _, err := cache.Refresh(context.Background(), adapter, lots, time.Minute); err == nil {
		t.Fatal("invalid refresh accepted")
	}
	after, afterTerms, err := cache.Read()
	if err != nil || !after.AsOf.Equal(second.AsOf) || after.Products[0].TickSize != "0.10" || afterTerms[42].MakerCommission != "0.0003" {
		t.Fatal("failed refresh changed cached observation")
	}
}
