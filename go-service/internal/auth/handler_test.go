package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"goservice/internal/client"
	"goservice/internal/models"
	"goservice/internal/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Mock IBackend ---
type mockBackend struct {
	loginFn func(ctx context.Context, username, password string) ([]*http.Cookie, error)
}

func (m *mockBackend) Login(ctx context.Context, username, password string) ([]*http.Cookie, error) {
	return m.loginFn(ctx, username, password)
}
func (m *mockBackend) GetStudentByID(ctx context.Context, id int, cookies []*http.Cookie) (*models.Student, error) {
	return nil, nil // not needed for this test
}

func TestHandler_Login_Success(t *testing.T) {
	mock := &mockBackend{
		loginFn: func(_ context.Context, username, password string) ([]*http.Cookie, error) {
			if username == "user" && password == "pass" {
				return []*http.Cookie{
					{Name: client.CSFRTokenName, Value: "csrf123"},
					{Name: client.AccesTokenName, Value: "access123"},
					{Name: client.RefreshTokenName, Value: "refresh123"},
				}, nil
			}
			return nil, errors.New("invalid credentials")
		},
	}
	h := NewHandler(mock)
	r := h.Routes()

	creds := map[string]string{"username": "user", "password": "pass"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	// Check cookies are set
	cookies := resp.Cookies()
	if len(cookies) != 3 {
		t.Errorf("expected 3 cookies, got %d", len(cookies))
	}
	// Check response body
	var out response.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	msg, ok := out.Data.(map[string]interface{})
	if !ok {
		t.Errorf("expected map data, got %T", out.Data)
	}
	if msg["message"] != "Login successful" {
		t.Errorf("unexpected message: %v", msg["message"])
	}
}

func TestHandler_Login_BadJSON(t *testing.T) {
	mock := &mockBackend{
		loginFn: func(_ context.Context, username, password string) ([]*http.Cookie, error) {
			return nil, nil
		},
	}
	h := NewHandler(mock)
	r := h.Routes()

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("{bad json"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var out response.APIResponse
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.Error == "" {
		t.Error("expected error in response")
	}
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	mock := &mockBackend{
		loginFn: func(_ context.Context, username, password string) ([]*http.Cookie, error) {
			return nil, errors.New("invalid credentials")
		},
	}
	h := NewHandler(mock)
	r := h.Routes()

	creds := map[string]string{"username": "user", "password": "badpass"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	var out response.APIResponse
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.Error != "invalid credentials" {
		t.Errorf("expected error 'invalid credentials', got %q", out.Error)
	}
}
