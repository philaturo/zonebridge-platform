// Package main is the entry point for the ZoneBridge backend application.
// It is intentionally minimal, delegating all orchestration to the application package.
package main

import (
	"context"
	"log"

	"github.com/philaturo/zonebridge-platform/internal/application"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := app.Run(context.Background()); err != nil {
		log.Fatalf("application terminated with error: %v", err)
	}
}