package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"

	"github.com/mashurimansur/goCMS/internal/utils/database"
)

// AppConfig aggregates all runtime configuration required by the application.
type AppConfig struct {
	HTTPAddr string
	GinMode  string
	Database database.Config
}

// Load reads the provided .env files (if present) and maps environment variables to AppConfig.
// Missing .env files are ignored so the service can still rely on real environment variables.
func Load(envFiles ...string) (AppConfig, error) {
	if len(envFiles) == 0 {
		envFiles = []string{".env"}
	}

	if err := loadEnvFiles(envFiles); err != nil {
		return AppConfig{}, err
	}

	cfg := AppConfig{
		HTTPAddr: envOrDefault("HTTP_ADDR", ":8080"),
		GinMode:  os.Getenv("GIN_MODE"),
		Database: database.Config{
			Driver: os.Getenv("DB_DRIVER"),
			DSN:    os.Getenv("DB_DSN"),
		},
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadEnvFiles(envFiles []string) error {
	filesToLoad := make([]string, 0, len(envFiles))

	for _, path := range envFiles {
		info, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}

		if info.IsDir() {
			continue
		}

		filesToLoad = append(filesToLoad, path)
	}

	if len(filesToLoad) == 0 {
		return nil
	}

	return godotenv.Load(filesToLoad...)
}
