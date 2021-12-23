package models

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/adhikaribishal/bookStoreBackend/database"
	"github.com/adhikaribishal/bookStoreBackend/helpers"
)

type Book struct {
	ID            int64  `json:"id" sql:"id"`
	Title         string `json:"title" sql:"title"`
	Author        string `json:"author" sql:"author"`
	Price         string `json:"price" sql:"price"`
	Publication   string `json:"publication" sql:"publication"`
	PublishedDate string `json:"published_date" sql:"published_date"`
	ISBN          string `json:"isbn" sql:"isbn"`
}

func (book *Book) Validate() (map[string]interface{}, bool) {
	if book.Title == "" {
		return helpers.Message(false, "Book title should be on the payload"), false
	}

	if book.Author == "" {
		return helpers.Message(false, "Book author should be on the payload"), false
	}

	if book.Price == "" {
		return helpers.Message(false, "Price should be on the payload"), false

	}

	if book.Publication == "" {
		return helpers.Message(false, "Book publication should be on the payload"), false

	}

	if book.PublishedDate == "" {
		return helpers.Message(false, "Published date should be on the payload"), false

	}

	if book.ISBN == "" {
		return helpers.Message(false, "ISBN should be on the payload"), false

	}

	temp := &Book{}

	db := database.CreateDatabseConnection()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS books(
		id SERIAL PRIMARY KEY,
		title text NOT NULL,
		author text NOT NULL,
		price text NOT NULL,
		publication text NOT NULL,
		published_date text NOT NULL,
		isbn text NOT NULL,
		UNIQUE(isbn)
	)`)
	if err != nil {
		log.Fatalf("Unable to create the books table. %v", err)
	}

	sqlStatement := `SELECT isbn FROM books WHERE isbn=$1`

	row := db.QueryRow(sqlStatement, book.ISBN)

	err = row.Scan(&temp.ISBN)
	if err != nil && err != sql.ErrNoRows {
		return helpers.Message(false, "Connection error. Please retry"), false
	}

	if temp.ISBN != "" {
		return helpers.Message(false, "Book ISBN should be unique."), false
	}

	return helpers.Message(true, "success"), true

}

func (book *Book) Create() (map[string]interface{}, int) {
	if resp, ok := book.Validate(); !ok {
		return resp, http.StatusBadRequest
	}

	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	sqlStatement := `INSERT INTO books
					(title, author, price, publication, published_date, isbn)
					VALUES($1, $2, $3, $4, $5, $6)
					RETURNING id`

	var id int64
	err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Price, book.Publication, book.PublishedDate, book.ISBN).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	if id <= 0 {
		return helpers.Message(false, "Failed to create book, connection error."), http.StatusInternalServerError
	}

	book.ID = id
	response := helpers.Message(true, "Book has been created")
	response["book"] = book
	return response, http.StatusCreated
}
