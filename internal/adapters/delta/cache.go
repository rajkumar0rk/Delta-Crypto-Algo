package delta

import (
	"context"
	"reflect"
	"slices"
	"sort"
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

// PublicTermsChange records changed public terms for an existing product.
// Listed and removed products are covered separately by ProductChanges.
type PublicTermsChange struct {
	ProductID int64
	Before    PublicTerms
	After     PublicTerms
}

type RefreshResult struct {
	AsOf           time.Time
	ProductChanges []catalog.ClassifiedChange
	TermsChanges   []PublicTermsChange
}

func copyTerms(t PublicTerms) PublicTerms {
	t.Tags = slices.Clone(t.Tags)
	return t
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
		terms[id] = copyTerms(value)
	}
	return snap, terms, nil
}

// Refresh performs network I/O and validation outside the publication lock.
// An older concurrent observation cannot overwrite a newer one. A failed
// refresh keeps both the prior product snapshot and its financial terms.
func (c *CatalogCache) Refresh(ctx context.Context, adapter PublicCatalog, lots map[int64]LotEvidence, maxAge time.Duration) (RefreshResult, error) {
	if ctx == nil {
		return RefreshResult{}, ErrCatalogIncomplete
	}
	raw, err := adapter.Fetch(ctx)
	if err != nil {
		return RefreshResult{}, err
	}
	next, terms, err := MapSnapshot(raw, lots, maxAge)
	if err != nil {
		return RefreshResult{}, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return RefreshResult{}, err
	}
	changes, err := c.store.PublishWithChanges(next)
	if err != nil {
		return RefreshResult{}, err
	}
	result := RefreshResult{AsOf: next.AsOf, ProductChanges: changes}
	for id, after := range terms {
		if before, ok := c.terms[id]; ok && !reflect.DeepEqual(before, after) {
			result.TermsChanges = append(result.TermsChanges, PublicTermsChange{
				ProductID: id, Before: copyTerms(before), After: copyTerms(after),
			})
		}
	}
	sort.Slice(result.TermsChanges, func(i, j int) bool {
		return result.TermsChanges[i].ProductID < result.TermsChanges[j].ProductID
	})
	c.terms = terms
	return result, nil
}
