package repository

import (
	"database/sql"
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/db"
)

type Customers struct {
	ID   int
	Name string
	Code string
}

func CreateCustomer(name string, code string) (int64, error) {
	schema, err := db.Connect()
	if err != nil {
		return 0, err
	}
	defer schema.Close()
	var customerID int64
	err = schema.QueryRow("INSERT INTO customers (name, code) VALUES ($1, $2) RETURNING id", name, code).Scan(&customerID)
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
	rows, err := schema.Query("SELECT id, name, code FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customers
	for rows.Next() {
		var customer Customers
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Code); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
func ReadCustomer(customerID int64) (string, string, error) {
	schema, err := db.Connect()
	if err != nil {
		return "", "", err
	}
	defer schema.Close()
	var name, code string
	err = schema.QueryRow("SELECT name, code FROM customers WHERE id = $1", customerID).Scan(&name, &code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("customer with ID %d not found", customerID)
		}
		return "", "", err
	}

	return name, code, nil
}
