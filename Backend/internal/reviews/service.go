package reviews

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/bookings"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
)

type Service struct {
	repo         *Repository
	bookingsRepo *bookings.Repository
}

func NewService(repo *Repository, bookingsRepo *bookings.Repository) *Service {
	return &Service{repo: repo, bookingsRepo: bookingsRepo}
}

func (s *Service) Create(ctx context.Context, bookingID, guestID string, req CreateReviewRequest) (*Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "rating must be between 1 and 5"}
	}

	b, err := s.bookingsRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, respond.ErrNotFound
	}
	if b.GuestID != guestID {
		return nil, respond.ErrForbidden
	}
	if b.Status != "completed" {
		return nil, &respond.AppError{Code: http.StatusConflict, Message: fmt.Sprintf("cannot review a booking with status %q", b.Status)}
	}

	rev := &Review{
		BookingID: bookingID,
		ListingID: b.ListingID,
		GuestID:   guestID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	if err := s.repo.Create(ctx, rev); err != nil {
		return nil, err
	}
	return rev, nil
}

func (s *Service) ListForListing(ctx context.Context, listingID string) (*ReviewsResponse, error) {
	reviews, avg, err := s.repo.ListByListing(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}
	return &ReviewsResponse{
		Data:          reviews,
		Total:         len(reviews),
		AverageRating: avg,
	}, nil
}
