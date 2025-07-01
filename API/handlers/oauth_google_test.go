package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

// helper to create handler with dummy config
func newTestGoogleHandler() *GoogleOAuthHandler {
	cfg := &oauth2.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://example.com/auth",
			TokenURL: "http://example.com/token",
		},
	}
	return &GoogleOAuthHandler{Config: cfg}
}

func TestGoogleOAuthLoginRedirect(t *testing.T) {
	h := newTestGoogleHandler()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/forum/api/auth/google", nil)

	h.Login(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusTemporaryRedirect {
		t.Fatalf("expected %d, got %d", http.StatusTemporaryRedirect, status)
	}
	loc := rr.Header().Get("Location")
	if !strings.Contains(loc, "example.com/auth") {
		t.Fatalf("unexpected redirect location %q", loc)
	}
}

func TestGoogleOAuthCallbackMissingCode(t *testing.T) {
	h := newTestGoogleHandler()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/forum/api/auth/google/callback", nil)

	h.Callback(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, status)
	}
}
