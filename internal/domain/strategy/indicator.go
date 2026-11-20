// Package strategy provides server-owned, draft-only semantic validation primitives.
// Validation here does not authorize a strategy for execution.
package strategy

import (
	"errors"
	"fmt"
)

var ErrInvalidIndicator = errors.New("unsupported indicator configuration")

// Indicator describes an untrusted strategy operand; it is never executable code.
type Indicator struct {
	ID string
	Output string
	Parameters map[string]int
}

// Supported v1 indicator outputs are explicit and independent of client-provided JSON.
// Period semantics are deliberately narrow until a versioned parameter registry exists.
var outputs = map[string]map[string]struct{}{
	"ema": {"value": {}},
	"sma": {"value": {}},
	"rsi": {"value": {}},
	"atr": {"value": {}},
	"vwap": {"value": {}},
	"supertrend": {"trend": {}, "band": {}},
	"macd": {"line": {}, "signal": {}, "histogram": {}},
	"stoch_rsi": {"k": {}, "d": {}},
	"adx": {"adx": {}, "plus_di": {}, "minus_di": {}},
	"volume_ratio": {"value": {}},
	"bollinger": {"upper": {}, "middle": {}, "lower": {}},
	"donchian": {"upper": {}, "lower": {}, "middle": {}},
}

// ValidateIndicator rejects unknown IDs, outputs and parameter shapes.
// This is a preliminary semantic gate, not full schema/type/capability validation.
func ValidateIndicator(i Indicator) error {
	allowed, ok := outputs[i.ID]
	if !ok {
		return fmt.Errorf("%w: unknown ID %q", ErrInvalidIndicator, i.ID)
	}
	if _, ok := allowed[i.Output]; !ok {
		return fmt.Errorf("%w: unsupported output %q for %q", ErrInvalidIndicator, i.Output, i.ID)
	}
	if len(i.Parameters) != 1 {
		return fmt.Errorf("%w: exactly one period parameter required", ErrInvalidIndicator)
	}
	period, ok := i.Parameters["period"]
	if !ok {
		return fmt.Errorf("%w: period parameter required", ErrInvalidIndicator)
	}
	if i.ID == "vwap" {
		if period != 0 {
			return fmt.Errorf("%w: vwap period must be zero", ErrInvalidIndicator)
		}
		return nil
	}
	if period < 1 || period > 10000 {
		return fmt.Errorf("%w: period must be 1..10000", ErrInvalidIndicator)
	}
	return nil
}
