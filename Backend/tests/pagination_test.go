package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestListingsPagination(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	token := registerAndLogin(t, "paghost@example.com", "password123")

	for i := 0; i < 3; i++ {
		body := sampleListing()
		body["title"] = fmt.Sprintf("Listing %d", i+1)
		resp := do(t, http.MethodPost, "/api/listings", body, token)
		resp.Body.Close()
	}

	t.Run("default limit=20 returns all 3", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings", nil, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if int(body["total"].(float64)) != 3 {
			t.Errorf("expected total 3, got %v", body["total"])
		}
		if len(body["data"].([]any)) != 3 {
			t.Errorf("expected 3 items, got %d", len(body["data"].([]any)))
		}
		if int(body["page"].(float64)) != 1 {
			t.Errorf("expected page 1, got %v", body["page"])
		}
	})

	t.Run("limit=2 page=1 returns 2 items", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings?limit=2&page=1", nil, "")
		var body map[string]any
		decode(t, resp, &body)
		if len(body["data"].([]any)) != 2 {
			t.Errorf("expected 2 items on page 1, got %d", len(body["data"].([]any)))
		}
		if int(body["total"].(float64)) != 3 {
			t.Errorf("expected total 3, got %v", body["total"])
		}
		if int(body["page"].(float64)) != 1 {
			t.Errorf("expected page 1, got %v", body["page"])
		}
	})

	t.Run("limit=2 page=2 returns 1 item", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings?limit=2&page=2", nil, "")
		var body map[string]any
		decode(t, resp, &body)
		if len(body["data"].([]any)) != 1 {
			t.Errorf("expected 1 item on page 2, got %d", len(body["data"].([]any)))
		}
		if int(body["page"].(float64)) != 2 {
			t.Errorf("expected page 2, got %v", body["page"])
		}
	})

	t.Run("limit=2 page=3 returns 0 items", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings?limit=2&page=3", nil, "")
		var body map[string]any
		decode(t, resp, &body)
		if len(body["data"].([]any)) != 0 {
			t.Errorf("expected 0 items on page 3, got %d", len(body["data"].([]any)))
		}
		if int(body["total"].(float64)) != 3 {
			t.Errorf("expected total still 3, got %v", body["total"])
		}
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/listings?limit=999", nil, "")
		var body map[string]any
		decode(t, resp, &body)
		if int(body["limit"].(float64)) != 100 {
			t.Errorf("expected limit capped at 100, got %v", body["limit"])
		}
	})
}

func TestBodySizeLimit(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	token := registerAndLogin(t, "sizehost@example.com", "password123")

	t.Run("body over 1MB is rejected", func(t *testing.T) {
		bigPayload := make([]byte, 1<<20+1)
		for i := range bigPayload {
			bigPayload[i] = 'a'
		}

		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/listings",
			bytes.NewReader(append([]byte(`{"title":"`), append(bigPayload, []byte(`"}`)...)...)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusRequestEntityTooLarge {
			t.Errorf("expected 413, got %d", resp.StatusCode)
		}
	})
}
