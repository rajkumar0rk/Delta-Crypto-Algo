package delta

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// rejectAmbiguousJSON refuses duplicate object keys at any nesting level.
// Go's default JSON decoder silently accepts duplicate keys and uses the last
// value. Exchange metadata must never depend on that ambiguous interpretation.
// The input is already bounded by the HTTP page byte limit.
func rejectAmbiguousJSON(data []byte) error {
	if !utf8.Valid(data) {
		return errors.New("invalid UTF-8 JSON")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := inspectJSONValue(dec, 0); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("trailing JSON data")
	}
	return nil
}

func inspectJSONValue(dec *json.Decoder, depth int) error {
	if depth > 64 {
		return errors.New("JSON nesting exceeds limit")
	}
	token, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := make(map[string]struct{})
		seenFolded := make([]string, 0)
		for dec.More() {
			if len(seen) >= 256 {
				return errors.New("JSON object exceeds key limit")
			}
			keyToken, err := dec.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errors.New("non-string object key")
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("duplicate JSON key %q", key)
			}
			for _, previous := range seenFolded {
				if strings.EqualFold(previous, key) {
					return fmt.Errorf("case-ambiguous JSON keys %q and %q", previous, key)
				}
			}
			seen[key] = struct{}{}
			seenFolded = append(seenFolded, key)
			if err := inspectJSONValue(dec, depth+1); err != nil {
				return err
			}
		}
		_, err := dec.Token()
		return err
	case '[':
		for dec.More() {
			if err := inspectJSONValue(dec, depth+1); err != nil {
				return err
			}
		}
		_, err := dec.Token()
		return err
	default:
		return errors.New("unexpected JSON delimiter")
	}
}
