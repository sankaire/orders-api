package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Connect() (*sql.DB, error) {
	//err := godotenv.Load()
	//if err != nil {
	//	panic("Error loading .env file")
	//}
	//dbURI := os.Getenv("DB_URI")
	//if dbURI == "" {
	//	panic("DB_URI environment variable not set")
	//}
	//
	//fmt.Println("DB_URI:", "dbURI")
	db, err := sql.Open("postgres", "postgres://gxvhqpskrrabic:18df20bd44f4ee663808faadc5395258e2ce539d82aad6da8bfe4cfbfd275644@ec2-54-78-142-10.eu-west-1.compute.amazonaws.com:5432/d4ppro9ijibn6c")
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
