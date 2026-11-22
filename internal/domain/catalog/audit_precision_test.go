package catalog

import (
	"testing"
	"time"
)

// Audit evidence must be detached all the way through nested expiry pointers,
// even when the caller mutates the returned record after classification.
func TestClassifiedChangeExpiryEvidenceIsDetached(t *testing.T) {
	before := fixture()
	after := fixture()
	after.AsOf = before.AsOf.Add(time.Second)
	changed := after.Products[1].Expiry.Add(time.Hour)
	after.Products[1].Expiry = &changed

	changes, err := ClassifyChanges(before, after)
	if err != nil || len(changes) != 1 || changes[0].Kind != ContractChanged {
		t.Fatalf("unexpected expiry audit: changes=%+v err=%v", changes, err)
	}
	*changes[0].Before.Expiry = changes[0].Before.Expiry.Add(48 * time.Hour)
	*changes[0].After.Expiry = changes[0].After.Expiry.Add(48 * time.Hour)
	if !before.Products[1].Expiry.Equal(before.AsOf.Add(24 * time.Hour)) {
		t.Fatal("before snapshot expiry was mutated through audit evidence")
	}
	if !after.Products[1].Expiry.Equal(changed) {
		t.Fatal("after snapshot expiry was mutated through audit evidence")
	}
}

// Equivalent expiry instants with different timezone representations are not
// contract changes. A genuine one-nanosecond shift remains observable.
func TestExpiryComparisonUsesInstantNotTimezone(t *testing.T) {
	before := fixture()
	after := fixture()
	after.AsOf = before.AsOf.Add(time.Second)
	sameInstant := after.Products[1].Expiry.In(time.FixedZone("offset", 5*3600+30*60))
	after.Products[1].Expiry = &sameInstant
	changes, err := ClassifyChanges(before, after)
	if err != nil || len(changes) != 0 {
		t.Fatalf("timezone-only change: changes=%+v err=%v", changes, err)
	}
	differentInstant := sameInstant.Add(time.Nanosecond)
	after.Products[1].Expiry = &differentInstant
	changes, err = ClassifyChanges(before, after)
	if err != nil || len(changes) != 1 || changes[0].Kind != ContractChanged {
		t.Fatalf("nanosecond expiry change: changes=%+v err=%v", changes, err)
	}
}

// Decimal formatting changes must not obscure a real status transition;
// very small actual changes must not disappear through float rounding.
func TestHighPrecisionMultiplierAndStatusAudit(t *testing.T) {
	before := fixture()
	after := fixture()
	after.AsOf = before.AsOf.Add(time.Second)
	after.Products[1].ContractMultiplier = "0.001000000000000000000000000000"
	after.Products[1].Status = Inactive
	changes, err := ClassifyChanges(before, after)
	if err != nil || len(changes) != 1 || changes[0].Kind != ProductDisabled {
		t.Fatalf("equivalent multiplier: changes=%+v err=%v", changes, err)
	}
	after.Products[1].ContractMultiplier = "0.001000000000000000000000000001"
	changes, err = ClassifyChanges(before, after)
	if err != nil || len(changes) != 1 || changes[0].Kind != ContractChanged {
		t.Fatalf("high precision contract change: changes=%+v err=%v", changes, err)
	}
}
