package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connectionString := os.Getenv("DATABASE_URL")

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping failed: ", err)
	}

	log.Println("Connected to database successfully")
}