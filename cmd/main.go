package main

import (
	"fmt"
	"github.com/sankaire/orders-api/internal/db/internal/db"
	"github.com/sankaire/orders-api/internal/db/internal/handlers"
	"log"
	"net/http"
)

func main() {
	db.CreateTables()
	http.HandleFunc("/api/customer", handlers.CreateCustomerHandler)
	http.HandleFunc("/api/customers", handlers.ReadCustomersHandler)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
