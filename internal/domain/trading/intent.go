package trading

import (
	"errors"
	"fmt"
	"strings"
)

// No floating-point financial values: external prices are parsed into fixed-point ticks.
// Later adapters must honor per-product tick size, lot multiplier and settlement currency.
type ProductKind string

const (
	Future ProductKind = "future"
	Option ProductKind = "option"
)

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type OrderIntent struct {
	ProductID       int64
	Kind            ProductKind
	Side            Side
	QuantityLots    int64
	LimitPriceTicks int64
	ClientOrderID   string
	ReduceOnly      bool
}

func (o OrderIntent) Validate() error {
	if o.ProductID <= 0 {
		return errors.New("product id must be positive")
	}
	if o.Kind != Future && o.Kind != Option {
		return fmt.Errorf("unsupported kind %q", o.Kind)
	}
	if o.Side != Buy && o.Side != Sell {
		return fmt.Errorf("unsupported side %q", o.Side)
	}
	if o.QuantityLots <= 0 || o.LimitPriceTicks <= 0 {
		return errors.New("quantity and limit price ticks must be positive")
	}
	if len(o.ClientOrderID) == 0 || len(o.ClientOrderID) > 32 || strings.TrimSpace(o.ClientOrderID) != o.ClientOrderID {
		return errors.New("client order id must be 1-32 characters without outer whitespace")
	}
	for _, c := range o.ClientOrderID {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return errors.New("client order id may only use ASCII letters, digits, _ and -")
		}
	}
	return nil
}
