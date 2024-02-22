package main

import (
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/db"
	"github.com/sankaire/orders-api/internal/db/internal/handlers"
	"github.com/sankaire/orders-api/internal/db/internal/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	db.CreateTables()
	port := os.Getenv("PORT")
	http.HandleFunc("/api/customer", handlers.CreateCustomerHandler)
	http.Handle("/api/customers", middleware.Authenticate(http.HandlerFunc(handlers.ReadCustomersHandler)))
	http.HandleFunc("/api/customer/login", handlers.LoginCustomerHandler)
	http.Handle("/api/order", middleware.Authenticate(http.HandlerFunc(handlers.CreateCustomerOrder)))
	http.Handle("/api/orders", middleware.Authenticate(http.HandlerFunc(handlers.ReadAllCustomerOrdersHandler)))
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
