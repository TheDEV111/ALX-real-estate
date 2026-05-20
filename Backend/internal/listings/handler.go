package listings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/auth"
	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit := intOrDefault(r.URL.Query().Get("limit"), 20)
	offset := intOrDefault(r.URL.Query().Get("offset"), 0)

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page := intOrDefault(pageStr, 1); page > 0 {
			offset = (page - 1) * limit
		}
	}

	p := SearchParams{
		Query:    nullableString(r.URL.Query().Get("q")),
		City:     nullableString(r.URL.Query().Get("city")),
		MinPrice: nullableFloat(r.URL.Query().Get("min_price")),
		MaxPrice: nullableFloat(r.URL.Query().Get("max_price")),
		Guests:   nullableInt(r.URL.Query().Get("guests")),
		CheckIn:  nullableString(r.URL.Query().Get("check_in")),
		CheckOut: nullableString(r.URL.Query().Get("check_out")),
		Limit:    limit,
		Offset:   offset,
	}

	if amenities := r.URL.Query()["amenities"]; len(amenities) > 0 {
		p.Amenities = amenities
	}

	result, err := h.svc.Search(r.Context(), p)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	listing, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, listing)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	var req CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err)
		return
	}

	listing, err := h.svc.Create(r.Context(), claims.UserID, req)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, listing)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	var req UpdateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err)
		return
	}

	listing, err := h.svc.Update(r.Context(), id, claims.UserID, claims.Role, req)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, listing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id, claims.UserID, claims.Role); err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "listing deleted"})
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nullableFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

func nullableInt(s string) *int {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func intOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
