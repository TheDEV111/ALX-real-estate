package listings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, hostID string, req CreateListingRequest) (*Listing, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	if req.Amenities == nil {
		req.Amenities = []string{}
	}
	if req.Images == nil {
		req.Images = []string{}
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}

	l := &Listing{
		HostID:        hostID,
		Title:         req.Title,
		Description:   req.Description,
		PropertyType:  req.PropertyType,
		Country:       req.Country,
		City:          req.City,
		Address:       req.Address,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		PricePerNight: req.PricePerNight,
		Currency:      req.Currency,
		MaxGuests:     req.MaxGuests,
		Bedrooms:      req.Bedrooms,
		Bathrooms:     req.Bathrooms,
		Amenities:     req.Amenities,
		Images:        req.Images,
	}

	if err := s.repo.Create(ctx, l); err != nil {
		return nil, fmt.Errorf("create listing: %w", err)
	}
	return l, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Listing, error) {
	l, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, respond.ErrNotFound
	}
	return l, nil
}

func (s *Service) Search(ctx context.Context, p SearchParams) (*ListingsResponse, error) {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}

	listings, total, err := s.repo.Search(ctx, p)
	if err != nil {
		return nil, err
	}

	return &ListingsResponse{
		Data:   listings,
		Total:  total,
		Limit:  p.Limit,
		Offset: p.Offset,
		Page:   p.Offset/p.Limit + 1,
	}, nil
}

func (s *Service) Update(ctx context.Context, id, requesterID, requesterRole string, req UpdateListingRequest) (*Listing, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, respond.ErrNotFound
	}
	if existing.HostID != requesterID && requesterRole != "admin" {
		return nil, respond.ErrForbidden
	}

	updated, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id, requesterID, requesterRole string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return respond.ErrNotFound
	}
	if existing.HostID != requesterID && requesterRole != "admin" {
		return respond.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

func validateCreate(req CreateListingRequest) error {
	if req.Title == "" || req.Description == "" || req.PropertyType == "" ||
		req.Country == "" || req.City == "" || req.Address == "" {
		return &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "title, description, property_type, country, city and address are required"}
	}
	if req.PricePerNight <= 0 {
		return &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "price_per_night must be greater than 0"}
	}
	if req.MaxGuests <= 0 {
		return &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "max_guests must be greater than 0"}
	}

	validTypes := map[string]bool{"apartment": true, "house": true, "villa": true, "cabin": true, "studio": true}
	if !validTypes[req.PropertyType] {
		return &respond.AppError{Code: http.StatusUnprocessableEntity, Message: "property_type must be one of: apartment, house, villa, cabin, studio"}
	}
	return nil
}
