package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sankaire/orders-api/internal/db/internal/handlers"
)

func TestCreateCustomerHandler(t *testing.T) {
	reqBody := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"phone":    "+1234567890",
		"password": "password",
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/create-customer", bytes.NewBuffer(reqJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handlers.CreateCustomerHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestLoginCustomerHandler(t *testing.T) {
	reqBody := map[string]string{
		"email":    "john@example.com",
		"password": "password",
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/login-customer", bytes.NewBuffer(reqJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handlers.LoginCustomerHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestCreateCustomerOrder(t *testing.T) {
	reqBody := map[string]interface{}{
		"item":   "Test Item",
		"amount": 100,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/create-customer-order", bytes.NewBuffer(reqJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handlers.CreateCustomerOrder(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestReadAllCustomerOrdersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/read-all-customer-orders", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handlers.ReadAllCustomerOrdersHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
