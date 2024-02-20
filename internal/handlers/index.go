package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/repository"
	"github.com/sankaire/orders-api/internal/db/internal/utils"
	"net/http"
)

func CreateCustomerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method not allowed")
		return
	}

	var customer struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&customer); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, "Invalid request body")
		return
	}
	hashedPassword, err := utils.EncryptPassword(customer.Password)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "An error occured")
		return
	}
	customerID, err := repository.CreateCustomer(customer.Name, customer.Phone, customer.Email, hashedPassword)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error creating customer: %v", err)
		return
	}

	name, email, phone, err := repository.ReadCustomer(customerID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Error reading customer: %v", err)
		return
	}
	createdCustomer := struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}{
		ID:    customerID,
		Name:  name,
		Phone: phone,
		Email: email,
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

func LoginCustomerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method not allowed")
		return
	}
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	id, email, phone, password, err := repository.ReadCustomerByEmail(creds.Email)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = utils.ComparePassword(password, creds.Password)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	Payload := struct {
		ID    int    `json:"ID"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}{
		ID:    id,
		Email: email,
		Phone: phone,
	}
	accessToken, err := utils.CreateToken(Payload)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResponse, err := json.Marshal(accessToken)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	response.Write(jsonResponse)
}
