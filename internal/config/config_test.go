package config

import (
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_DSN", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("ENV", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 8080 || cfg.LogLevel != LogInfo || len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "*" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	t.Setenv("PORT", "abc")
	t.Setenv("ALLOWED_ORIGINS", "*")
	t.Setenv("LOG_LEVEL", "info")
	if _, err := Load(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("LOG_LEVEL", "nope")
	t.Setenv("ALLOWED_ORIGINS", "*")
	t.Setenv("ENV", "dev")
	if _, err := Load(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoad_OriginsList(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("ALLOWED_ORIGINS", "http://a.com, http://b.com")
	t.Setenv("ENV", "prod")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Fatalf("unexpected origins: %+v", cfg.AllowedOrigins)
	}
}

func TestLoad_ProdWildcardRejected(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("ENV", "prod")
	t.Setenv("ALLOWED_ORIGINS", "*")
	if _, err := Load(); err == nil {
		t.Fatalf("expected error for wildcard in prod")
	}
}
