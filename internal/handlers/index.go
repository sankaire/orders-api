package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/repository"
	"net/http"
)

func CreateCustomerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method not allowed")
		return
	}

	var customer struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	if err := json.NewDecoder(request.Body).Decode(&customer); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, "Invalid request body")
		return
	}

	customerID, err := repository.CreateCustomer(customer.Name, customer.Code)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error creating customer: %v", err)
		return
	}

	name, code, err := repository.ReadCustomer(customerID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error reading customer: %v", err)
		return
	}
	createdCustomer := struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}{
		ID:   customerID,
		Name: name,
		Code: code,
	}
	jsonResponse, err := json.Marshal(createdCustomer)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error marshaling JSON: %v", err)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	response.Write(jsonResponse)
}
func ReadCustomersHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method not allowed")
		return
	}
	customers, err := repository.ReadCustomers()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error reading customer: %v", err)
		return
	}
	jsonResponse, err := json.Marshal(customers)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error marshaling JSON: %v", err)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	response.Write(jsonResponse)
}
