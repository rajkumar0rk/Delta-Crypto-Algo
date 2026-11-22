package catalog

import (
	"errors"
	"fmt"
	"time"
)

// ErrCapabilityUnverified is a metadata-only admission denial, not an order decision.
var ErrCapabilityUnverified = errors.New("venue product capability unverified")

// Capability records independently verified, product-specific exchange rules.
// Absence is deliberately different from a zero fee or zero margin requirement.
// No account permission, risk approval, or trading authorization is implied.
type Capability struct {
	ProductID int64
	VerifiedAt time.Time
	Source string
	AllowMarket bool
	AllowLimit bool
	AllowPostOnly bool
	AllowReduceOnly bool
	MakerFeeRate string
	TakerFeeRate string
	InitialMarginRate string
	MaintenanceMarginRate string
}

// CapabilitySet is intentionally keyed by numeric venue product ID rather than
// underlying symbol. Different expiries and option strikes have distinct rules.
type CapabilitySet struct {
	AsOf time.Time
	Source string
	Products map[int64]Capability
}

// VerifiedFor checks that the product's independently sourced fee/margin/order
// capabilities are complete and fresh. It does NOT grant trading permission.
func (s CapabilitySet) VerifiedFor(p Product, now time.Time, maxAge time.Duration) (Capability, error) {
	var denied Capability
	if err := p.Validate(); err != nil {
		return denied, fmt.Errorf("%w: invalid product: %v", ErrCapabilityUnverified, err)
	}
	if maxAge <= 0 || now.IsZero() || s.AsOf.IsZero() || now.Before(s.AsOf) || now.Sub(s.AsOf) > maxAge || s.Source == "" {
		return denied, fmt.Errorf("%w: missing or stale capability set", ErrCapabilityUnverified)
	}
	c, ok := s.Products[p.ID]
	if !ok || c.ProductID != p.ID || c.Source == "" || c.Source != s.Source || c.VerifiedAt.IsZero() || c.VerifiedAt.After(now) || c.VerifiedAt.After(s.AsOf) || s.AsOf.Sub(c.VerifiedAt) > maxAge {
		return denied, fmt.Errorf("%w: missing or stale product evidence", ErrCapabilityUnverified)
	}
	if !c.AllowMarket && !c.AllowLimit {
		return denied, fmt.Errorf("%w: no verified entry order type", ErrCapabilityUnverified)
	}
	if c.AllowPostOnly && !c.AllowLimit {
		return denied, fmt.Errorf("%w: post-only without limit capability", ErrCapabilityUnverified)
	}
	// These are ratios, not prices or percentages. They must be explicit exact
	// decimals; no defaults are inferred from product kind or underlying.
	for _, f := range []struct{name, value string}{
		{"maker fee", c.MakerFeeRate},
		{"taker fee", c.TakerFeeRate},
		{"initial margin", c.InitialMarginRate},
		{"maintenance margin", c.MaintenanceMarginRate},
	} {
		if !nonnegativeRatio(f.value) {
			return denied, fmt.Errorf("%w: invalid %s", ErrCapabilityUnverified, f.name)
		}
	}
	return c, nil
}

// nonnegativeRatio accepts bounded plain ASCII decimals only; this avoids
// floating point, scientific notation, signs, and unbounded parser inputs.
func nonnegativeRatio(v string) bool {
	if len(v) == 0 || len(v) > 64 {
		return false
	}
	dot := false
	digits := 0
	for i := 0; i < len(v); i++ {
		switch {
		case v[i] >= '0' && v[i] <= '9':
			digits++
		case v[i] == '.' && !dot:
			dot = true
		default:
			return false
		}
	}
	if digits == 0 || v[0] == '.' || v[len(v)-1] == '.' {
		return false
	}
	// Ratios are bounded to [0,1] without lossy float conversion.
	if v[0] != '0' {
		return v == "1" || (v[0] == '1' && dot && allZeros(v[2:]))
	}
	if dot {
		return len(v) >= 3 && allDecimalDigits(v[2:])
	}
	return allZeros(v)
}

func allZeros(v string) bool {
	for i := range v {
		if v[i] != '0' { return false }
	}
	return true
}

func allDecimalDigits(v string) bool {
	for i := range v {
		if v[i] < '0' || v[i] > '9' { return false }
	}
	return true
}
