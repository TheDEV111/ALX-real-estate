package tests

import (
	"context"
	"net/http"
	"testing"
)

func completeBooking(t *testing.T, bookingID string) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		"UPDATE bookings SET status = 'completed' WHERE id = $1", bookingID)
	if err != nil {
		t.Fatalf("complete booking: %v", err)
	}
}

func TestCreateReview(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	hostToken := registerAndLogin(t, "revhost@example.com", "password123")
	guestToken := registerAndLogin(t, "revguest@example.com", "password123")
	otherToken := registerAndLogin(t, "revother@example.com", "password123")
	listingID := createListing(t, hostToken)

	bookingResp := do(t, http.MethodPost, "/api/bookings", map[string]any{
		"listing_id": listingID,
		"check_in":   "2026-08-01",
		"check_out":  "2026-08-05",
	}, guestToken)
	var booking map[string]any
	decode(t, bookingResp, &booking)
	bookingID := booking["id"].(string)

	t.Run("cannot review pending booking", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
			"rating": 5, "comment": "great",
		}, guestToken)
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected 409, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	completeBooking(t, bookingID)

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
			"rating": 4, "comment": "lovely stay",
		}, guestToken)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected 201, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["rating"].(float64)) != 4 {
			t.Errorf("expected rating 4, got %v", body["rating"])
		}
		if body["listing_id"] != listingID {
			t.Errorf("expected listing_id %s, got %v", listingID, body["listing_id"])
		}
	})

	t.Run("duplicate review", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
			"rating": 3,
		}, guestToken)
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected 409, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("wrong guest forbidden", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
			"rating": 5,
		}, otherToken)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("invalid rating", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
			"rating": 6,
		}, guestToken)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("booking not found", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/00000000-0000-0000-0000-000000000000/review", map[string]any{
			"rating": 5,
		}, guestToken)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
			"rating": 5,
		}, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestListReviews(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	hostToken := registerAndLogin(t, "listrevhost@example.com", "password123")
	guestToken := registerAndLogin(t, "listrevguest@example.com", "password123")
	listingID := createListing(t, hostToken)

	bookingResp := do(t, http.MethodPost, "/api/bookings", map[string]any{
		"listing_id": listingID,
		"check_in":   "2026-07-01",
		"check_out":  "2026-07-04",
	}, guestToken)
	var booking map[string]any
	decode(t, bookingResp, &booking)
	bookingID := booking["id"].(string)
	completeBooking(t, bookingID)

	do(t, http.MethodPost, "/api/bookings/"+bookingID+"/review", map[string]any{
		"rating": 5, "comment": "perfect",
	}, guestToken).Body.Close()

	t.Run("returns reviews with average", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+listingID+"/reviews", nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 1 {
			t.Errorf("expected 1 review, got %v", body["total"])
		}
		if body["average_rating"].(float64) != 5.0 {
			t.Errorf("expected average_rating 5.0, got %v", body["average_rating"])
		}
	})

	t.Run("empty for listing with no reviews", func(t *testing.T) {
		otherToken := registerAndLogin(t, "listrevother@example.com", "password123")
		otherListingID := createListing(t, otherToken)

		resp := do(t, http.MethodGet, "/api/listings/"+otherListingID+"/reviews", nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 0 {
			t.Errorf("expected 0 reviews, got %v", body["total"])
		}
		if body["average_rating"].(float64) != 0.0 {
			t.Errorf("expected average_rating 0.0, got %v", body["average_rating"])
		}
	})
}
