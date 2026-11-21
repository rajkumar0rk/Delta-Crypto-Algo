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
	for _, addr := range []string{":8080", "0.0.0.0:8080", "[::]:8080", "192.168.0.10:8080", "bad", "127.0.0.1:0"} {
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

func TestCanonicalListenerPortBounds(t *testing.T) {
	t.Setenv("DELTA_MODE", "paper")
	for _, port := range []string{"", "0", "00", "00080", "080", "-1", "+80", "65536", "99999", "http", "1.5"} {
		t.Run("reject_"+port, func(t *testing.T) {
			t.Setenv("DELTA_LISTEN_ADDR", "127.0.0.1:"+port)
			if _, err := FromEnv(); err == nil {
				t.Fatalf("noncanonical or invalid port %q accepted", port)
			}
		})
	}
	for _, port := range []string{"1", "80", "65535"} {
		t.Run("accept_"+port, func(t *testing.T) {
			t.Setenv("DELTA_LISTEN_ADDR", "127.0.0.1:"+port)
			if _, err := FromEnv(); err != nil {
				t.Fatalf("valid port %q rejected: %v", port, err)
			}
		})
	}
}
