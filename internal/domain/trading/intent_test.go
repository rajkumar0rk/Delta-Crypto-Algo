package trading

import "testing"

func TestIntentValidation(t *testing.T) {
	good := OrderIntent{ProductID: 1, Kind: Future, Side: Buy, QuantityLots: 1, LimitPriceTicks: 100, ClientOrderID: "id_123"}
	if err := good.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		mutate func(*OrderIntent)
	}{
		{"product", func(o *OrderIntent) { o.ProductID = 0 }},
		{"quantity", func(o *OrderIntent) { o.QuantityLots = 0 }},
		{"price", func(o *OrderIntent) { o.LimitPriceTicks = -1 }},
		{"type", func(o *OrderIntent) { o.Kind = "spot" }},
		{"side", func(o *OrderIntent) { o.Side = "hold" }},
		{"id", func(o *OrderIntent) { o.ClientOrderID = "unsafe whitespace" }},
		{"id_length", func(o *OrderIntent) { o.ClientOrderID = "123456789012345678901234567890123" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := good
			tc.mutate(&o)
			if err := o.Validate(); err == nil {
				t.Fatal("bad intent accepted")
			}
		})
	}
}
