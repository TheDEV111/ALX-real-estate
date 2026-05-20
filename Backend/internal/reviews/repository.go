package reviews

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, rev *Review) error {
	query := `
		INSERT INTO reviews (booking_id, listing_id, guest_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query,
		rev.BookingID, rev.ListingID, rev.GuestID, rev.Rating, rev.Comment,
	).Scan(&rev.ID, &rev.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &respond.AppError{Code: http.StatusConflict, Message: "you have already reviewed this booking"}
		}
		return fmt.Errorf("Create review: %w", err)
	}
	return nil
}

func (r *Repository) GetByBookingID(ctx context.Context, bookingID string) (*Review, error) {
	query := `SELECT id, booking_id, listing_id, guest_id, rating, comment, created_at FROM reviews WHERE booking_id = $1`
	rev := &Review{}
	err := r.pool.QueryRow(ctx, query, bookingID).Scan(
		&rev.ID, &rev.BookingID, &rev.ListingID, &rev.GuestID,
		&rev.Rating, &rev.Comment, &rev.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByBookingID: %w", err)
	}
	return rev, nil
}

func (r *Repository) ListByListing(ctx context.Context, listingID string) ([]Review, float64, error) {
	query := `
		SELECT id, booking_id, listing_id, guest_id, rating, comment, created_at,
		       AVG(rating) OVER() AS avg_rating
		FROM reviews
		WHERE listing_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, 0, fmt.Errorf("ListByListing: %w", err)
	}
	defer rows.Close()

	var results []Review
	var avg float64

	for rows.Next() {
		rev := Review{}
		if err := rows.Scan(
			&rev.ID, &rev.BookingID, &rev.ListingID, &rev.GuestID,
			&rev.Rating, &rev.Comment, &rev.CreatedAt,
			&avg,
		); err != nil {
			return nil, 0, fmt.Errorf("ListByListing scan: %w", err)
		}
		results = append(results, rev)
	}

	if results == nil {
		results = []Review{}
	}
	return results, avg, nil
}
