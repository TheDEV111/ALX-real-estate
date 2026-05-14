package main

import (
	"context"
	"log"
	"net/http"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/auth"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/bookings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/config"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/db"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/listings"
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

	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authSvc)
	authMiddleware := auth.Authenticator(authSvc)

	listingsRepo := listings.NewRepository(pool)
	listingsSvc := listings.NewService(listingsRepo)
	listingsHandler := listings.NewHandler(listingsSvc)

	bookingsRepo := bookings.NewRepository(pool)
	bookingsSvc := bookings.NewService(bookingsRepo, listingsRepo)
	bookingsHandler := bookings.NewHandler(bookingsSvc)

	router := server.NewRouter(cfg, authHandler, authMiddleware, listingsHandler, bookingsHandler)

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}
