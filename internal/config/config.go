package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type LogLevel string

const (
	LogDebug LogLevel = "debug"
	LogInfo  LogLevel = "info"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
)

type Config struct {
	Port           int
	DatabaseDSN    string
	LogLevel       LogLevel
	AllowedOrigins []string
	Env            string
}

func Load() (Config, error) {
	var cfg Config

	// Port
	portStr := getenv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, errors.New("invalid PORT")
	}
	cfg.Port = port

	// DB
	cfg.DatabaseDSN = os.Getenv("DB_DSN")

	// Log level
	level := strings.ToLower(getenv("LOG_LEVEL", string(LogInfo)))
	switch LogLevel(level) {
	case LogDebug, LogInfo, LogWarn, LogError:
		cfg.LogLevel = LogLevel(level)
	default:
		return Config{}, errors.New("invalid LOG_LEVEL")
	}

	// Allowed origins (comma-separated). "*" means any origin.
	origins := getenv("ALLOWED_ORIGINS", "*")
	if origins == "" {
		return Config{}, errors.New("ALLOWED_ORIGINS cannot be empty")
	}
	if origins == "*" {
		cfg.AllowedOrigins = []string{"*"}
	} else {
		parts := strings.Split(origins, ",")
		trimmed := make([]string, 0, len(parts))
		for _, p := range parts {
			s := strings.TrimSpace(p)
			if s != "" {
				trimmed = append(trimmed, s)
			}
		}
		if len(trimmed) == 0 {
			return Config{}, errors.New("ALLOWED_ORIGINS invalid")
		}
		cfg.AllowedOrigins = trimmed
	}

	// Environment: dev or prod (default dev)
	env := strings.ToLower(getenv("ENV", "dev"))
	switch env {
	case "dev", "prod":
		cfg.Env = env
	default:
		return Config{}, errors.New("invalid ENV (must be dev or prod)")
	}

	// In prod, wildcard origins are not allowed
	if cfg.Env == "prod" && len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*" {
		return Config{}, errors.New("ALLOWED_ORIGINS cannot be * in prod")
	}

	return cfg, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
