// Package catalog validates exchange product metadata without enabling trading.
package catalog

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Kind string

const (
	Perpetual   Kind = "perpetual_future"
	DatedFuture Kind = "dated_future"
	CallOption  Kind = "call_option"
	PutOption   Kind = "put_option"
)

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Expired  Status = "expired"
	Delisted Status = "delisted"
	Unknown  Status = "unknown"
)

// Product records venue-provided product properties. Decimal values are exact strings.
type Product struct {
	ID                 int64      `json:"id"`
	Symbol             string     `json:"symbol"`
	Underlying         string     `json:"underlying"`
	Kind               Kind       `json:"kind"`
	Status             Status     `json:"status"`
	TickSize           string     `json:"tick_size"`
	LotSize            string     `json:"lot_size"`
	ContractMultiplier string     `json:"contract_multiplier"`
	QuoteCurrency      string     `json:"quote_currency"`
	SettlementCurrency string     `json:"settlement_currency"`
	Expiry             *time.Time `json:"expiry,omitempty"`
	Strike             string     `json:"strike,omitempty"`
}

type Snapshot struct {
	Source   string    `json:"source"`
	AsOf     time.Time `json:"as_of"`
	Complete bool      `json:"complete"`
	Products []Product `json:"products"`
}

var ErrNotTradeable = errors.New("product not confirmed tradeable")

// positiveDecimal accepts only plain base-10 decimals. big.Rat.SetString also
// accepts non-venue representations (such as fractions and base-prefixed
// integers), so validate the lexical form before parsing exactly.
func positiveDecimal(s string) bool {
	if s == "" || len(s) > 128 {
		return false
	}
	dot := false
	digits := 0
	for _, ch := range s {
		switch {
		case ch >= '0' && ch <= '9':
			digits++
		case ch == '.' && !dot:
			dot = true
		default:
			return false
		}
	}
	if digits == 0 || s[0] == '.' || s[len(s)-1] == '.' {
		return false
	}
	r, ok := new(big.Rat).SetString(s)
	return ok && r.Sign() > 0
}

func (p Product) Validate() error {
	if p.ID <= 0 {
		return errors.New("product ID must be positive")
	}
	if len(p.Symbol) == 0 || len(p.Symbol) > 128 || strings.TrimSpace(p.Symbol) != p.Symbol {
		return errors.New("invalid symbol")
	}
	if len(p.Underlying) == 0 || len(p.Underlying) > 128 || strings.TrimSpace(p.Underlying) != p.Underlying {
		return errors.New("invalid underlying")
	}
	switch p.Kind {
	case Perpetual, DatedFuture, CallOption, PutOption:
	default:
		return fmt.Errorf("unsupported kind %q", p.Kind)
	}
	switch p.Status {
	case Active, Inactive, Expired, Delisted, Unknown:
	default:
		return fmt.Errorf("unsupported status %q", p.Status)
	}
	for _, f := range []struct{ name, value string }{
		{"tick_size", p.TickSize},
		{"lot_size", p.LotSize},
		{"contract_multiplier", p.ContractMultiplier},
	} {
		if !positiveDecimal(f.value) {
			return fmt.Errorf("%s must be a positive exact decimal", f.name)
		}
	}
	for _, f := range []struct{ name, value string }{
		{"quote_currency", p.QuoteCurrency},
		{"settlement_currency", p.SettlementCurrency},
	} {
		if len(f.value) == 0 || len(f.value) > 32 || strings.TrimSpace(f.value) != f.value {
			return fmt.Errorf("invalid %s", f.name)
		}
	}
	if p.Kind == CallOption || p.Kind == PutOption {
		if p.Expiry == nil || p.Expiry.IsZero() || !positiveDecimal(p.Strike) {
			return errors.New("option requires expiry and positive strike")
		}
	} else if p.Strike != "" {
		return errors.New("futures cannot specify option strike")
	}
	if p.Kind == DatedFuture && (p.Expiry == nil || p.Expiry.IsZero()) {
		return errors.New("dated future requires expiry")
	}
	if p.Kind == Perpetual && p.Expiry != nil {
		return errors.New("perpetual future cannot expire")
	}
	return nil
}

func (s Snapshot) Validate() error {
	if s.Source == "" || len(s.Source) > 128 || strings.TrimSpace(s.Source) != s.Source {
		return errors.New("invalid source")
	}
	if s.AsOf.IsZero() {
		return errors.New("snapshot missing as_of")
	}
	if !s.Complete {
		return errors.New("incomplete product pagination")
	}
	if len(s.Products) == 0 {
		return errors.New("empty product catalog")
	}
	seenID := make(map[int64]struct{}, len(s.Products))
	seenSymbol := make(map[string]struct{}, len(s.Products))
	for i, p := range s.Products {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("products[%d]: %w", i, err)
		}
		if _, ok := seenID[p.ID]; ok {
			return fmt.Errorf("duplicate product id %d", p.ID)
		}
		seenID[p.ID] = struct{}{}
		if _, ok := seenSymbol[p.Symbol]; ok {
			return fmt.Errorf("duplicate product symbol %q", p.Symbol)
		}
		seenSymbol[p.Symbol] = struct{}{}
	}
	return nil
}

// Tradeable checks metadata only. Account permissions and risk checks are separate.
func (s Snapshot) Tradeable(productID int64, now time.Time, maxAge time.Duration) error {
	if err := s.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrNotTradeable, err)
	}
	if maxAge <= 0 || now.Before(s.AsOf) || now.Sub(s.AsOf) > maxAge {
		return fmt.Errorf("%w: stale or future snapshot", ErrNotTradeable)
	}
	for _, p := range s.Products {
		if p.ID != productID {
			continue
		}
		if p.Status != Active || (p.Expiry != nil && !p.Expiry.After(now)) {
			return fmt.Errorf("%w: inactive or expired product", ErrNotTradeable)
		}
		return nil
	}
	return fmt.Errorf("%w: product missing", ErrNotTradeable)
}
