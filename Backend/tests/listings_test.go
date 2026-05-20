package tests

import (
	"net/http"
	"testing"
)

func sampleListing() map[string]any {
	return map[string]any{
		"title":           "Cozy Accra Apartment",
		"description":     "A beautiful apartment in the heart of Accra",
		"property_type":   "apartment",
		"country":         "Ghana",
		"city":            "Accra",
		"address":         "14 Liberation Road",
		"price_per_night": 85.00,
		"max_guests":      4,
		"bedrooms":        2,
		"bathrooms":       1.0,
		"amenities":       []string{"wifi", "pool"},
		"images":          []string{},
	}
}

func createListing(t *testing.T, token string) string {
	t.Helper()
	resp := do(t, http.MethodPost, "/api/listings", sampleListing(), token)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 creating listing, got %d", resp.StatusCode)
	}
	var body map[string]any
	decode(t, resp, &body)
	return body["id"].(string)
}

func TestListListings(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	t.Run("empty", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings", nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 0 {
			t.Errorf("expected total 0, got %v", body["total"])
		}
	})

	t.Run("with results", func(t *testing.T) {
		token := registerAndLogin(t, "listhost@example.com", "password123")
		createListing(t, token)

		resp := do(t, http.MethodGet, "/api/listings", nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 1 {
			t.Errorf("expected total 1, got %v", body["total"])
		}
	})

	t.Run("filter by city", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings?city=Accra", nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		data := body["data"].([]any)
		for _, item := range data {
			listing := item.(map[string]any)
			if listing["city"] != "Accra" {
				t.Errorf("expected city Accra, got %v", listing["city"])
			}
		}
	})
}

func TestCreateListing(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	token := registerAndLogin(t, "createhost@example.com", "password123")

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/listings", sampleListing(), token)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected 201, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["id"] == nil {
			t.Error("expected id in response")
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/listings", sampleListing(), "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("missing required fields", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/listings", map[string]any{"title": "only title"}, token)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("invalid property type", func(t *testing.T) {
		body := sampleListing()
		body["property_type"] = "castle"
		resp := do(t, http.MethodPost, "/api/listings", body, token)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("zero price", func(t *testing.T) {
		body := sampleListing()
		body["price_per_night"] = 0
		resp := do(t, http.MethodPost, "/api/listings", body, token)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestGetListingByID(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	token := registerAndLogin(t, "gethost@example.com", "password123")
	id := createListing(t, token)

	t.Run("found", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+id, nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["id"] != id {
			t.Errorf("expected id %s, got %v", id, body["id"])
		}
	})

	t.Run("not found", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/00000000-0000-0000-0000-000000000000", nil, "")
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestUpdateListing(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	ownerToken := registerAndLogin(t, "updateowner@example.com", "password123")
	otherToken := registerAndLogin(t, "updateother@example.com", "password123")
	id := createListing(t, ownerToken)

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/listings/"+id, map[string]any{
			"title": "Updated Title",
		}, ownerToken)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["title"] != "Updated Title" {
			t.Errorf("expected updated title, got %v", body["title"])
		}
	})

	t.Run("wrong owner", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/listings/"+id, map[string]any{
			"title": "Hijacked Title",
		}, otherToken)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("not found", func(t *testing.T) {
		resp := do(t, http.MethodPatch, "/api/listings/00000000-0000-0000-0000-000000000000", map[string]any{
			"title": "Ghost",
		}, ownerToken)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestDeleteListing(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	ownerToken := registerAndLogin(t, "deleteowner@example.com", "password123")
	otherToken := registerAndLogin(t, "deleteother@example.com", "password123")
	id := createListing(t, ownerToken)

	t.Run("wrong owner", func(t *testing.T) {
		resp := do(t, http.MethodDelete, "/api/listings/"+id, nil, otherToken)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodDelete, "/api/listings/"+id, nil, ownerToken)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("not found after delete", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings/"+id, nil, "")
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}
