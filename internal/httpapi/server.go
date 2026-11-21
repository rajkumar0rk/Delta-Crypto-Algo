package httpapi

import (
	"encoding/json"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type health struct {
	State        string `json:"state"`
	TradeEnabled bool   `json:"trade_enabled"`
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, health{State: "alive", TradeEnabled: false})
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusServiceUnavailable, health{State: "foundation_only", TradeEnabled: false})
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		var requestID [16]byte
		if _, err := rand.Read(requestID[:]); err != nil {
			http.Error(w, "request unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("X-Request-ID", hex.EncodeToString(requestID[:]))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		mux.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
