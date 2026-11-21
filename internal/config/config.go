package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Mode string

const (
	ModePaper Mode = "paper"
	ModeDemo  Mode = "demo"
	ModeLive  Mode = "live"
)

type Config struct {
	Mode       Mode
	ListenAddr string
}

func FromEnv() (Config, error) {
	mode := Mode(strings.ToLower(strings.TrimSpace(os.Getenv("DELTA_MODE"))))
	if mode == "" {
		mode = ModePaper
	}
	switch mode {
	case ModePaper, ModeDemo:
	case ModeLive:
		return Config{}, errors.New("live trading blocked: release gates not implemented")
	default:
		return Config{}, fmt.Errorf("unsupported DELTA_MODE %q", mode)
	}
	addr := strings.TrimSpace(os.Getenv("DELTA_LISTEN_ADDR"))
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil || port == "" {
		return Config{}, fmt.Errorf("invalid DELTA_LISTEN_ADDR %q (expected host:port)", addr)
	}
	if host != "127.0.0.1" && host != "::1" && host != "localhost" {
		return Config{}, errors.New("foundation API must bind to loopback")
	}
	// A syntactically valid host:port is not necessarily a usable TCP port.
	// Reject service names, signed values, leading zeroes and out-of-range ports
	// before any listener is created.
	n, parseErr := strconv.Atoi(port)
	if parseErr != nil || n < 1 || n > 65535 || strconv.Itoa(n) != port {
		return Config{}, errors.New("foundation API requires a canonical TCP port in 1..65535")
	}
	return Config{Mode: mode, ListenAddr: addr}, nil
}
