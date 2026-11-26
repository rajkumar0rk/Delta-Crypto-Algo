package catalog

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// AccountScope identifies an account and environment, never its credentials.
type AccountScope struct {
	AccountID   string
	Environment string
}

// EligibilityEvidence is independently observed account metadata. It is tied
// to a specific catalog observation so a contract change invalidates it.
type EligibilityEvidence struct {
	Scope         AccountScope
	ProductID     int64
	CatalogSource string
	CatalogAsOf   time.Time
	VerifiedAt    time.Time
	Source        string
	Allowed       bool
}

var ErrEligibilityUnverified = errors.New("account product eligibility unverified")

// CheckMetadataAdmission combines three independent prerequisites: current
// product metadata, verified product capabilities and scoped account evidence.
// It is deliberately not a risk decision, live authorization or order route.
func CheckMetadataAdmission(s Snapshot, capabilities CapabilitySet, evidence EligibilityEvidence, scope AccountScope, id int64, now time.Time, maxAge time.Duration) error {
	if strings.TrimSpace(scope.AccountID) == "" || strings.TrimSpace(scope.Environment) == "" || scope != evidence.Scope {
		return fmt.Errorf("%w: account scope mismatch", ErrEligibilityUnverified)
	}
	if err := s.Tradeable(id, now, maxAge); err != nil {
		return err
	}
	if !evidence.Allowed || evidence.ProductID != id || evidence.CatalogSource != s.Source || !evidence.CatalogAsOf.Equal(s.AsOf) || strings.TrimSpace(evidence.Source) == "" {
		return fmt.Errorf("%w: absent or mismatched product evidence", ErrEligibilityUnverified)
	}
	if evidence.VerifiedAt.IsZero() || evidence.VerifiedAt.Before(s.AsOf) || evidence.VerifiedAt.After(now) || now.Sub(evidence.VerifiedAt) > maxAge {
		return fmt.Errorf("%w: stale or invalid observation", ErrEligibilityUnverified)
	}
	for _, p := range s.Products {
		if p.ID != id {
			continue
		}
		capability, err := capabilities.VerifiedFor(p, now, maxAge)
		if err != nil {
			return err
		}
		if capability.VerifiedAt.Before(s.AsOf) {
			return fmt.Errorf("%w: capability predates current catalog", ErrCapabilityUnverified)
		}
		return nil
	}
	return ErrNotTradeable
}
