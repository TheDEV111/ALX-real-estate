package tests

import (
	"net/http"
	"testing"
)

func TestListListingBookings(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	hostToken := registerAndLogin(t, "dashhost@example.com", "password123")
	guestToken := registerAndLogin(t, "dashguest@example.com", "password123")
	otherToken := registerAndLogin(t, "dashother@example.com", "password123")
	listingID := createListing(t, hostToken)

	do(t, http.MethodPost, "/api/bookings", map[string]any{
		"listing_id": listingID,
		"check_in":   "2026-09-01",
		"check_out":  "2026-09-05",
	}, guestToken).Body.Close()

	t.Run("host can see bookings", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+listingID+"/bookings", nil, hostToken)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 1 {
			t.Errorf("expected 1 booking, got %v", body["total"])
		}
	})

	t.Run("non-owner is forbidden", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+listingID+"/bookings", nil, otherToken)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("guest is forbidden", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+listingID+"/bookings", nil, guestToken)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+listingID+"/bookings", nil, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("listing not found", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/00000000-0000-0000-0000-000000000000/bookings", nil, hostToken)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}
