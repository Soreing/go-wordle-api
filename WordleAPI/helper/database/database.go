package database

import (
	"os"
	"fmt"
	"database/sql"
	"github.com/lib/pq"
)

// Global Reference to Database
var db *sql.DB

// Connects to the Wordle database with credentials
func Connect() {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPASS"),
		os.Getenv("DBHOST"),
		os.Getenv("DBSSL"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	
	if err != nil {
		panic("Error: Database connection failed")
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		panic("Error: Database connection failed")
	}
}

// Disconnects from the Wordle database
func Disconnect(){
	db.Close()
}

// Returns the database instance
func Pool() *sql.DB {
	return db
}

// Gets the Postgres Error Code from the error
func GetErrorCode(err error) string {
	pqerr := err.(*pq.Error)
	return string(pqerr.Code)
}