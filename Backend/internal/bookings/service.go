package bookings

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/listings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
)

type Service struct {
	repo         *Repository
	listingsRepo *listings.Repository
}

func NewService(repo *Repository, listingsRepo *listings.Repository) *Service {
	return &Service{repo: repo, listingsRepo: listingsRepo}
}

func (s *Service) Create(ctx context.Context, guestID string, req CreateBookingRequest) (*Booking, error) {
	if req.ListingID == "" || req.CheckIn == "" || req.CheckOut == "" {
		return nil, &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "listing_id, check_in and check_out are required"}
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		return nil, &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "invalid check_in date, use YYYY-MM-DD"}
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		return nil, &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "invalid check_out date, use YYYY-MM-DD"}
	}
	if !checkOut.After(checkIn) {
		return nil, &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "check_out must be after check_in"}
	}
	if checkIn.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "check_in cannot be in the past"}
	}

	listing, err := s.listingsRepo.GetByID(ctx, req.ListingID)
	if err != nil {
		return nil, err
	}
	if listing == nil || !listing.IsActive {
		return nil, respond.ErrNotFound
	}

	nights := checkOut.Sub(checkIn).Hours() / 24
	totalPrice := math.Round(listing.PricePerNight*nights*100) / 100

	b := &Booking{
		ListingID:  req.ListingID,
		GuestID:    guestID,
		CheckIn:    req.CheckIn,
		CheckOut:   req.CheckOut,
		TotalPrice: totalPrice,
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetByID(ctx context.Context, id, requesterID, requesterRole string) (*Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, respond.ErrNotFound
	}
	if b.GuestID != requesterID && requesterRole != "admin" {
		return nil, respond.ErrForbidden
	}
	return b, nil
}

func (s *Service) ListMine(ctx context.Context, guestID string, limit, offset int) (*BookingsResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	bookings, total, err := s.repo.ListByGuest(ctx, guestID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}

	return &BookingsResponse{
		Data:   bookings,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) Cancel(ctx context.Context, id, requesterID, requesterRole string) (*Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, respond.ErrNotFound
	}
	if b.GuestID != requesterID && requesterRole != "admin" {
		return nil, respond.ErrForbidden
	}
	if b.Status == "cancelled" || b.Status == "completed" {
		return nil, &respond.AppError{Code: http.StatusConflict, Message: fmt.Sprintf("booking is already %s", b.Status)}
	}

	if err := s.repo.Cancel(ctx, id); err != nil {
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	b.Status = "cancelled"
	return b, nil
}
