package catalog

import (
	"testing"
	"time"
)

func TestClassifyChangesDeterministicAndConservative(t *testing.T) {
	before := fixture()
	after := fixture()
	after.AsOf = after.AsOf.Add(time.Minute)
	after.Products[0].TickSize = "0.00000002"
	after.Products[0].Status = Delisted
	after.Products[1].Status = Inactive
	changes, err := ClassifyChanges(before, after)
	if err != nil || len(changes) != 2 {
		t.Fatalf("changes=%+v err=%v", changes, err)
	}
	if changes[0].ID != 10 || changes[0].Kind != ContractChanged {
		t.Fatalf("contract change must take precedence: %+v", changes[0])
	}
	if changes[1].ID != 11 || changes[1].Kind != ProductDisabled {
		t.Fatalf("disabled option: %+v", changes[1])
	}
	// Audit copies must be detached from both snapshots.
	changes[1].After.Symbol = "MUTATED"
	if after.Products[1].Symbol == "MUTATED" {
		t.Fatal("audit aliases original product")
	}
	reactivated := after
	reactivated.AsOf = after.AsOf.Add(time.Second)
	reactivated.Products = append([]Product(nil), after.Products...)
	reactivated.Products[1].Status = Active
	changes, err = ClassifyChanges(after, reactivated)
	if err != nil || len(changes) != 1 || changes[0].Kind != ProductEnabled {
		t.Fatalf("reactivation: %+v err=%v", changes, err)
	}
}

func TestClassifyChangesListingsRemovalsAndExpiry(t *testing.T) {
	before := fixture()
	after := fixture()
	after.AsOf = after.AsOf.Add(time.Minute)
	after.Products = after.Products[1:]
	added := before.Products[0]
	added.ID = 99
	added.Symbol = "NEW-PERP"
	after.Products = append(after.Products, added)
	changes, err := ClassifyChanges(before, after)
	if err != nil || len(changes) != 2 || changes[0].Kind != ProductRemoved || changes[1].Kind != ProductListed {
		t.Fatalf("listing/removal: %+v err=%v", changes, err)
	}
	later := fixture()
	later.AsOf = later.AsOf.Add(time.Minute)
	expiry := later.Products[1].Expiry.Add(time.Hour)
	later.Products[1].Expiry = &expiry
	changes, err = ClassifyChanges(before, later)
	if err != nil || len(changes) != 1 || changes[0].Kind != ContractChanged {
		t.Fatalf("expiry change: %+v err=%v", changes, err)
	}
	if _, err := ClassifyChanges(before, before); err == nil {
		t.Fatal("non-monotonic snapshots accepted")
	}
}

func TestClassifyExactDecimalScaleEquivalenceAndRealChange(t *testing.T) {
	before := fixture()
	after := fixture()
	after.AsOf = after.AsOf.Add(time.Minute)
	after.Products[0].TickSize = decimalWithExtraScale(after.Products[0].TickSize)
	after.Products[0].LotSize = decimalWithExtraScale(after.Products[0].LotSize)
	if after.Products[1].Strike != "" {
		after.Products[1].Strike = decimalWithExtraScale(after.Products[1].Strike)
	}
	changes, err := ClassifyChanges(before, after)
	if err != nil || len(changes) != 0 {
		t.Fatalf("scale-only change should not classify: %+v err=%v", changes, err)
	}
	after.Products[0].Status = Delisted
	changes, err = ClassifyChanges(before, after)
	if err != nil || len(changes) != 1 || changes[0].Kind != ProductDisabled {
		t.Fatalf("status change with decimal formatting must classify as disabled: %+v err=%v", changes, err)
	}
}

 
// Every field used by the shared contract comparison must cause both a Diff
// event and a ContractChanged classification. Status-only changes must not.
func TestContractFieldClassificationParity(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Product)
	}{
		{"symbol", func(p *Product) { p.Symbol = "CHANGED-PERP" }},
		{"underlying", func(p *Product) { p.Underlying = "CHANGED" }},
		{"kind", func(p *Product) {
			p.Kind = DatedFuture
			expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
			p.Expiry = &expiry
		}},
		{"tick", func(p *Product) { p.TickSize = "0.00000003" }},
		{"lot", func(p *Product) { p.LotSize = "3" }},
		{"multiplier", func(p *Product) { p.ContractMultiplier = "3" }},
		{"quote", func(p *Product) { p.QuoteCurrency = "INR" }},
		{"settlement", func(p *Product) { p.SettlementCurrency = "INR" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			before := fixture()
			after := fixture()
			after.AsOf = after.AsOf.Add(time.Second)
			tc.edit(&after.Products[0])
			changes, err := ClassifyChanges(before, after)
			if err != nil || len(changes) != 1 || changes[0].Kind != ContractChanged {
				t.Fatalf("contract field %s: changes=%+v err=%v", tc.name, changes, err)
			}
			if changes[0].ID != before.Products[0].ID {
				t.Fatalf("wrong product id: %+v", changes[0])
			}
		})
	}
	t.Run("status only", func(t *testing.T) {
		before := fixture()
		after := fixture()
		after.AsOf = after.AsOf.Add(time.Second)
		after.Products[0].Status = Inactive
		changes, err := ClassifyChanges(before, after)
		if err != nil || len(changes) != 1 || changes[0].Kind != ProductDisabled {
			t.Fatalf("status should not count as contract change: %+v err=%v", changes, err)
		}
	})
}

func TestOptionContractFieldClassificationParity(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*Product)
	}{
		{"strike", func(p *Product) { p.Strike = "999999.00001" }},
		{"expiry", func(p *Product) {
			expiry := p.Expiry.Add(time.Second)
			p.Expiry = &expiry
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := fixture()
			after := fixture()
			after.AsOf = after.AsOf.Add(time.Second)
			tc.edit(&after.Products[1])
			changes, err := ClassifyChanges(before, after)
			if err != nil || len(changes) != 1 || changes[0].Kind != ContractChanged {
				t.Fatalf("option contract %s: changes=%+v err=%v", tc.name, changes, err)
			}
		})
	}
}
