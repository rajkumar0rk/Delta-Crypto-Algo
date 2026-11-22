// Package delta provides public, read-only Delta Exchange India adapters.
// No private credentials or order endpoints are used in this package.
package delta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	maxPageBytes = 4 << 20
	maxPages     = 100
	maxProducts  = 100000
)

var ErrCatalogIncomplete = errors.New("delta catalog incomplete")

// PublicCatalog retrieves all pages from /v2/products. It deliberately returns
// raw product records: a separate, verified contract mapper must classify and
// validate venue fields before anything can enter the tradeability catalog.
type PublicCatalog struct {
	BaseURL     string
	Client      *http.Client
	MaxPages    int
	MaxProducts int
}

type RawSnapshot struct {
	AsOf     time.Time
	Products []json.RawMessage
	Pages    int
}

type pageResponse struct {
	Success *bool             `json:"success"`
	Result  []json.RawMessage `json:"result"`
	Meta    struct {
		After json.RawMessage `json:"after"`
	} `json:"meta"`
}

func (c PublicCatalog) Fetch(ctx context.Context) (RawSnapshot, error) {
	var empty RawSnapshot
	if ctx == nil {
		return empty, fmt.Errorf("%w: nil context", ErrCatalogIncomplete)
	}
	base, err := url.Parse(c.BaseURL)
	if err != nil || base == nil || base.Scheme != "https" || base.Host == "" || base.User != nil || base.RawQuery != "" || base.Fragment != "" || (base.Path != "" && base.Path != "/") {
		return empty, fmt.Errorf("%w: expected HTTPS origin without path or query", ErrCatalogIncomplete)
	}
	if c.Client == nil || c.Client.Timeout <= 0 {
		return empty, fmt.Errorf("%w: HTTP client requires explicit timeout", ErrCatalogIncomplete)
	}
	pageLimit := c.MaxPages
	if pageLimit == 0 {
		pageLimit = maxPages
	}
	productLimit := c.MaxProducts
	if productLimit == 0 {
		productLimit = maxProducts
	}
	if pageLimit < 1 || pageLimit > maxPages || productLimit < 1 || productLimit > maxProducts {
		return empty, fmt.Errorf("%w: invalid bounded limits", ErrCatalogIncomplete)
	}
	out := RawSnapshot{AsOf: time.Now().UTC()}
	seen := map[string]struct{}{}
	cursor := ""
	for page := 0; page < pageLimit; page++ {
		endpoint := *base
		endpoint.Path = "/v2/products"
		q := url.Values{}
		if cursor != "" {
			q.Set("after", cursor)
		}
		endpoint.RawQuery = q.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
		if err != nil {
			return empty, fmt.Errorf("%w: request: %v", ErrCatalogIncomplete, err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "identity")
		// Never follow redirects: a configurable HTTPS origin must not be able
		// to redirect catalog reads to another host (or an untrusted path).
		// Copy the client to preserve caller configuration without mutation.
		client := *c.Client
		client.CheckRedirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
		resp, err := client.Do(req)
		if err != nil {
			return empty, fmt.Errorf("%w: transport: %v", ErrCatalogIncomplete, err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return empty, fmt.Errorf("%w: HTTP %d", ErrCatalogIncomplete, resp.StatusCode)
		}
		// Do not trust content negotiation alone. A proxy must not supply a
		// compressed or non-JSON body that could be interpreted inconsistently.
		// Missing Content-Type is tolerated for legacy venue responses, but a
		// declared type must be JSON. The byte cap still applies in every case.
		if encoding := strings.TrimSpace(resp.Header.Get("Content-Encoding")); encoding != "" && !strings.EqualFold(encoding, "identity") {
			resp.Body.Close()
			return empty, fmt.Errorf("%w: unsupported content encoding", ErrCatalogIncomplete)
		}
		if declared := resp.Header.Get("Content-Type"); declared != "" {
			mediaType := strings.ToLower(strings.TrimSpace(strings.SplitN(declared, ";", 2)[0]))
			if mediaType != "application/json" {
				resp.Body.Close()
				return empty, fmt.Errorf("%w: unexpected content type", ErrCatalogIncomplete)
			}
		}
		// Enforce the byte limit independently of JSON parsing. A valid JSON
		// document followed by excess whitespace must not bypass the cap.
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxPageBytes+1))
		closeErr := resp.Body.Close()
		if readErr != nil || closeErr != nil || len(body) > maxPageBytes {
			return empty, fmt.Errorf("%w: oversized or unreadable page: %v %v", ErrCatalogIncomplete, readErr, closeErr)
		}
		var data pageResponse
		decoder := json.NewDecoder(bytes.NewReader(body))
		err = decoder.Decode(&data)
		if err == nil {
			var extra json.RawMessage
			if trailing := decoder.Decode(&extra); trailing != io.EOF {
				err = errors.New("trailing response")
			}
		}
		if err != nil {
			return empty, fmt.Errorf("%w: invalid JSON: %v", ErrCatalogIncomplete, err)
		}
		if data.Success == nil || !*data.Success || len(data.Result) == 0 {
			return empty, fmt.Errorf("%w: failed or empty page", ErrCatalogIncomplete)
		}
		if len(data.Result) > productLimit-len(out.Products) {
			return empty, fmt.Errorf("%w: product limit exceeded", ErrCatalogIncomplete)
		}
		for _, product := range data.Result {
			var obj map[string]json.RawMessage
			if len(product) == 0 || json.Unmarshal(product, &obj) != nil || len(obj) == 0 {
				return empty, fmt.Errorf("%w: malformed product", ErrCatalogIncomplete)
			}
			out.Products = append(out.Products, append(json.RawMessage(nil), product...))
		}
		out.Pages++
		raw := strings.TrimSpace(string(data.Meta.After))
		if raw == "" || raw == "null" {
			return out, nil
		}
		var next string
		if err := json.Unmarshal(data.Meta.After, &next); err != nil {
			var num json.Number
			if err2 := json.Unmarshal(data.Meta.After, &num); err2 != nil {
				return empty, fmt.Errorf("%w: unsupported cursor", ErrCatalogIncomplete)
			}
			next = num.String()
			if _, err := strconv.ParseInt(next, 10, 64); err != nil {
				return empty, fmt.Errorf("%w: invalid numeric cursor", ErrCatalogIncomplete)
			}
		}
		if next == "" || len(next) > 256 {
			return empty, fmt.Errorf("%w: empty or oversized cursor", ErrCatalogIncomplete)
		}
		if _, exists := seen[next]; exists {
			return empty, fmt.Errorf("%w: repeated cursor", ErrCatalogIncomplete)
		}
		seen[next] = struct{}{}
		cursor = next
	}
	return empty, fmt.Errorf("%w: pagination page limit reached", ErrCatalogIncomplete)
}
