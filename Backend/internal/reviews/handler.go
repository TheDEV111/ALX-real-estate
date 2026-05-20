package reviews

import (
	"encoding/json"
	"net/http"

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

	bookingID := chi.URLParam(r, "id")
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err)
		return
	}

	rev, err := h.svc.Create(r.Context(), bookingID, claims.UserID, req)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, rev)
}

func (h *Handler) ListForListing(w http.ResponseWriter, r *http.Request) {
	listingID := chi.URLParam(r, "id")
	result, err := h.svc.ListForListing(r.Context(), listingID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, result)
}
