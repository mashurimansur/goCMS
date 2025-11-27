package main

import (
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mashurimansur/goCMS/internal/app"
	"github.com/mashurimansur/goCMS/internal/utils/config"
)

func main() {
	ctx := context.Background()

	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	application, err := app.New(ctx, appConfig)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer application.Close()

	log.Printf("HTTP server listening on %s", appConfig.HTTPAddr)
	if err := application.Run(); err != nil {
		log.Fatalf("failed to run HTTP server: %v", err)
	}
}
