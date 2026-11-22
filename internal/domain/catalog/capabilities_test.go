package catalog

import (
	"errors"
	"testing"
	"time"
)

func TestCapabilityAdmissionFailClosed(t *testing.T) {
	now := time.Date(2026, 11, 22, 7, 0, 0, 0, time.UTC)
	product := Product{ID: 987654, Symbol: "EXAMPLE-RWA-2027", Underlying: "RWA", Kind: DatedFuture, Status: Active, TickSize: "0.01", LotSize: "1", ContractMultiplier: "0.1", QuoteCurrency: "USD", SettlementCurrency: "USD"}
	exp := now.Add(24 * time.Hour)
	product.Expiry = &exp
	good := Capability{ProductID: product.ID, VerifiedAt: now.Add(-time.Minute), Source: "verified-fixture", AllowLimit: true, AllowPostOnly: true, MakerFeeRate: "0", TakerFeeRate: "0.001", InitialMarginRate: "0.1", MaintenanceMarginRate: "0.05"}
	base := CapabilitySet{AsOf: now, Source: "verified-fixture", Products: map[int64]Capability{product.ID: good}}
	if _, err := base.VerifiedFor(product, now, time.Hour); err != nil {
		t.Fatalf("expected complete fixture: %v", err)
	}
	cases := []struct {
		name   string
		change func(*CapabilitySet)
	}{
		{"missing product", func(s *CapabilitySet) { s.Products = nil }},
		{"wrong product id", func(s *CapabilitySet) {
			c := good
			c.ProductID++
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"missing provenance", func(s *CapabilitySet) { s.Source = "" }},
		{"source mismatch", func(s *CapabilitySet) {
			c := good
			c.Source = "other"
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"future verification", func(s *CapabilitySet) {
			c := good
			c.VerifiedAt = now.Add(time.Second)
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"stale evidence", func(s *CapabilitySet) {
			c := good
			c.VerifiedAt = now.Add(-2 * time.Hour)
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"unknown fees", func(s *CapabilitySet) {
			c := good
			c.TakerFeeRate = ""
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"negative margin", func(s *CapabilitySet) {
			c := good
			c.InitialMarginRate = "-0.1"
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"exponent", func(s *CapabilitySet) {
			c := good
			c.MakerFeeRate = "1e-3"
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"over one", func(s *CapabilitySet) {
			c := good
			c.MaintenanceMarginRate = "1.01"
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"post only without limit", func(s *CapabilitySet) {
			c := good
			c.AllowLimit = false
			s.Products = map[int64]Capability{product.ID: c}
		}},
		{"no order type", func(s *CapabilitySet) {
			c := good
			c.AllowLimit = false
			c.AllowPostOnly = false
			s.Products = map[int64]Capability{product.ID: c}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := base
			tc.change(&s)
			if _, err := s.VerifiedFor(product, now, time.Hour); !errors.Is(err, ErrCapabilityUnverified) {
				t.Fatalf("expected denied, got %v", err)
			}
		})
	}
	if _, err := base.VerifiedFor(product, now.Add(2*time.Hour), time.Hour); !errors.Is(err, ErrCapabilityUnverified) {
		t.Fatalf("stale set: %v", err)
	}
	if _, err := base.VerifiedFor(product, now, -time.Second); !errors.Is(err, ErrCapabilityUnverified) {
		t.Fatalf("invalid max age: %v", err)
	}
}

func TestNonnegativeRatio(t *testing.T) {
	for _, v := range []string{"0", "0.0", "0.001", "1", "1.000"} {
		if !nonnegativeRatio(v) {
			t.Errorf("valid %q rejected", v)
		}
	}
	for _, v := range []string{"", "-0.1", "+0.1", "2", "1.001", "0.", "0..1", ".5", "00.1", "0x1", "1e-3", " 0", "0/1", "0_1"} {
		if nonnegativeRatio(v) {
			t.Errorf("invalid %q accepted", v)
		}
	}
}
