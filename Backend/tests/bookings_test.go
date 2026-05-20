package tests

import (
	"net/http"
	"testing"
)

func TestCreateBooking(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	hostToken := registerAndLogin(t, "bookhost@example.com", "password123")
	guestToken := registerAndLogin(t, "bookguest@example.com", "password123")
	listingID := createListing(t, hostToken)

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
			"check_in":   "2026-12-01",
			"check_out":  "2026-12-05",
		}, guestToken)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected 201, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["total_price"].(float64) != 340.00 {
			t.Errorf("expected total_price 340.00 (85 x 4 nights), got %v", body["total_price"])
		}
		if int(body["total_nights"].(float64)) != 4 {
			t.Errorf("expected total_nights 4, got %v", body["total_nights"])
		}
		if body["status"] != "pending" {
			t.Errorf("expected status pending, got %v", body["status"])
		}
	})

	t.Run("double booking overlapping dates", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
			"check_in":   "2026-12-03",
			"check_out":  "2026-12-07",
		}, guestToken)
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected 409, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("check_out before check_in", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
			"check_in":   "2026-12-10",
			"check_out":  "2026-12-08",
		}, guestToken)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("same day check_in and check_out", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
			"check_in":   "2026-12-20",
			"check_out":  "2026-12-20",
		}, guestToken)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("past date", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
			"check_in":   "2020-01-01",
			"check_out":  "2020-01-05",
		}, guestToken)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("missing fields", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
		}, guestToken)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": listingID,
			"check_in":   "2026-11-01",
			"check_out":  "2026-11-05",
		}, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("listing not found", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/bookings", map[string]any{
			"listing_id": "00000000-0000-0000-0000-000000000000",
			"check_in":   "2026-11-01",
			"check_out":  "2026-11-05",
		}, guestToken)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestCancelBooking(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	hostToken := registerAndLogin(t, "cancelhost@example.com", "password123")
	guestToken := registerAndLogin(t, "cancelguest@example.com", "password123")
	otherToken := registerAndLogin(t, "cancelother@example.com", "password123")
	listingID := createListing(t, hostToken)

	bookingResp := do(t, http.MethodPost, "/api/bookings", map[string]any{
		"listing_id": listingID,
		"check_in":   "2026-11-01",
		"check_out":  "2026-11-05",
	}, guestToken)
	var booking map[string]any
	decode(t, bookingResp, &booking)
	bookingID := booking["id"].(string)

	t.Run("wrong user", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/bookings/"+bookingID+"/cancel", nil, otherToken)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/bookings/"+bookingID+"/cancel", nil, guestToken)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["status"] != "cancelled" {
			t.Errorf("expected status cancelled, got %v", body["status"])
		}
	})

	t.Run("already cancelled", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/bookings/"+bookingID+"/cancel", nil, guestToken)
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected 409, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("not found", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/bookings/00000000-0000-0000-0000-000000000000/cancel", nil, guestToken)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestListMyBookings(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	hostToken := registerAndLogin(t, "listbookhost@example.com", "password123")
	guestToken := registerAndLogin(t, "listbookguest@example.com", "password123")
	listingID := createListing(t, hostToken)

	do(t, http.MethodPost, "/api/bookings", map[string]any{
		"listing_id": listingID,
		"check_in":   "2026-10-01",
		"check_out":  "2026-10-05",
	}, guestToken).Body.Close()

	t.Run("returns my bookings", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/bookings", nil, guestToken)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 1 {
			t.Errorf("expected 1 booking, got %v", body["total"])
		}
	})

	t.Run("does not return other users bookings", func(t *testing.T) {
		otherToken := registerAndLogin(t, "listbookother@example.com", "password123")
		resp := do(t, http.MethodGet, "/api/bookings", nil, otherToken)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 0 {
			t.Errorf("expected 0 bookings for other user, got %v", body["total"])
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/bookings", nil, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}
