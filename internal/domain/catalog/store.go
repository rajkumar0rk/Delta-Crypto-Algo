package catalog

import (
	"errors"
	"fmt"
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
	if err := next.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrCatalogUnavailable, err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ready && !next.AsOf.After(s.current.AsOf) {
		return fmt.Errorf("%w: non-monotonic as_of", ErrCatalogUnavailable)
	}
	s.current = copySnapshot(next)
	s.ready = true
	return nil
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
