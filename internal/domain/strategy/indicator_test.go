package strategy

import (
	"errors"
	"testing"
)

func TestValidateIndicator(t *testing.T) {
	for _, tc := range []struct {
		name string
		indicator Indicator
		valid bool
	}{
		{"EMA", Indicator{"ema", "value", map[string]int{"period": 20}}, true},
		{"VWAP", Indicator{"vwap", "value", map[string]int{"period": 0}}, true},
		{"Supertrend band", Indicator{"supertrend", "band", map[string]int{"period": 10}}, true},
		{"MACD histogram", Indicator{"macd", "histogram", map[string]int{"period": 12}}, true},
		{"Bollinger lower", Indicator{"bollinger", "lower", map[string]int{"period": 20}}, true},
		{"unknown", Indicator{"custom_eval", "value", map[string]int{"period": 20}}, false},
		{"prototype name", Indicator{"__proto__", "value", map[string]int{"period": 20}}, false},
		{"inherited name", Indicator{"toString", "value", map[string]int{"period": 20}}, false},
		{"wrong output", Indicator{"ema", "upper", map[string]int{"period": 20}}, false},
		{"empty output", Indicator{"ema", "", map[string]int{"period": 20}}, false},
		{"missing params", Indicator{"ema", "value", nil}, false},
		{"wrong param", Indicator{"ema", "value", map[string]int{"length": 20}}, false},
		{"extra param", Indicator{"ema", "value", map[string]int{"period": 20, "eval": 1}}, false},
		{"negative", Indicator{"ema", "value", map[string]int{"period": -1}}, false},
		{"zero", Indicator{"ema", "value", map[string]int{"period": 0}}, false},
		{"too large", Indicator{"ema", "value", map[string]int{"period": 10001}}, false},
		{"vwap positive", Indicator{"vwap", "value", map[string]int{"period": 20}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateIndicator(tc.indicator)
			if tc.valid && err != nil {
				t.Fatalf("valid indicator rejected: %v", err)
			}
			if !tc.valid && !errors.Is(err, ErrInvalidIndicator) {
				t.Fatalf("invalid indicator accepted or wrong error: %v", err)
			}
		})
	}
}
