package catalog

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

var ErrCatalogUnavailable = errors.New("validated catalog unavailable")

// Store publishes only complete, validated snapshots. Reads return defensive
// copies; neither callers nor failed refreshes can mutate the published data.
// This is metadata storage, NOT an exchange authorization or live trading gate.
type Store struct {
	mu sync.RWMutex

	current Snapshot

	ready bool
}

func copySnapshot(s Snapshot) Snapshot {
	out := s
	out.Products = make([]Product, len(s.Products))
	copy(out.Products, s.Products)
	for i := range out.Products {
		if s.Products[i].Expiry != nil {
			expiry := *s.Products[i].Expiry
			out.Products[i].Expiry = &expiry
		}
	}
	return out
}

// Publish atomically replaces a catalog after all validation has passed.
// Older or equally dated data is rejected to prevent refresh rollback.
// A failed refresh never replaces the last good snapshot; consumers MUST still
// check freshness on every use.
func (s *Store) Publish(next Snapshot) error {
	_, err := s.PublishWithChanges(next)
	return err
}

// PublishWithChanges computes the audit transition and publishes under the
// same lock. An invalid, older or cross-source refresh changes neither state
// nor audit result. Initial publication reports every product as listed.
func (s *Store) PublishWithChanges(next Snapshot) ([]ClassifiedChange, error) {
	if err := next.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCatalogUnavailable, err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var changes []ClassifiedChange
	if s.ready {
		var err error
		changes, err = ClassifyChanges(s.current, next)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCatalogUnavailable, err)
		}
	} else {
		changes = make([]ClassifiedChange, 0, len(next.Products))
		for _, p := range next.Products {
			changes = append(changes, ClassifiedChange{Change: Change{ID: p.ID, After: copyProduct(p)}, Kind: ProductListed})
		}
		sort.Slice(changes, func(i, j int) bool { return changes[i].ID < changes[j].ID })
	}
	s.current = copySnapshot(next)
	s.ready = true
	return changes, nil
}

// Snapshot returns an independent copy of the latest complete snapshot.
// A stale catalog remains inspectable for diagnostics, not automatically safe.
func (s *Store) Snapshot() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ready {
		return Snapshot{}, ErrCatalogUnavailable
	}
	return copySnapshot(s.current), nil
}

// Tradeable uses the same fail-closed metadata precondition as Snapshot.Tradeable.
// It does not grant account permissions, margin or order execution authorization.
func (s *Store) Tradeable(id int64, now time.Time, maxAge time.Duration) error {
	snap, err := s.Snapshot()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotTradeable, err)
	}
	return snap.Tradeable(id, now, maxAge)
}
