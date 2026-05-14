package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string { return e.Message }

var (
	ErrNotFound     = &AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrUnauthorized = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden    = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrConflict     = &AppError{Code: http.StatusConflict, Message: "resource conflict"}
	ErrInvalidInput = &AppError{Code: http.StatusUnprocessableEntity, Message: "invalid input"}
)

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RespondError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		RespondJSON(w, appErr.Code, appErr)
		return
	}
	log.Printf("internal error: %v", err)
	RespondJSON(w, http.StatusInternalServerError, &AppError{Message: "internal server error"})
}
