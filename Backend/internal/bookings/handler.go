package bookings

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err)
		return
	}

	booking, err := h.svc.Create(r.Context(), claims.UserID, req)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, booking)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	limit := intOrDefault(r.URL.Query().Get("limit"), 20)
	offset := intOrDefault(r.URL.Query().Get("offset"), 0)

	result, err := h.svc.ListMine(r.Context(), claims.UserID, limit, offset)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	booking, err := h.svc.GetByID(r.Context(), id, claims.UserID, claims.Role)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, booking)
}

func (h *Handler) ListForListing(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	listingID := chi.URLParam(r, "id")
	limit := intOrDefault(r.URL.Query().Get("limit"), 20)
	offset := intOrDefault(r.URL.Query().Get("offset"), 0)

	result, err := h.svc.ListForListing(r.Context(), listingID, claims.UserID, claims.Role, limit, offset)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	booking, err := h.svc.Cancel(r.Context(), id, claims.UserID, claims.Role)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, booking)
}

func intOrDefault(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
