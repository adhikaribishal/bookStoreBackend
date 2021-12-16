package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adhikaribishal/bookStoreBackend/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("BACKEND_PORT")
	fmt.Println("Port: ", port)

	router := http.HandlerFunc(router.Serve)
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
