package repository

import (
	"database/sql"
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/db"
)

type Customers struct {
	ID       int
	Name     string
	Phone    string
	Email    string
	Password string
}

func CreateCustomer(name string, phone string, email string, password string) (int64, error) {
	schema, err := db.Connect()
	if err != nil {
		return 0, err
	}
	defer schema.Close()
	var customerID int64
	err = schema.QueryRow("INSERT INTO customers (name, phone, email, password) VALUES ($1, $2,$3,$4 ) RETURNING id", name, phone, email, password).Scan(&customerID)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	return customerID, nil
}
func ReadCustomers() ([]Customers, error) {
	schema, err := db.Connect()
	if err != nil {
		return []Customers{}, err
	}
	defer schema.Close()
	rows, err := schema.Query("SELECT id, name, email, phone,password FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customers
	for rows.Next() {
		var customer Customers
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Password); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
func ReadCustomer(customerID int64) (string, string, string, error) {
	schema, err := db.Connect()
	if err != nil {
		return "", "", "", err
	}
	defer schema.Close()
	var name, phone, email string
	err = schema.QueryRow("SELECT name, phone, email FROM customers WHERE id = $1", customerID).Scan(&name, &phone, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", fmt.Errorf("customer with ID %d not found", customerID)
		}
		return "", "", "", err
	}

	return name, email, phone, nil
}
func ReadCustomerByEmail(email string) (int, string, string, string, error) {
	schema, err := db.Connect()
	if err != nil {
		return 0, "", "", "", nil
	}
	defer schema.Close()

	var customer Customers
	err = schema.QueryRow("SELECT id,phone,email,password FROM customers WHERE email = $1", email).Scan(&customer.ID, &customer.Phone, &customer.Email, &customer.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", "", "", fmt.Errorf("customer with email %s not found", email)
		}
		return 0, "", "", "", nil
	}

	return customer.ID, customer.Email, customer.Phone, customer.Password, nil
}
