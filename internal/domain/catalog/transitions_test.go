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
