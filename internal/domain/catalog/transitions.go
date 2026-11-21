package catalog

import "time"

// TransitionKind describes the meaning of a validated snapshot change for
// operators. It is metadata only, never a trading or order authorization.
type TransitionKind string

const (
	ProductListed	TransitionKind = "listed"
	ProductRemoved	TransitionKind = "removed"
	ProductDisabled	TransitionKind = "disabled"
	ProductEnabled	TransitionKind = "enabled"
	ContractChanged	TransitionKind = "contract_changed"
	ProductUpdated	TransitionKind = "updated"
)

// ClassifiedChange keeps the complete before/after evidence, not just a label.
type ClassifiedChange struct {
	Change
	Kind TransitionKind
}

// ClassifyChanges preserves Diff's stable ID ordering and propagates all
// validation errors. Structural changes outrank status-only transitions so
// operators cannot mistake a changed contract for a routine reactivation.
func ClassifyChanges(before, after Snapshot) ([]ClassifiedChange, error) {
	changes, err := Diff(before, after)
	if err != nil {
		return nil, err
	}
	result := make([]ClassifiedChange, 0, len(changes))
	for _, change := range changes {
		result = append(result, ClassifiedChange{Change: change, Kind: classify(change)})
	}
	return result, nil
}

func classify(c Change) TransitionKind {
	if c.Before == nil {
		return ProductListed
	}
	if c.After == nil {
		return ProductRemoved
	}
	a, b := c.Before, c.After
	if a.Symbol != b.Symbol || a.Underlying != b.Underlying ||
		a.Kind != b.Kind || a.TickSize != b.TickSize ||
		a.LotSize != b.LotSize || a.ContractMultiplier != b.ContractMultiplier ||
		a.QuoteCurrency != b.QuoteCurrency ||
		a.SettlementCurrency != b.SettlementCurrency ||
		a.Strike != b.Strike || !sameExpiry(a.Expiry, b.Expiry) {
		return ContractChanged
	}
	if a.Status == Active && b.Status != Active {
		return ProductDisabled
	}
	if a.Status != Active && b.Status == Active {
		return ProductEnabled
	}
	return ProductUpdated
}

func sameExpiry(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Equal(*b)
}
