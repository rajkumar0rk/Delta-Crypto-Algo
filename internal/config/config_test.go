package config

import "testing"

func TestPaperDefaults(t *testing.T) {
	t.Setenv("DELTA_MODE", "")
	t.Setenv("DELTA_LISTEN_ADDR", "")
	cfg, err := FromEnv()
	if err != nil || cfg.Mode != ModePaper || cfg.ListenAddr != "127.0.0.1:8080" {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
}

func TestRefuseLiveAndPublicListen(t *testing.T) {
	t.Setenv("DELTA_LISTEN_ADDR", "")
	t.Setenv("DELTA_MODE", "live")
	if _, err := FromEnv(); err == nil {
		t.Fatal("live mode must fail")
	}
	t.Setenv("DELTA_MODE", "paper")
	for _, addr := range []string{":8080", "0.0.0.0:8080", "[::]:8080", "192.168.0.10:8080", "bad", "127.0.0.1:0", "127.0.0.1:abc", "127.0.0.1:-1", "127.0.0.1:65536", "localhost:999999999999999999999999"} {
		t.Setenv("DELTA_LISTEN_ADDR", addr)
		if _, err := FromEnv(); err == nil {
			t.Fatalf("listener %q must fail", addr)
		}
	}
}

func TestDemoLoopbackAllowed(t *testing.T) {
	t.Setenv("DELTA_MODE", "DEMO")
	t.Setenv("DELTA_LISTEN_ADDR", "localhost:9090")
	cfg, err := FromEnv()
	if err != nil || cfg.Mode != ModeDemo {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
}

func TestValidLoopbackPortBoundaries(t *testing.T) {
	t.Setenv("DELTA_MODE", "paper")
	for _, addr := range []string{"127.0.0.1:1", "localhost:65535", "[::1]:443"} {
		t.Setenv("DELTA_LISTEN_ADDR", addr)
		cfg, err := FromEnv()
		if err != nil || cfg.ListenAddr != addr {
			t.Fatalf("listener %q: cfg=%+v err=%v", addr, cfg, err)
		}
	}
}
