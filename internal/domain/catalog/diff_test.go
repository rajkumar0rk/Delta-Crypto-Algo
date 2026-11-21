package catalog

import (
 "reflect"
 "testing"
 "time"
)

func TestDiffAddedRemovedModifiedAndStableOrder(t *testing.T) {
 old := fixture()
 next := fixture()
 next.AsOf = next.AsOf.Add(time.Minute)
 next.Products = []Product{old.Products[1]}
 next.Products[0].Status = Delisted
 added := old.Products[0]
 added.ID = 2
 added.Symbol = "NEW-RWA-PERP"
 next.Products = append(next.Products, added)
 changes, err := Diff(old, next)
 if err != nil { t.Fatal(err) }
 if len(changes) != 3 { t.Fatalf("want three changes, got %+v", changes) }
 if changes[0].ID != 2 || changes[0].Before != nil || changes[0].After == nil { t.Fatalf("added product: %+v", changes[0]) }
 if changes[1].ID != 10 || changes[1].Before == nil || changes[1].After != nil { t.Fatalf("removed product: %+v", changes[1]) }
 if changes[2].ID != 11 || changes[2].Before.Status != Active || changes[2].After.Status != Delisted { t.Fatalf("status transition: %+v", changes[2]) }
 if !reflect.DeepEqual(old, fixture()) { t.Fatal("source snapshot mutated") }
 // Audit copies must not alias expiry pointers in either snapshot.
 *changes[2].After.Expiry = time.Time{}
 if next.Products[0].Expiry.IsZero() { t.Fatal("audit result aliases source expiry") }
}

func TestDiffRejectsInvalidOrNonMonotonicSnapshots(t *testing.T) {
 old := fixture()
 for _, tc := range []struct{name string; mutate func(*Snapshot)}{
  {"same timestamp", func(s *Snapshot) {}},
  {"older timestamp", func(s *Snapshot) { s.AsOf = s.AsOf.Add(-time.Second) }},
  {"different source", func(s *Snapshot) { s.AsOf = s.AsOf.Add(time.Second); s.Source = "other" }},
  {"partial pagination", func(s *Snapshot) { s.AsOf = s.AsOf.Add(time.Second); s.Complete = false }},
  {"duplicate product", func(s *Snapshot) { s.AsOf = s.AsOf.Add(time.Second); s.Products = append(s.Products, s.Products[0]) }},
 } {
  t.Run(tc.name, func(t *testing.T) {
   next := fixture(); tc.mutate(&next)
   if changes, err := Diff(old,next); err == nil || changes != nil { t.Fatalf("expected rejection: changes=%v err=%v",changes,err) }
  })
 }
 next := fixture(); next.AsOf = next.AsOf.Add(time.Second)
 changes, err := Diff(old,next)
 if err != nil || len(changes) != 0 { t.Fatalf("unchanged products: %+v, %v",changes,err) }
}
