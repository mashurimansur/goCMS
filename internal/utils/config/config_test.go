package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_FromEnvFile(t *testing.T) {
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "app.env")

	content := []byte("HTTP_ADDR=127.0.0.1:9090\nGIN_MODE=release\nDB_DRIVER=postgres\nDB_DSN=postgres://user:pass@localhost/db\n")
	if err := os.WriteFile(envFile, content, 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	for _, key := range []string{"HTTP_ADDR", "GIN_MODE", "DB_DRIVER", "DB_DSN"} {
		t.Setenv(key, "")
		// Unset so godotenv can populate values from the file.
		os.Unsetenv(key)
	}

	cfg, err := Load(envFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.HTTPAddr != "127.0.0.1:9090" {
		t.Fatalf("unexpected HTTPAddr: %s", cfg.HTTPAddr)
	}
	if cfg.GinMode != "release" {
		t.Fatalf("unexpected GinMode: %s", cfg.GinMode)
	}
	if cfg.Database.Driver != "postgres" {
		t.Fatalf("unexpected database config: %+v", cfg.Database)
	}
}

func TestLoad_UsesDefaultsWhenMissing(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("GIN_MODE", "")
	t.Setenv("DB_DRIVER", "")
	t.Setenv("DB_DSN", "")

	cfg, err := Load(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default HTTPAddr, got %s", cfg.HTTPAddr)
	}
	if cfg.GinMode != "" {
		t.Fatalf("expected empty GinMode, got %s", cfg.GinMode)
	}
	if cfg.Database.Driver != "" {
		t.Fatalf("expected empty database config, got %+v", cfg.Database)
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("SAMPLE_KEY", "value")
	if got := envOrDefault("SAMPLE_KEY", "fallback"); got != "value" {
		t.Fatalf("expected env value, got %s", got)
	}

	t.Setenv("SAMPLE_KEY", "")
	if got := envOrDefault("SAMPLE_KEY", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %s", got)
	}
}

func TestLoadEnvFiles(t *testing.T) {
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "vars.env")
	if err := os.WriteFile(envFile, []byte("FROM_ENV=value"), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	t.Setenv("FROM_ENV", "")
	os.Unsetenv("FROM_ENV")
	if err := loadEnvFiles([]string{envFile}); err != nil {
		t.Fatalf("loadEnvFiles returned error: %v", err)
	}

	if got := os.Getenv("FROM_ENV"); got != "value" {
		t.Fatalf("expected env var to be loaded, got %q", got)
	}

	// Missing and directory entries should be ignored without error.
	dir := filepath.Join(tempDir, "dir")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := loadEnvFiles([]string{filepath.Join(tempDir, "missing.env"), dir}); err != nil {
		t.Fatalf("expected missing files and directories to be skipped: %v", err)
	}
}
