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
	err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//
	//}
	bdUri := os.Getenv("DB_URI")
	var db *sql.DB
	//var err error
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, user, password, dbname)
	db, err = sql.Open("postgres", bdUri)
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
