package server

import (
	"net/http"
	"time"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/auth"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/bookings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/config"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/docs"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/listings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func NewRouter(
	cfg *config.Config,
	authHandler *auth.Handler,
	authMiddleware func(http.Handler) http.Handler,
	listingsHandler *listings.Handler,
	bookingsHandler *bookings.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(30 * time.Second))
	if !cfg.DisableRateLimit {
		r.Use(httprate.LimitByIP(100, time.Minute))
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		r.Get("/docs", docs.UIHandler)
		r.Get("/docs/spec", docs.SpecHandler)

		r.Group(func(r chi.Router) {
			if !cfg.DisableRateLimit {
				r.Use(httprate.LimitByIP(5, time.Minute))
			}
			r.Post("/auth/register", authHandler.Register)
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.Refresh)
		})

		r.Get("/listings", listingsHandler.List)
		r.Get("/listings/{id}", listingsHandler.GetByID)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/users/me", authHandler.Me)

			r.Post("/listings", listingsHandler.Create)
			r.Patch("/listings/{id}", listingsHandler.Update)
			r.Delete("/listings/{id}", listingsHandler.Delete)

			r.Post("/bookings", bookingsHandler.Create)
			r.Get("/bookings", bookingsHandler.ListMine)
			r.Get("/bookings/{id}", bookingsHandler.GetByID)
			r.Patch("/bookings/{id}/cancel", bookingsHandler.Cancel)
		})
	})

	return r
}
