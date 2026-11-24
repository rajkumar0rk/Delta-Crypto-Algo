package delta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rajkumar0rk/Delta-Crypto-Algo/internal/domain/catalog"
)

// PublicTerms preserves the venue's published values without interpreting
// margin units, fee tiers, account discounts or option-writing obligations.
type PublicTerms struct {
	ContractUnitCurrency string
	NotionalType         string
	InitialMargin        string
	MaintenanceMargin    string
	MakerCommission      string
	TakerCommission      string
}

// LotEvidence is separate from the Products response, which does not document
// a per-product lot-step field. It must be supplied by a verified source.
type LotEvidence struct {
	ProductID     int64
	Symbol        string
	ContractValue string
	LotSize       string
	Source        string
	VerifiedAt    time.Time
}

type mappedProduct struct {
	product catalog.Product
	terms   PublicTerms
}

// decodeProduct maps documented Delta India Product fields, not symbol
// heuristics. Unknown notional models are rejected rather than guessed.
func decodeProduct(raw json.RawMessage, lot LotEvidence, observed time.Time, maxAge time.Duration) (mappedProduct, error) {
	var out mappedProduct
	if len(raw) > maxPageBytes {
		return out, errors.New("oversized product")
	}
	if err := rejectAmbiguousJSON(raw); err != nil {
		return out, err
	}
	var v struct {
		ID                   int64      `json:"id"`
		Symbol               string     `json:"symbol"`
		ContractType         string     `json:"contract_type"`
		State                string     `json:"state"`
		TradingStatus        string     `json:"trading_status"`
		TickSize             string     `json:"tick_size"`
		ContractValue        string     `json:"contract_value"`
		ContractUnitCurrency string     `json:"contract_unit_currency"`
		NotionalType         string     `json:"notional_type"`
		IsQuanto             *bool      `json:"is_quanto"`
		SettlementTime       *time.Time `json:"settlement_time"`
		StrikePrice          *string    `json:"strike_price"`
		InitialMargin        string     `json:"initial_margin"`
		MaintenanceMargin    string     `json:"maintenance_margin"`
		MakerCommission      string     `json:"maker_commission_rate"`
		TakerCommission      string     `json:"taker_commission_rate"`
		UnderlyingAsset      struct {
			Symbol string `json:"symbol"`
		} `json:"underlying_asset"`
		QuotingAsset struct {
			Symbol string `json:"symbol"`
		} `json:"quoting_asset"`
		SettlingAsset struct {
			Symbol string `json:"symbol"`
		} `json:"settling_asset"`
		ProductSpecs struct {
			ReduceOnly *bool `json:"only_reduce_only_orders_allowed"`
		} `json:"product_specs"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return out, fmt.Errorf("invalid product: %w", err)
	}
	if maxAge <= 0 || lot.ProductID != v.ID || lot.Symbol != v.Symbol || lot.ContractValue != v.ContractValue || strings.TrimSpace(lot.Source) == "" || lot.VerifiedAt.IsZero() || lot.VerifiedAt.After(observed) || observed.Sub(lot.VerifiedAt) > maxAge {
		return out, errors.New("missing or stale matching lot evidence")
	}
	if v.NotionalType != "vanilla" || v.IsQuanto == nil || *v.IsQuanto || v.ContractUnitCurrency != v.UnderlyingAsset.Symbol {
		return out, errors.New("unverified contract value model")
	}
	p := catalog.Product{
		ID:                 v.ID,
		Symbol:             v.Symbol,
		Underlying:         v.UnderlyingAsset.Symbol,
		TickSize:           v.TickSize,
		LotSize:            lot.LotSize,
		ContractMultiplier: v.ContractValue,
		QuoteCurrency:      v.QuotingAsset.Symbol,
		SettlementCurrency: v.SettlingAsset.Symbol,
		Expiry:             v.SettlementTime,
	}
	switch v.ContractType {
	case "perpetual_futures":
		p.Kind = catalog.Perpetual
	case "futures":
		p.Kind = catalog.DatedFuture
	case "call_options":
		p.Kind = catalog.CallOption
	case "put_options":
		p.Kind = catalog.PutOption
	default:
		return out, errors.New("unsupported contract type")
	}
	if v.StrikePrice != nil {
		p.Strike = *v.StrikePrice
	}
	switch v.State {
	case "live":
		p.Status = catalog.Inactive
		if v.TradingStatus == "operational" && v.ProductSpecs.ReduceOnly != nil && !*v.ProductSpecs.ReduceOnly {
			p.Status = catalog.Active
		}
	case "expired":
		p.Status = catalog.Expired
	default:
		p.Status = catalog.Unknown
	}
	if p.Expiry != nil && !p.Expiry.After(observed) {
		p.Status = catalog.Expired
	}
	if err := p.Validate(); err != nil {
		return out, err
	}
	out.product = p
	out.terms = PublicTerms{
		ContractUnitCurrency: v.ContractUnitCurrency,
		NotionalType:         v.NotionalType,
		InitialMargin:        v.InitialMargin,
		MaintenanceMargin:    v.MaintenanceMargin,
		MakerCommission:      v.MakerCommission,
		TakerCommission:      v.TakerCommission,
	}
	return out, nil
}

// MapSnapshot rejects the entire snapshot on any unverified mapping or
// duplicate identity. PublicTerms are informational, not Capability evidence.
func MapSnapshot(raw RawSnapshot, lots map[int64]LotEvidence, maxAge time.Duration) (catalog.Snapshot, map[int64]PublicTerms, error) {
	var empty catalog.Snapshot
	if raw.Pages <= 0 || raw.Pages > maxPages || raw.AsOf.IsZero() || len(raw.Products) == 0 || len(raw.Products) > maxProducts {
		return empty, nil, ErrCatalogIncomplete
	}
	next := catalog.Snapshot{Source: "delta-india-public", AsOf: raw.AsOf, Complete: true}
	terms := make(map[int64]PublicTerms, len(raw.Products))
	for i, record := range raw.Products {
		var identity struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(record, &identity); err != nil {
			return empty, nil, fmt.Errorf("%w: product %d identity", ErrCatalogIncomplete, i)
		}
		mapped, err := decodeProduct(record, lots[identity.ID], raw.AsOf, maxAge)
		if err != nil {
			return empty, nil, fmt.Errorf("%w: product %d: %v", ErrCatalogIncomplete, i, err)
		}
		next.Products = append(next.Products, mapped.product)
		terms[mapped.product.ID] = mapped.terms
	}
	if err := next.Validate(); err != nil {
		return empty, nil, fmt.Errorf("%w: %v", ErrCatalogIncomplete, err)
	}
	return next, terms, nil
}

// RefreshCatalog is read-only. Failed fetch or mapping leaves the last
// validated store snapshot unchanged. Returned terms belong to that refresh,
// not to a later store read; they are never permission to place an order.
func (c PublicCatalog) RefreshCatalog(ctx context.Context, store *catalog.Store, lots map[int64]LotEvidence, maxAge time.Duration) ([]catalog.ClassifiedChange, map[int64]PublicTerms, error) {
	if store == nil || ctx == nil {
		return nil, nil, ErrCatalogIncomplete
	}
	raw, err := c.Fetch(ctx)
	if err != nil {
		return nil, nil, err
	}
	next, terms, err := MapSnapshot(raw, lots, maxAge)
	if err != nil {
		return nil, nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	changes, err := store.PublishWithChanges(next)
	if err != nil {
		return nil, nil, err
	}
	return changes, terms, nil
}
