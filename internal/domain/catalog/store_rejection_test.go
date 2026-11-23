package catalog

import (
 "errors"
 "testing"
 "time"
)

// A failed venue refresh must never replace a previously validated snapshot.
// Even after a failed refresh, stale data must not become tradeable again.
func TestStoreRejectedRefreshPreservesPublishedCatalog(t *testing.T) {
 base := time.Date(2026, 11, 23, 12, 0, 0, 0, time.UTC)
 cases := []struct {
  name string
  mutate func(*Snapshot)
 }{
  {"incomplete pagination", func(s *Snapshot) { s.Complete = false }},
  {"missing source", func(s *Snapshot) { s.Source = "" }},
  {"duplicate id", func(s *Snapshot) { duplicate := s.Products[0]; duplicate.Symbol = "OTHER-PERP"; s.Products = append(s.Products, duplicate) }},
  {"duplicate symbol", func(s *Snapshot) { duplicate := s.Products[0]; duplicate.ID = 2; s.Products = append(s.Products, duplicate) }},
  {"unknown product kind", func(s *Snapshot) { s.Products[0].Kind = Kind("unrecognized") }},
  {"unknown product status", func(s *Snapshot) { s.Products[0].Status = Status("unrecognized") }},
  {"nondecimal tick", func(s *Snapshot) { s.Products[0].TickSize = "1/10" }},
  {"zero multiplier", func(s *Snapshot) { s.Products[0].ContractMultiplier = "0" }},
  {"empty catalog", func(s *Snapshot) { s.Products = nil }},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   var store Store
   original := testSnapshot(base)
   if err := store.Publish(original); err != nil { t.Fatal(err) }
   rejected := testSnapshot(base.Add(time.Second))
   tc.mutate(&rejected)
   if err := store.Publish(rejected); !errors.Is(err, ErrCatalogUnavailable) {
    t.Fatalf("expected unavailable refresh rejection, got %v", err)
   }
   got, err := store.Snapshot()
   if err != nil { t.Fatal(err) }
   if got.AsOf != base || len(got.Products) != 1 || got.Products[0].Symbol != original.Products[0].Symbol {
    t.Fatalf("failed refresh mutated published snapshot: %+v", got)
   }
   if err := store.Tradeable(1, base.Add(30*time.Second), time.Minute); err != nil {
    t.Fatalf("last good snapshot should remain usable while fresh: %v", err)
   }
   if err := store.Tradeable(1, base.Add(2*time.Minute), time.Minute); !errors.Is(err, ErrNotTradeable) {
    t.Fatalf("last good snapshot must fail closed when stale: %v", err)
   }
  })
 }
}
