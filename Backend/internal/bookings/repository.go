package bookings

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func (r *Repository) Create(ctx context.Context, b *Booking) error {
	checkIn, err := time.Parse("2006-01-02", b.CheckIn)
	if err != nil {
		return &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "invalid check_in date, use YYYY-MM-DD"}
	}
	checkOut, err := time.Parse("2006-01-02", b.CheckOut)
	if err != nil {
		return &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "invalid check_out date, use YYYY-MM-DD"}
	}

	query := `
		INSERT INTO bookings (listing_id, guest_id, check_in, check_out, total_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, total_nights, status, created_at, updated_at
	`
	err = r.pool.QueryRow(ctx, query,
		b.ListingID, b.GuestID, checkIn, checkOut, b.TotalPrice,
	).Scan(&b.ID, &b.TotalNights, &b.Status, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
			return &respond.AppError{Code: http.StatusConflict, Message: "listing is already booked for these dates"}
		}
		return fmt.Errorf("Create: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Booking, error) {
	query := `
		SELECT id, listing_id, guest_id, check_in, check_out, total_nights, total_price, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`
	b := &Booking{}
	var checkIn, checkOut time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.ListingID, &b.GuestID,
		&checkIn, &checkOut, &b.TotalNights,
		&b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}

	b.CheckIn = checkIn.Format("2006-01-02")
	b.CheckOut = checkOut.Format("2006-01-02")
	return b, nil
}

func (r *Repository) ListByGuest(ctx context.Context, guestID string, limit, offset int) ([]Booking, int, error) {
	query := `
		SELECT id, listing_id, guest_id, check_in, check_out, total_nights, total_price, status, created_at, updated_at,
		       COUNT(*) OVER() AS total_count
		FROM bookings
		WHERE guest_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, guestID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ListByGuest: %w", err)
	}
	defer rows.Close()

	var results []Booking
	var total int

	for rows.Next() {
		b := Booking{}
		var checkIn, checkOut time.Time
		err := rows.Scan(
			&b.ID, &b.ListingID, &b.GuestID,
			&checkIn, &checkOut, &b.TotalNights,
			&b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("ListByGuest scan: %w", err)
		}
		b.CheckIn = checkIn.Format("2006-01-02")
		b.CheckOut = checkOut.Format("2006-01-02")
		results = append(results, b)
	}

	if results == nil {
		results = []Booking{}
	}
	return results, total, nil
}

func (r *Repository) Cancel(ctx context.Context, id string) error {
	query := `UPDATE bookings SET status = 'cancelled', updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
