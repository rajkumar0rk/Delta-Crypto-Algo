package delta

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/domain/catalog"
)

// Synthetic fixtures follow the documented Product field names, not captured
// exchange responses. No venue lot size or fee schedule is asserted.
const productFixture = `{"id":42,"symbol":"SYNTH","contract_type":"perpetual_futures","state":"live","trading_status":"operational","tick_size":"0.05","contract_value":"0.001","contract_unit_currency":"SYNTH","notional_type":"vanilla","is_quanto":false,"settlement_time":null,"strike_price":null,"initial_margin":"0.5","maintenance_margin":"0.25","maker_commission_rate":"0.0002","taker_commission_rate":"0.0005","underlying_asset":{"symbol":"SYNTH"},"quoting_asset":{"symbol":"USD"},"settling_asset":{"symbol":"INR"},"product_specs":{"only_reduce_only_orders_allowed":false}}`

func fixtureLots(now time.Time) map[int64]LotEvidence {
	return map[int64]LotEvidence{42: {ProductID: 42, Symbol: "SYNTH", ContractValue: "0.001", LotSize: "1", Source: "synthetic-test", VerifiedAt: now}}
}

func TestMapSnapshotDocumentedFields(t *testing.T) {
	now := time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		raw    string
		kind   catalog.Kind
		status catalog.Status
	}{
		{"perpetual", productFixture, catalog.Perpetual, catalog.Active},
		{"RWA underlying", strings.ReplaceAll(productFixture, "SYNTH", "RWA"), catalog.Perpetual, catalog.Active},
		{"reduce only", strings.Replace(productFixture, "allowed\":false", "allowed\":true", 1), catalog.Perpetual, catalog.Inactive},
		{"unknown trading status", strings.Replace(productFixture, "operational", "unknown", 1), catalog.Perpetual, catalog.Inactive},
		{"dated future", strings.Replace(strings.Replace(productFixture, "perpetual_futures", "futures", 1), "settlement_time\":null", "settlement_time\":\"2026-12-01T00:00:00Z\"", 1), catalog.DatedFuture, catalog.Active},
		{"call", strings.Replace(strings.Replace(strings.Replace(productFixture, "perpetual_futures", "call_options", 1), "settlement_time\":null", "settlement_time\":\"2026-12-01T00:00:00Z\"", 1), "strike_price\":null", "strike_price\":\"100\"", 1), catalog.CallOption, catalog.Active},
		{"put expired", strings.Replace(strings.Replace(strings.Replace(productFixture, "perpetual_futures", "put_options", 1), "settlement_time\":null", "settlement_time\":\"2026-11-01T00:00:00Z\"", 1), "strike_price\":null", "strike_price\":\"100\"", 1), catalog.PutOption, catalog.Expired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lots := fixtureLots(now)
			if tc.name == "RWA underlying" {
				lot := lots[42]
				lot.Symbol = "RWA"
				lots[42] = lot
			}
			s, terms, err := MapSnapshot(RawSnapshot{AsOf: now, Pages: 1, Products: []json.RawMessage{json.RawMessage(tc.raw)}}, lots, time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			p := s.Products[0]
			if p.Kind != tc.kind || p.Status != tc.status || p.TickSize != "0.05" || p.ContractMultiplier != "0.001" || p.SettlementCurrency != "INR" || terms[42].InitialMargin != "0.5" {
				t.Fatalf("unexpected mapping: %+v %+v", p, terms)
			}
		})
	}
}

func TestMapSnapshotDeniesUnverifiedRecords(t *testing.T) {
	now := time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)
	for _, raw := range []string{
		strings.Replace(productFixture, "vanilla", "inverse", 1),
		strings.Replace(productFixture, "perpetual_futures", "unknown", 1),
		strings.Replace(productFixture, "0.05", "1/2", 1),
		strings.Replace(productFixture, "0.001", "0.002", 1),
	} {
		s, terms, err := MapSnapshot(RawSnapshot{AsOf: now, Pages: 1, Products: []json.RawMessage{json.RawMessage(productFixture), json.RawMessage(raw)}}, fixtureLots(now), time.Minute)
		if err == nil || len(s.Products) != 0 || terms != nil {
			t.Fatalf("partial or unverified snapshot published: %+v %v", s, err)
		}
	}
	s, terms, err := MapSnapshot(RawSnapshot{AsOf: now, Pages: 1, Products: []json.RawMessage{json.RawMessage(productFixture)}}, nil, time.Minute)
	if err == nil || len(s.Products) != 0 || terms != nil {
		t.Fatal("missing lot evidence accepted")
	}
}

func TestRefreshCatalogPreservesLastGoodSnapshot(t *testing.T) {
	var store catalog.Store
	body := productFixture
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"result":[` + body + `],"meta":{"after":null}}`))
	}))
	defer srv.Close()
	client := srv.Client()
	client.Timeout = time.Second
	adapter := PublicCatalog{BaseURL: srv.URL, Client: client}
	lots := fixtureLots(time.Now().Add(-time.Second))
	changes, _, err := adapter.RefreshCatalog(context.Background(), &store, lots, time.Minute)
	if err != nil || len(changes) != 1 || changes[0].Kind != catalog.ProductListed {
		t.Fatalf("initial refresh: %+v %v", changes, err)
	}
	before, err := store.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	body = strings.Replace(productFixture, "0.05", "bad", 1)
	if _, _, err := adapter.RefreshCatalog(context.Background(), &store, lots, time.Minute); err == nil {
		t.Fatal("invalid refresh accepted")
	}
	after, err := store.Snapshot()
	if err != nil || !after.AsOf.Equal(before.AsOf) || after.Products[0].TickSize != before.Products[0].TickSize {
		t.Fatalf("last valid snapshot lost: %+v %v", after, err)
	}
}
