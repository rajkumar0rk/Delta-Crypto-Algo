package delta

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/domain/catalog"
)

// Captured public fields are genuine observations. Lot evidence is synthetic
// and deliberately test-only: the samples do not establish a size step.
func TestMapCapturedIndiaProductFields(t *testing.T) {
	data, err := os.ReadFile("testdata/india_public_samples_20261124.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		ObservedAt time.Time         `json:"observed_at"`
		Products   []json.RawMessage `json:"products"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	if len(fixture.Products) != 6 || fixture.ObservedAt.IsZero() {
		t.Fatal("invalid captured sample fixture")
	}
	lots := make(map[int64]LotEvidence)
	for _, row := range fixture.Products {
		var fields struct {
			ID            int64  `json:"id"`
			Symbol        string `json:"symbol"`
			ContractValue string `json:"contract_value"`
		}
		if err := json.Unmarshal(row, &fields); err != nil {
			t.Fatal(err)
		}
		lots[fields.ID] = LotEvidence{ProductID: fields.ID, Symbol: fields.Symbol, ContractValue: fields.ContractValue, LotSize: "1", Source: "synthetic-test-only-lot", VerifiedAt: fixture.ObservedAt}
	}
	raw := RawSnapshot{AsOf: fixture.ObservedAt, Pages: 1, Products: fixture.Products}
	if _, _, err := MapSnapshot(raw, nil, time.Minute); err == nil {
		t.Fatal("public sample admitted without lot evidence")
	}
	// Treat the sample as a test input only; it is not a full venue snapshot.
	snapshot, terms, err := MapSnapshot(raw, lots, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	kinds := make(map[catalog.Kind]int)
	for _, p := range snapshot.Products {
		kinds[p.Kind]++
		if p.Status != catalog.Active || p.TickSize != "0.01" || p.SettlementCurrency != "USD" {
			t.Fatalf("unexpected mapped product: %+v", p)
		}
		v := terms[p.ID]
		if p.Kind == catalog.Perpetual {
			if !slices.Contains(v.Tags, "xStock") || v.InitialMargin != "4" || v.MakerCommission != "0.0002" || v.PriceBand != "10" || v.PositionSizeLimit != "80000" {
				t.Fatalf("xStock public terms: %+v", v)
			}
		} else {
			if p.Expiry == nil || p.Strike == "" || v.InitialMargin != "1" || v.MakerCommission != "0.0001" || v.PositionSizeLimit != "120000" {
				t.Fatalf("option public terms: %+v %+v", p, v)
			}
		}
	}
	if kinds[catalog.Perpetual] != 2 || kinds[catalog.CallOption] != 2 || kinds[catalog.PutOption] != 2 {
		t.Fatalf("unexpected product types: %+v", kinds)
	}
}
