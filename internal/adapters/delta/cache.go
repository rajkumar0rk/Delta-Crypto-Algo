package delta

import (
	"context"
	"sync"
	"time"

	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/domain/catalog"
)

// CatalogCache keeps normalized products and raw financial terms from the
// same observation. The zero value is empty, never implicitly tradeable.
type CatalogCache struct {
	mu    sync.RWMutex
	store catalog.Store
	terms map[int64]PublicTerms
}

// Read returns detached product and terms copies from one publication.
// Consumers must still enforce snapshot freshness at the time of use.
func (c *CatalogCache) Read() (catalog.Snapshot, map[int64]PublicTerms, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snap, err := c.store.Snapshot()
	if err != nil {
		return catalog.Snapshot{}, nil, err
	}
	terms := make(map[int64]PublicTerms, len(c.terms))
	for id, value := range c.terms {
		terms[id] = value
	}
	return snap, terms, nil
}

// Refresh performs network I/O and validation outside the publication lock.
// An older concurrent observation cannot overwrite a newer one. A failed
// refresh keeps both the prior product snapshot and its financial terms.
func (c *CatalogCache) Refresh(ctx context.Context, adapter PublicCatalog, lots map[int64]LotEvidence, maxAge time.Duration) ([]catalog.ClassifiedChange, error) {
	if ctx == nil {
		return nil, ErrCatalogIncomplete
	}
	raw, err := adapter.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	next, terms, err := MapSnapshot(raw, lots, maxAge)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	changes, err := c.store.PublishWithChanges(next)
	if err != nil {
		return nil, err
	}
	c.terms = terms
	return changes, nil
}
