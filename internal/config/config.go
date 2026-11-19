package config

import (
	"errors"
	"fmt"
	"net"
	"os"
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
	if port == "0" {
		return Config{}, errors.New("port 0 is disallowed")
	}
	return Config{Mode: mode, ListenAddr: addr}, nil
}
