package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sankaire/orders-api/internal/db/internal/db"
	"github.com/sankaire/orders-api/internal/db/internal/handlers"
	"github.com/sankaire/orders-api/internal/db/internal/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	db.CreateTables()
	http.HandleFunc("/api/customer", handlers.CreateCustomerHandler)
	http.HandleFunc("/api/customer/login", handlers.LoginCustomerHandler)
	http.Handle("/api/order", middleware.Authenticate(http.HandlerFunc(handlers.CreateCustomerOrder)))
	http.Handle("/api/orders", middleware.Authenticate(http.HandlerFunc(handlers.ReadAllCustomerOrdersHandler)))
	http.Handle("/api/order/", middleware.Authenticate(http.HandlerFunc(handlers.ReadCustomerOrder)))
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
