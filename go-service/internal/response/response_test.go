package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}

	JSON(rec, http.StatusOK, data)

	resp := rec.Result()
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check content type
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected content-type application/json, got %s", resp.Header.Get("Content-Type"))
	}

	// Check body
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	// Data should match
	got, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map data, got %T", apiResp.Data)
	}
	if got["foo"] != "bar" {
		t.Errorf("expected foo=bar, got foo=%v", got["foo"])
	}
	if apiResp.Error != "" {
		t.Errorf("expected error to be empty, got %q", apiResp.Error)
	}
}

func TestError(t *testing.T) {
	rec := httptest.NewRecorder()
	errMsg := "something bad happened"

	Error(rec, http.StatusBadRequest, errors.New(errMsg))

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected content-type application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if apiResp.Error != errMsg {
		t.Errorf("expected error %q, got %q", errMsg, apiResp.Error)
	}
	if apiResp.Data != nil {
		t.Errorf("expected data to be nil, got %v", apiResp.Data)
	}
}
