package main

import (
	"log"

	"factorial-cal-services/migrations"
	"factorial-cal-services/pkg/config"
)

func main() {
	log.Println("Starting database migrations...")

	cfg := config.LoadConfig()

	if err := migrations.RunMigrations(cfg.DSN()); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✓ All migrations completed successfully")
}
