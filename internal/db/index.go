package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func Connect() (*sql.DB, error) {
	godotenv.Load()
	dbURI := os.Getenv("DB_URI")
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database")
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
