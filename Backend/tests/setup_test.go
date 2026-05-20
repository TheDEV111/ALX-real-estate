package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/auth"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/bookings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/config"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/db"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/listings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/reviews"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/server"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	testServer *httptest.Server
	testPool   *pgxpool.Pool
)

func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err != nil {
		panic("config: " + err.Error())
	}

	if testDB := os.Getenv("TEST_DATABASE_URL"); testDB != "" {
		cfg.DatabaseURL = testDB
	}

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic("db connect: " + err.Error())
	}
	testPool = pool

	if err := db.RunMigrations(context.Background(), pool, "../migrations"); err != nil {
		panic("migrations: " + err.Error())
	}

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

	reviewsRepo := reviews.NewRepository(pool)
	reviewsSvc := reviews.NewService(reviewsRepo, bookingsRepo)
	reviewsHandler := reviews.NewHandler(reviewsSvc)

	router := server.NewRouter(cfg, authHandler, authMiddleware, listingsHandler, bookingsHandler, reviewsHandler)
	testServer = httptest.NewServer(router)

	code := m.Run()

	testServer.Close()
	pool.Close()
	os.Exit(code)
}

func truncate(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		"TRUNCATE reviews, bookings, listings, users CASCADE")
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func do(t *testing.T, method, path string, body any, token string) *http.Response {
	t.Helper()

	var b *bytes.Buffer
	if body != nil {
		data, _ := json.Marshal(body)
		b = bytes.NewBuffer(data)
	} else {
		b = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, testServer.URL+path, b)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func registerAndLogin(t *testing.T, email, password string) string {
	t.Helper()

	resp := do(t, http.MethodPost, "/api/auth/register", map[string]string{
		"email":     email,
		"password":  password,
		"full_name": "Test User",
	}, "")
	resp.Body.Close()

	resp = do(t, http.MethodPost, "/api/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, "")
	var result map[string]any
	decode(t, resp, &result)
	return result["access_token"].(string)
}
