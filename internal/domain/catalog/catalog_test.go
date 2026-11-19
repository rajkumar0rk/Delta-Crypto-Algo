package catalog

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func fixture() Snapshot {
	now := time.Date(2026, 11, 19, 9, 0, 0, 0, time.UTC)
	expiry := now.Add(24 * time.Hour)
	return Snapshot{
		Source:   "delta-india",
		AsOf:     now,
		Complete: true,
		Products: []Product{
			{ID: 10, Symbol: "ANY-RWA-PERP", Underlying: "ANY-RWA", Kind: Perpetual, Status: Active, TickSize: "0.00000001", LotSize: "0.01", ContractMultiplier: "1", QuoteCurrency: "USD", SettlementCurrency: "USD"},
			{ID: 11, Symbol: "ALT-C-1", Underlying: "ALT", Kind: CallOption, Status: Active, TickSize: "0.5", LotSize: "1", ContractMultiplier: "0.001", QuoteCurrency: "USD", SettlementCurrency: "USD", Strike: "100.25", Expiry: &expiry},
		},
	}
}

func TestSnapshotAndTradeability(t *testing.T) {
	s := fixture()
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := s.Tradeable(10, s.AsOf.Add(time.Minute), time.Hour); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		mutate func(*Snapshot)
	}{
		{"partial pagination", func(s *Snapshot) { s.Complete = false }},
		{"empty catalog", func(s *Snapshot) { s.Products = nil }},
		{"duplicate ID", func(s *Snapshot) { s.Products[1].ID = s.Products[0].ID }},
		{"duplicate symbol", func(s *Snapshot) { s.Products[1].Symbol = s.Products[0].Symbol }},
		{"invalid tick", func(s *Snapshot) { s.Products[0].TickSize = "0" }},
		{"exponent rejected", func(s *Snapshot) { s.Products[0].LotSize = "1e-8" }},
		{"negative multiplier", func(s *Snapshot) { s.Products[0].ContractMultiplier = "-1" }},
		{"hexadecimal quantity", func(s *Snapshot) { s.Products[0].LotSize = "0x10" }},
		{"binary quantity", func(s *Snapshot) { s.Products[0].LotSize = "0b10" }},
		{"underscored quantity", func(s *Snapshot) { s.Products[0].LotSize = "1_000" }},
		{"signed quantity", func(s *Snapshot) { s.Products[0].LotSize = "+1" }},
		{"leading decimal point", func(s *Snapshot) { s.Products[0].LotSize = ".5" }},
		{"trailing decimal point", func(s *Snapshot) { s.Products[0].LotSize = "5." }},
		{"multiple decimal points", func(s *Snapshot) { s.Products[0].LotSize = "1.2.3" }},
		{"overlong decimal", func(s *Snapshot) { s.Products[0].LotSize = strings.Repeat("1", 129) }},
		{"missing option expiry", func(s *Snapshot) { s.Products[1].Expiry = nil }},
		{"missing option strike", func(s *Snapshot) { s.Products[1].Strike = "" }},
		{"unknown product type", func(s *Snapshot) { s.Products[0].Kind = "mystery" }},
		{"unknown status", func(s *Snapshot) { s.Products[0].Status = "maybe" }},
		{"future has strike", func(s *Snapshot) { s.Products[0].Strike = "1" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := fixture()
			tc.mutate(&s)
			if err := s.Validate(); err == nil {
				t.Fatal("invalid snapshot accepted")
			}
		})
	}
	for _, tc := range []struct {
		name   string
		id     int64
		now    time.Time
		age    time.Duration
		mutate func(*Snapshot)
	}{
		{"unknown ID", 99, s.AsOf.Add(time.Minute), time.Hour, func(*Snapshot) {}},
		{"stale", 10, s.AsOf.Add(2 * time.Hour), time.Hour, func(*Snapshot) {}},
		{"future dated snapshot", 10, s.AsOf.Add(-time.Minute), time.Hour, func(*Snapshot) {}},
		{"zero freshness budget", 10, s.AsOf.Add(time.Minute), 0, func(*Snapshot) {}},
		{"inactive", 10, s.AsOf.Add(time.Minute), time.Hour, func(s *Snapshot) { s.Products[0].Status = Inactive }},
		{"expired option", 11, s.Products[1].Expiry.Add(time.Second), 48 * time.Hour, func(*Snapshot) {}},
		{"partial", 10, s.AsOf.Add(time.Minute), time.Hour, func(s *Snapshot) { s.Complete = false }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := fixture()
			tc.mutate(&s)
			if err := s.Tradeable(tc.id, tc.now, tc.age); !errors.Is(err, ErrNotTradeable) {
				t.Fatalf("wanted fail-closed: %v", err)
			}
		})
	}
}

func TestProductKinds(t *testing.T) {
	s := fixture()
	expiry := s.Products[1].Expiry
	p := s.Products[0]
	p.Kind = DatedFuture
	p.Expiry = expiry
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	p.Kind = PutOption
	p.Strike = "42.5"
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	p.Kind = Perpetual
	if err := p.Validate(); err == nil {
		t.Fatal("perpetual with expiry accepted")
	}
}

func TestPositiveDecimalPlainBase10(t *testing.T) {
	for _, value := range []string{"1", "0.00000001", "001.2500", "1000000000000000000000000"} {
		if !positiveDecimal(value) {
			t.Errorf("rejected valid decimal %q", value)
		}
	}
	for _, value := range []string{"", "0", "0.000", "0x10", "0b10", "1_000", "1/2", "1e3", "+1", "-1", ".5", "5.", "1.2.3", "NaN", "Inf", " 1", "1 ", "１２"} {
		if positiveDecimal(value) {
			t.Errorf("accepted non-positive or non-decimal quantity %q", value)
		}
	}
}
