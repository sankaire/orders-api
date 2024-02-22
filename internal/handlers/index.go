package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/repository"
	"github.com/sankaire/orders-api/internal/db/internal/utils"
	"net/http"
	"strconv"
	time2 "time"
)

func CreateCustomerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		utils.WriteResponse(response, http.StatusBadRequest, false, "Method not allowed", nil)
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
		utils.WriteResponse(response, http.StatusBadRequest, false, err.Error(), nil)
		return
	}
	hashedPassword, err := utils.EncryptPassword(customer.Password)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		utils.WriteResponse(response, http.StatusBadRequest, false, err.Error(), nil)
		return
	}
	customerID, err := repository.CreateCustomer(customer.Name, customer.Phone, customer.Email, hashedPassword)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, "An error occurred", nil)
		return
	}

	name, email, phone, err := repository.ReadCustomer(customerID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		utils.WriteResponse(response, http.StatusInternalServerError, false, "An error occurred", nil)
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
	var customerResults interface{}
	customerResponse, err := json.Marshal(createdCustomer)
	err = json.Unmarshal(customerResponse, &customerResults)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	utils.WriteResponse(response, http.StatusCreated, true, "Account created successfully", customerResults)
	return
}

func LoginCustomerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		utils.WriteResponse(response, http.StatusBadRequest, false, "Method Not allowed", nil)
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
	var loginResult interface{}
	loginResponse, err := json.Marshal(accessToken)
	err = json.Unmarshal(loginResponse, &loginResult)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	utils.WriteResponse(response, http.StatusOK, true, "Login successful", loginResult)
	return
}
func CreateCustomerOrder(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		utils.WriteResponse(response, http.StatusBadRequest, false, "Method not allowed", nil)
		return
	}
	var Order struct {
		Item   string `json:"item"`
		Amount int    `json:"amount"`
	}
	err := json.NewDecoder(request.Body).Decode(&Order)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	phone := ""
	if phoneValue := request.Context().Value("phone"); phoneValue != nil {
		if phoneStr, ok := phoneValue.(string); ok {
			phone = phoneStr
		}
	}
	id := request.Context().Value("id")

	orderID, err := repository.CreateOrder(id, Order.Item, Order.Amount)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	_, _, item, amount, time, err := repository.ReadCustomerOrder(orderID)
	var Payload = struct {
		ID         int64      `json:"ID"`
		CustomerID any        `json:"customerID"`
		Item       string     `json:"item"`
		Amount     int64      `json:"amount"`
		Time       time2.Time `json:"time"`
	}{
		ID:         orderID,
		CustomerID: id,
		Item:       item,
		Amount:     amount,
		Time:       time,
	}
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	var ordersResult interface{}
	ordersResponse, err := json.Marshal(Payload)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	err = json.Unmarshal(ordersResponse, &ordersResult)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	err = utils.SendSms(phone, item, amount)
	if err != nil {
		fmt.Println(err)
	}
	utils.WriteResponse(response, http.StatusCreated, true, "Order created successfully", ordersResult)
	return
}
func ReadCustomerOrder(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		utils.WriteResponse(response, http.StatusBadRequest, false, "Method Not allowed", nil)
		return
	}
	orderID := request.URL.Query().Get("id")
	customerID := request.Context().Value("id")
	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	_, _, item, amount, time, err := repository.ReadCustomerOrder(int64(orderIDInt))
	var Payload = struct {
		ID         int64      `json:"ID"`
		CustomerID any        `json:"customerID"`
		Item       string     `json:"item"`
		Amount     int64      `json:"amount"`
		Time       time2.Time `json:"time"`
	}{
		ID:         int64(orderIDInt),
		CustomerID: customerID,
		Item:       item,
		Amount:     amount,
		Time:       time,
	}
	var orderResult interface{}
	orderResponse, err := json.Marshal(Payload)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	err = json.Unmarshal(orderResponse, &orderResult)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	utils.WriteResponse(response, http.StatusOK, true, "Order created successfully", orderResult)
	return
}
func ReadAllCustomerOrdersHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		utils.WriteResponse(response, http.StatusBadRequest, false, "Method Not allowed", nil)
		return
	}
	customerID := request.Context().Value("id")

	orders, err := repository.ReadAllCustomerOrder(customerID)
	if err != nil {
		utils.WriteResponse(response, http.StatusOK, false, "Customers not found", nil)
		return
	}
	orderResponse, err := json.Marshal(orders)
	var ordersResult interface{}

	err = json.Unmarshal(orderResponse, &ordersResult)
	if err != nil {
		utils.WriteResponse(response, http.StatusInternalServerError, false, err.Error(), nil)
	}
	utils.WriteResponse(response, http.StatusOK, true, "Orders fetched successfully", ordersResult)

}
