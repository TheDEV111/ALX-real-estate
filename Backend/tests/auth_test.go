package tests

import (
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/auth/register", map[string]string{
			"email": "register@example.com", "password": "password123", "full_name": "Henry Doe",
		}, "")
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected 201, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["access_token"] == nil {
			t.Error("expected access_token in response")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		do(t, http.MethodPost, "/api/auth/register", map[string]string{
			"email": "dup@example.com", "password": "password123", "full_name": "User",
		}, "").Body.Close()

		resp := do(t, http.MethodPost, "/api/auth/register", map[string]string{
			"email": "dup@example.com", "password": "password123", "full_name": "User",
		}, "")
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected 409, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("missing fields", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/auth/register", map[string]string{
			"email": "missing@example.com",
		}, "")
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestLogin(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	do(t, http.MethodPost, "/api/auth/register", map[string]string{
		"email": "login@example.com", "password": "password123", "full_name": "User",
	}, "").Body.Close()

	t.Run("success", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/auth/login", map[string]string{
			"email": "login@example.com", "password": "password123",
		}, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["access_token"] == nil {
			t.Error("expected access_token in response")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/auth/login", map[string]string{
			"email": "login@example.com", "password": "wrongpassword",
		}, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("user not found", func(t *testing.T) {
		resp := do(t, http.MethodPost, "/api/auth/login", map[string]string{
			"email": "ghost@example.com", "password": "password123",
		}, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}

func TestMe(t *testing.T) {
	t.Cleanup(func() { truncate(t) })

	token := registerAndLogin(t, "me@example.com", "password123")

	t.Run("with valid token", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/users/me", nil, token)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		decode(t, resp, &body)
		if body["email"] != "me@example.com" {
			t.Errorf("expected email me@example.com, got %v", body["email"])
		}
	})

	t.Run("without token", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/users/me", nil, "")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("invalid token", func(t *testing.T) {
		resp := do(t, http.MethodGet, "/api/users/me", nil, "not.a.valid.token")
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})
}
