package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheDEV111/ALX-real-estate/backend/internal/respond"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err)
		return
	}
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		respond.Error(w, &respond.AppError{Code: 422, Message: "email, password and full_name are required"})
		return
	}

	resp, refreshToken, err := h.svc.Register(r.Context(), req)
	if err != nil {
		respond.Error(w, err)
		return
	}

	setRefreshCookie(w, refreshToken)
	respond.JSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err)
		return
	}

	resp, refreshToken, err := h.svc.Login(r.Context(), req)
	if err != nil {
		respond.Error(w, err)
		return
	}

	setRefreshCookie(w, refreshToken)
	respond.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}

	resp, newRefreshToken, err := h.svc.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		respond.Error(w, err)
		return
	}

	setRefreshCookie(w, newRefreshToken)
	respond.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		respond.Error(w, respond.ErrUnauthorized)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{
		"id":    claims.UserID,
		"email": claims.Email,
		"role":  claims.Role,
	})
}

func setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/auth/refresh",
		MaxAge:   7 * 24 * 60 * 60,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}
