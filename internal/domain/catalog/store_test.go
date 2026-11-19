package catalog

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func testSnapshot(t time.Time) Snapshot {
	p := Product{}
	p.ID = 1
	p.Symbol = "TEST-PERP"
	p.Underlying = "TEST"
	p.Kind = Perpetual
	p.Status = Active
	p.TickSize = "0.01"
	p.LotSize = "1"
	p.ContractMultiplier = "1"
	p.QuoteCurrency = "USD"
	p.SettlementCurrency = "USD"
	return Snapshot{Source: "fixture", AsOf: t, Complete: true, Products: []Product{p}}
}

func TestStoreFailClosedAndAtomic(t *testing.T) {
	var store Store
	now := time.Now().UTC()
	if _, err := store.Snapshot(); !errors.Is(err, ErrCatalogUnavailable) {
		t.Fatalf("uninitialized store: %v", err)
	}
	if err := store.Tradeable(1, now, time.Minute); !errors.Is(err, ErrNotTradeable) {
		t.Fatalf("uninitialized tradeability: %v", err)
	}
	good := testSnapshot(now)
	if err := store.Publish(good); err != nil {
		t.Fatal(err)
	}
	good.Products[0].Symbol = "CORRUPTED"
	got, err := store.Snapshot()
	if err != nil || got.Products[0].Symbol != "TEST-PERP" {
		t.Fatalf("input alias: %v %+v", err, got)
	}
	got.Products[0].Status = Delisted
	got.Products = nil
	got2, _ := store.Snapshot()
	if len(got2.Products) != 1 || got2.Products[0].Status != Active {
		t.Fatal("output alias")
	}
	if err := store.Tradeable(1, now.Add(time.Second), time.Minute); err != nil {
		t.Fatal(err)
	}
	bad := testSnapshot(now.Add(time.Second))
	bad.Complete = false
	if err := store.Publish(bad); err == nil {
		t.Fatal("accepted incomplete refresh")
	}
	older := testSnapshot(now)
	if err := store.Publish(older); err == nil {
		t.Fatal("accepted rollback")
	}
	if err := store.Tradeable(1, now.Add(2*time.Minute), time.Minute); !errors.Is(err, ErrNotTradeable) {
		t.Fatalf("stale data must fail closed: %v", err)
	}
	if err := store.Tradeable(1, now.Add(time.Second), time.Minute); err != nil {
		t.Fatalf("failed refresh changed catalog: %v", err)
	}
}

func TestStoreExpiryPointerIsolation(t *testing.T) {
	var store Store
	now := time.Now().UTC()
	exp := now.Add(24*time.Hour)
	expectedExpiry := exp
	snap := testSnapshot(now)
	snap.Products[0].Kind = CallOption
	snap.Products[0].Strike = "100"
	snap.Products[0].Expiry = &exp
	if err := store.Publish(snap); err != nil {
		t.Fatal(err)
	}
	*snap.Products[0].Expiry = now.Add(-time.Hour)
	out, _ := store.Snapshot()
	if !out.Products[0].Expiry.Equal(expectedExpiry) {
		t.Fatal("input expiry pointer aliased")
	}
	*out.Products[0].Expiry = now.Add(-time.Hour)
	if err := store.Tradeable(1, now.Add(time.Second), time.Minute); err != nil {
		t.Fatalf("output expiry pointer aliased: %v", err)
	}
}

func TestStoreConcurrentReadersAndPublishers(t *testing.T) {
	var store Store
	start := time.Now().UTC()
	if err := store.Publish(testSnapshot(start)); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				got, err := store.Snapshot()
				if err != nil || len(got.Products) != 1 || got.Products[0].ID != 1 {
					t.Errorf("torn snapshot: %v %+v", err, got)
					return
				}
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			if err := store.Publish(testSnapshot(start.Add(time.Duration(i)*time.Second))); err != nil {
				t.Error(err)
				return
			}
		}
	}()
	wg.Wait()
}
