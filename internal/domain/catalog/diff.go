package catalog

import (
	"errors"
	"fmt"
	"reflect"
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
		if had && has && reflect.DeepEqual(prev, next) {
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
