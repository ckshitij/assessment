package client

import (
	"context"
	"encoding/json"
	"goservice/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Helper to create a student for test responses
func sampleStudent() *models.Student {
	return &models.Student{
		ID:                 1,
		Name:               "Alice",
		Email:              "alice@example.com",
		SystemAccess:       true,
		Phone:              "1234567890",
		Gender:             "F",
		Class:              "10",
		Section:            "A",
		Roll:               5,
		FatherName:         "Bob",
		FatherPhone:        "1111111111",
		MotherName:         "Carol",
		MotherPhone:        "2222222222",
		GuardianName:       "Eve",
		GuardianPhone:      "3333333333",
		RelationOfGuardian: "Aunt",
		CurrentAddress:     "Current Addr",
		PermanentAddress:   "Perm Addr",
		AdmissionDate:      time.Date(2018, 6, 10, 0, 0, 0, 0, time.UTC),
		ReporterName:       "Reporter",
		DOB:                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestBackendClient_Login_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		// Simulate setting cookies
		http.SetCookie(w, &http.Cookie{Name: CSFRTokenName, Value: "csrf123"})
		http.SetCookie(w, &http.Cookie{Name: AccesTokenName, Value: "access123"})
		http.SetCookie(w, &http.Cookie{Name: RefreshTokenName, Value: "refresh123"})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	client := NewBackendClient(ts.URL)
	cookies, err := client.Login(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cookies) != 3 {
		t.Errorf("expected 3 cookies, got %d", len(cookies))
	}
	// check each cookie exists
	expected := map[string]bool{
		CSFRTokenName:    false,
		AccesTokenName:   false,
		RefreshTokenName: false,
	}
	for _, c := range cookies {
		expected[c.Name] = true
	}
	for name, got := range expected {
		if !got {
			t.Errorf("expected cookie %s to be present", name)
		}
	}
}

func TestBackendClient_Login_Failure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
	}))
	defer ts.Close()

	client := NewBackendClient(ts.URL)
	cookies, err := client.Login(context.Background(), "user", "badpass")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Errorf("expected error to contain 'invalid credentials', got %v", err)
	}
	if len(cookies) != 0 {
		t.Errorf("expected no cookies, got %d", len(cookies))
	}
}

func TestBackendClient_GetStudentByID_Success(t *testing.T) {
	stu := sampleStudent()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for correct path and header
		if !strings.HasSuffix(r.URL.Path, "/api/v1/students/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-csrf-token") != "csrf123" {
			t.Errorf("expected CSRF token header, got %s", r.Header.Get("x-csrf-token"))
		}
		// Return JSON
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(stu)
	}))
	defer ts.Close()

	client := NewBackendClient(ts.URL)
	cookie := &http.Cookie{Name: CSFRTokenName, Value: "csrf123"}
	got, err := client.GetStudentByID(context.Background(), 1, []*http.Cookie{cookie})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != stu.ID || got.Name != stu.Name {
		t.Errorf("expected student %+v, got %+v", stu, got)
	}
}

func TestBackendClient_GetStudentByID_Failure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer ts.Close()

	client := NewBackendClient(ts.URL)
	cookie := &http.Cookie{Name: CSFRTokenName, Value: "csrf123"}
	got, err := client.GetStudentByID(context.Background(), 2, []*http.Cookie{cookie})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected error to contain 'not found', got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil student, got %+v", got)
	}
}

// Optionally, add test for failure on decoding JSON
func TestBackendClient_GetStudentByID_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "not json")
	}))
	defer ts.Close()
	client := NewBackendClient(ts.URL)
	cookie := &http.Cookie{Name: CSFRTokenName, Value: "csrf123"}
	got, err := client.GetStudentByID(context.Background(), 1, []*http.Cookie{cookie})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode student") {
		t.Errorf("expected decode error, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil student, got %+v", got)
	}
}
