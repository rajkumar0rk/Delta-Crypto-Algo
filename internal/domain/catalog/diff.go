package catalog

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
)

// Change describes a catalog transition without interpreting it as permission to trade.
// Before and After are detached copies so callers cannot mutate source snapshots.
type Change struct {
	ID     int64
	Before *Product
	After  *Product
}

// Diff compares two complete, validated catalog snapshots. It rejects non-monotonic
// observation times and preserves every removed product as an explicit transition.
// Changes are ordered by product ID to keep audit logs and tests deterministic.
func Diff(before, after Snapshot) ([]Change, error) {
	if err := before.Validate(); err != nil {
		return nil, fmt.Errorf("before snapshot: %w", err)
	}
	if err := after.Validate(); err != nil {
		return nil, fmt.Errorf("after snapshot: %w", err)
	}
	if before.Source != after.Source {
		return nil, errors.New("cannot compare different catalog sources")
	}
	if !after.AsOf.After(before.AsOf) {
		return nil, errors.New("catalog observation time must advance")
	}
	old := make(map[int64]Product, len(before.Products))
	current := make(map[int64]Product, len(after.Products))
	for _, p := range before.Products {
		old[p.ID] = p
	}
	for _, p := range after.Products {
		current[p.ID] = p
	}
	ids := make([]int64, 0, len(old)+len(current))
	for id := range old {
		ids = append(ids, id)
	}
	for id := range current {
		if _, ok := old[id]; !ok {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	changes := make([]Change, 0)
	for _, id := range ids {
		prev, had := old[id]
		next, has := current[id]
		if had && has && sameProduct(prev, next) {
			continue
		}
		change := Change{ID: id}
		if had {
			change.Before = copyProduct(prev)
		}
		if has {
			change.After = copyProduct(next)
		}
		changes = append(changes, change)
	}
	return changes, nil
}

func copyProduct(p Product) *Product {
	if p.Expiry != nil {
		expiry := *p.Expiry
		p.Expiry = &expiry
	}
	return &p
}

// sameProduct compares exchange metadata semantically. time.Time's location
// and monotonic internals are not contract changes when instants are equal.
func sameProduct(a, b Product) bool {
	return a.ID == b.ID && sameContractSpec(a, b) && a.Status == b.Status
}

// sameContractSpec is shared by diffing and classification to prevent drift
// when venue contract metadata evolves. Status is deliberately excluded.
func sameContractSpec(a, b Product) bool {
	return a.Symbol == b.Symbol &&
		a.Underlying == b.Underlying && a.Kind == b.Kind &&
		sameDecimal(a.TickSize, b.TickSize) &&
		sameDecimal(a.LotSize, b.LotSize) && sameDecimal(a.ContractMultiplier, b.ContractMultiplier) &&
		a.QuoteCurrency == b.QuoteCurrency &&
		a.SettlementCurrency == b.SettlementCurrency &&
		sameOptionalDecimal(a.Strike, b.Strike) && sameExpiry(a.Expiry, b.Expiry)
}

// sameDecimal compares validated exact decimal strings by numerical value,
// not display scale (e.g. 0.10 and 0.100 are the same contract tick).
func sameDecimal(a, b string) bool {
	if a == b {
		return true
	}
	x, ok := new(big.Rat).SetString(a)
	if !ok {
		return false
	}
	y, ok := new(big.Rat).SetString(b)
	return ok && x.Cmp(y) == 0
}

func sameOptionalDecimal(a, b string) bool {
	if a == "" || b == "" {
		return a == b
	}
	return sameDecimal(a, b)
}
