package main

import (
	"context"
	"log"
	"net/http"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/config"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/db"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()
	log.Println("database connection established")

	if err := db.RunMigrations(ctx, pool, "migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations applied")

	router := server.NewRouter(cfg)

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}
