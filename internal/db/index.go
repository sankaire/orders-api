package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Connect() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "0000"
		dbname   = "db1"
	)
	var db *sql.DB
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to db")
	return db, nil
}
func CreateTables() {
	db, err := Connect()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		phone TEXT,
		email TEXT,
		password TEXT
	)
`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Customers table created successfully")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		customer_id INTEGER,
		item TEXT,
		amount INTEGER,
		time TIMESTAMP,
		FOREIGN KEY (customer_id) REFERENCES customers(id)
	)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Orders table created successfully")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
}
