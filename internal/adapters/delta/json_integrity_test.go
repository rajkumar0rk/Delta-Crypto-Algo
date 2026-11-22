package delta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRejectAmbiguousJSON(t *testing.T) {
	cases := []struct {
		name string
		json string
		bad  bool
	}{
		{"valid nested", `{"success":true,"result":[{"id":1,"meta":{"strike":"0.01"}}],"meta":{"after":null}}`, false},
		{"case-ambiguous success", `{"Success":false,"success":true,"result":[{"id":1}],"meta":{"after":null}}`, true},
		{"escaped-equivalent key", `{"id":1,"\u0069d":2}`, true},
		{"case-ambiguous nested", `{"meta":{"After":1,"after":2}}`, true},
		{"invalid UTF-8", string([]byte{'{', '"', 0xff, '"', ':', '1', '}'}), true},
		{"object key limit", oversizedKeyObject(), true},
		{"duplicate success", `{"success":false,"success":true,"result":[{"id":1}],"meta":{"after":null}}`, true},
		{"duplicate cursor", `{"success":true,"result":[{"id":1}],"meta":{"after":"next","after":null}}`, true},
		{"duplicate product ID", `{"success":true,"result":[{"id":1,"id":2}],"meta":{"after":null}}`, true},
		{"duplicate nested spec", `{"success":true,"result":[{"id":1,"spec":{"lot":"1","lot":"2"}}],"meta":{"after":null}}`, true},
		{"same key in separate objects", `{"result":[{"id":1},{"id":2}]}`, false},
		{"malformed JSON", `{"result":[1,}`, true},
		{"extra document", `{} {}`, true},
		{"deep nesting", strings.Repeat("[", 66) + "0" + strings.Repeat("]", 66), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := rejectAmbiguousJSON([]byte(tc.json))
			if (err != nil) != tc.bad {
				t.Fatalf("invalidity mismatch: error=%v wantBad=%v", err, tc.bad)
			}
		})
	}
}

func TestPublicCatalogRejectsAmbiguousMetadataWithoutPartialPublication(t *testing.T) {
	for _, response := range []string{
		`{"success":true,"success":false,"result":[{"id":1}],"meta":{"after":null}}`,
		`{"success":true,"result":[{"id":1,"id":2}],"meta":{"after":null}}`,
		`{"success":true,"result":[{"id":1}],"meta":{"after":"next","after":null}}`,
	} {
		t.Run(response, func(t *testing.T) {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, response)
			}))
			defer srv.Close()
			client := srv.Client()
			client.Timeout = time.Second
			got, err := (PublicCatalog{BaseURL: srv.URL, Client: client}).Fetch(context.Background())
			if !errors.Is(err, ErrCatalogIncomplete) || len(got.Products) != 0 || got.Pages != 0 {
				t.Fatalf("ambiguous page must fail closed: %+v, %v", got, err)
			}
		})
	}
}

func oversizedKeyObject() string {
	keys := make([]string, 257)
	for i := range keys {
		keys[i] = fmt.Sprintf("%q:1", fmt.Sprintf("key%d", i))
	}
	return "{" + strings.Join(keys, ",") + "}"
}
