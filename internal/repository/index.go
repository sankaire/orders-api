package repository

import (
	"database/sql"
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/db"
	time2 "time"
)

type Customers struct {
	ID       int
	Name     string
	Phone    string
	Email    string
	Password string
}
type Orders struct {
	ID         int64
	Item       string
	Amount     int64
	CustomerID any
	Time       time2.Time
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
func CreateOrder(customerID any, item string, amount int) (int64, error) {
	schema, err := db.Connect()
	if err != nil {
		return 0, err
	}
	defer schema.Close()

	var orders Orders
	time := time2.Now()
	err = schema.QueryRow("INSERT INTO orders (customer_id,item,amount,time) VALUES ($1, $2,$3,$4 ) RETURNING id ", customerID, item, amount, time).Scan(&orders.ID)
	if err != nil {
		return 0, err
	}
	return orders.ID, nil
}
func ReadCustomerOrder(orderID int64) (int64, any, string, int64, time2.Time, error) {
	schema, err := db.Connect()
	if err != nil {
		return 0, 0, "", 0, time2.Time{}, nil
	}
	defer schema.Close()

	var orders Orders
	err = schema.QueryRow("SELECT id,customer_id,item,amount,time FROM orders WHERE id = $1", orderID).Scan(&orders.ID, &orders.CustomerID, &orders.Item, &orders.Amount, &orders.Time)
	if err != nil {
		return 0, 0, "", 0, time2.Time{}, err
	}

	return orders.ID, orders.CustomerID, orders.Item, orders.Amount, orders.Time, nil
}
func ReadAllCustomerOrder(customerID any) ([]Orders, error) {
	schema, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer schema.Close()

	var orders []Orders
	rows, err := schema.Query("SELECT id,customer_id item,amount,time FROM orders WHERE customer_id = $1", customerID)
	for rows.Next() {
		var order Orders
		if err := rows.Scan(&order.ID, &order.Item, &order.Amount, &order.Time); err != nil {
			return nil, err
		}
		order.CustomerID = customerID
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
