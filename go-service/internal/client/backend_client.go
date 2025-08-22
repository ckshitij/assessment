package client

import (
	"context"
	"encoding/json"
	"fmt"
	"goservice/internal/models"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	accesTokenName   = "accessToken"
	refreshTokenName = "refreshToken"
	csrfTokenName    = "csrfToken"
)

type BackendClient struct {
	BaseURL string
	Client  *http.Client
}

func NewBackendClient(baseURL string) *BackendClient {
	return &BackendClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (b *BackendClient) Login(ctx context.Context, username, password string) ([]*http.Cookie, error) {
	loginURL := fmt.Sprintf("%s/api/v1/auth/login", b.BaseURL)
	payload := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

	log.Println(loginURL)

	req, _ := http.NewRequestWithContext(ctx, "POST", loginURL, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	var cookies []*http.Cookie
	// Extract csrfToken, accessToken and refreshToken from cookies
	for _, c := range resp.Cookies() {
		if c.Name == csrfTokenName && c.Value != "" {
			cookies = append(cookies, c)
		}
		if c.Name == accesTokenName && c.Value != "" {
			cookies = append(cookies, c)
		}
		if c.Name == refreshTokenName && c.Value != "" {
			cookies = append(cookies, c)
		}
	}

	return cookies, nil
}

func (b *BackendClient) GetStudentByID(ctx context.Context, id int, rawCookies []*http.Cookie) (*models.Student, error) {
	url := fmt.Sprintf("%s/api/v1/students/%d", b.BaseURL, id)
	var csrfToken string

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	for _, c := range rawCookies {
		if c.Name == csrfTokenName {
			csrfToken = c.Value
		}
		req.AddCookie(c)
	}
	req.Header.Set("x-csrf-token", csrfToken)

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get student: %s", string(body))
	}

	var student models.Student
	if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
		return nil, fmt.Errorf("failed to decode student: %v", err)
	}

	return &student, nil
}
