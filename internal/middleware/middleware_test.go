package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sankaire/orders-api/internal/db/internal/middleware"
)

func TestAuthenticate_NoToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/protected", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := middleware.Authenticate(handler)
	authMiddleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := middleware.Authenticate(handler)
	authMiddleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
