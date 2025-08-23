package student

import (
	"bytes"
	"context"
	"errors"
	"goservice/internal/models"
	"io"
	"net/http"
	"testing"
	"time"
)

// --- Mock BackendClient ---
type mockBackendClient struct {
	loginFn        func(ctx context.Context, username, password string) ([]*http.Cookie, error)
	getStudentByID func(ctx context.Context, id int, cookies []*http.Cookie) (*models.Student, error)
}

func (m *mockBackendClient) Login(ctx context.Context, username, password string) ([]*http.Cookie, error) {
	return m.loginFn(ctx, username, password)
}
func (m *mockBackendClient) GetStudentByID(ctx context.Context, id int, cookies []*http.Cookie) (*models.Student, error) {
	return m.getStudentByID(ctx, id, cookies)
}

// --- Fakes for client.BackendClient interface ---
func fakeBackendClient(loginFn func(context.Context, string, string) ([]*http.Cookie, error),
	getStudentByIDFn func(context.Context, int, []*http.Cookie) (*models.Student, error)) *mockBackendClient {
	return &mockBackendClient{
		loginFn:        loginFn,
		getStudentByID: getStudentByIDFn,
	}
}

// --- Tests ---

func TestService_Login(t *testing.T) {
	wantCookies := []*http.Cookie{
		{Name: "csrfToken", Value: "csrf123"},
	}
	svc := &service{
		backend: fakeBackendClient(
			func(_ context.Context, u, p string) ([]*http.Cookie, error) {
				if u == "user" && p == "pass" {
					return wantCookies, nil
				}
				return nil, errors.New("invalid credentials")
			},
			nil, // not needed for this test
		),
	}
	got, err := svc.Login(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Value != "csrf123" {
		t.Errorf("expected cookies %v, got %v", wantCookies, got)
	}
}

func TestService_GetStudent(t *testing.T) {
	wantStudent := &models.Student{
		ID:                 1,
		Name:               "Test Student",
		Email:              "test@example.com",
		SystemAccess:       true,
		Phone:              "1234567890",
		Gender:             "M",
		DOB:                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Class:              "10",
		Section:            "A",
		Roll:               1,
		FatherName:         "Father",
		FatherPhone:        "1111111111",
		MotherName:         "Mother",
		MotherPhone:        "2222222222",
		GuardianName:       "Guardian",
		GuardianPhone:      "3333333333",
		RelationOfGuardian: "Uncle",
		CurrentAddress:     "Current Addr",
		PermanentAddress:   "Perm Addr",
		AdmissionDate:      time.Date(2018, 6, 10, 0, 0, 0, 0, time.UTC),
		ReporterName:       "Reporter",
	}
	svc := &service{
		backend: fakeBackendClient(
			nil,
			func(_ context.Context, id int, _ []*http.Cookie) (*models.Student, error) {
				if id == 1 {
					return wantStudent, nil
				}
				return nil, errors.New("not found")
			},
		),
	}
	got, err := svc.GetStudent(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != wantStudent.ID || got.Name != wantStudent.Name {
		t.Errorf("expected %v, got %v", wantStudent, got)
	}
}

func TestService_GenerateReport(t *testing.T) {
	wantStudent := &models.Student{
		ID:                 1,
		Name:               "Test Student",
		Email:              "test@example.com",
		SystemAccess:       true,
		Phone:              "1234567890",
		Gender:             "M",
		DOB:                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Class:              "10",
		Section:            "A",
		Roll:               1,
		FatherName:         "Father",
		FatherPhone:        "1111111111",
		MotherName:         "Mother",
		MotherPhone:        "2222222222",
		GuardianName:       "Guardian",
		GuardianPhone:      "3333333333",
		RelationOfGuardian: "Uncle",
		CurrentAddress:     "Current Addr",
		PermanentAddress:   "Perm Addr",
		AdmissionDate:      time.Date(2018, 6, 10, 0, 0, 0, 0, time.UTC),
		ReporterName:       "Reporter",
	}
	svc := &service{
		backend: fakeBackendClient(
			nil,
			func(_ context.Context, id int, _ []*http.Cookie) (*models.Student, error) {
				return wantStudent, nil
			},
		),
	}
	rep, err := svc.GenerateReport(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Make sure Output does not error and writes some PDF bytes
	buf := new(bytes.Buffer)
	if err := rep.Output(buf); err != nil && err != io.EOF {
		t.Fatalf("report output error: %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("expected some report output, got 0 bytes")
	}
}
