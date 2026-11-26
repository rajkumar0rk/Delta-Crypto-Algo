package catalog

import (
	"errors"
	"testing"
	"time"
)

func TestMetadataAdmissionRequiresIndependentEvidence(t *testing.T) {
	now := time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)
	s := testSnapshot(now)
	scope := AccountScope{AccountID: "synthetic-account", Environment: "demo"}
	e := EligibilityEvidence{Scope: scope, ProductID: 1, CatalogSource: s.Source, CatalogAsOf: now, VerifiedAt: now, Source: "synthetic-account-fixture", Allowed: true}
	c := Capability{ProductID: 1, VerifiedAt: now, Source: "synthetic-capability-fixture", AllowLimit: true, MakerFeeRate: "0.001", TakerFeeRate: "0.002", InitialMarginRate: "0.1", MaintenanceMarginRate: "0.05"}
	caps := CapabilitySet{AsOf: now, Source: c.Source, Products: map[int64]Capability{1: c}}
	if err := CheckMetadataAdmission(s, caps, e, scope, 1, now, time.Minute); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name   string
		mutate func(*EligibilityEvidence)
	}{
		{"denied", func(x *EligibilityEvidence) { x.Allowed = false }},
		{"wrong account", func(x *EligibilityEvidence) { x.Scope.AccountID = "other" }},
		{"wrong environment", func(x *EligibilityEvidence) { x.Scope.Environment = "live" }},
		{"wrong product", func(x *EligibilityEvidence) { x.ProductID = 2 }},
		{"wrong source", func(x *EligibilityEvidence) { x.CatalogSource = "other" }},
		{"previous observation", func(x *EligibilityEvidence) { x.CatalogAsOf = now.Add(-time.Second) }},
		{"missing source", func(x *EligibilityEvidence) { x.Source = "" }},
		{"old evidence", func(x *EligibilityEvidence) { x.VerifiedAt = now.Add(-time.Second) }},
		{"future evidence", func(x *EligibilityEvidence) { x.VerifiedAt = now.Add(time.Second) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bad := e
			tc.mutate(&bad)
			if err := CheckMetadataAdmission(s, caps, bad, scope, 1, now, time.Minute); !errors.Is(err, ErrEligibilityUnverified) {
				t.Fatalf("expected denial: %v", err)
			}
		})
	}
	if err := CheckMetadataAdmission(s, CapabilitySet{}, e, scope, 1, now, time.Minute); !errors.Is(err, ErrCapabilityUnverified) {
		t.Fatalf("missing capability not rejected: %v", err)
	}
	c.VerifiedAt = now.Add(-time.Second)
	caps.Products[1] = c
	if err := CheckMetadataAdmission(s, caps, e, scope, 1, now, time.Minute); !errors.Is(err, ErrCapabilityUnverified) {
		t.Fatalf("prior contract capability not rejected: %v", err)
	}
	s.Products[0].Status = Delisted
	if err := CheckMetadataAdmission(s, caps, e, scope, 1, now, time.Minute); !errors.Is(err, ErrNotTradeable) {
		t.Fatalf("delisted product not rejected: %v", err)
	}
}

func TestCapabilityAgeMeasuredAtUseNotSetPublication(t *testing.T) {
	now := time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)
	p := testSnapshot(now).Products[0]
	c := Capability{ProductID: p.ID, Source: "fixture", VerifiedAt: now.Add(-90 * time.Second), AllowMarket: true, MakerFeeRate: "0", TakerFeeRate: "0", InitialMarginRate: "0.1", MaintenanceMarginRate: "0.05"}
	set := CapabilitySet{Source: c.Source, AsOf: now.Add(-45 * time.Second), Products: map[int64]Capability{p.ID: c}}
	if _, err := set.VerifiedFor(p, now, time.Minute); !errors.Is(err, ErrCapabilityUnverified) {
		t.Fatalf("stale individual evidence accepted: %v", err)
	}
}
