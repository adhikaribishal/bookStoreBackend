package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func CreateDatabseConnection() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Could not connect to database %s at host %s:%s", dbName, host, port)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
