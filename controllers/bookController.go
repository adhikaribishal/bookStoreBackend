package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adhikaribishal/bookStoreBackend/database"
	"github.com/adhikaribishal/bookStoreBackend/helpers"
	"github.com/adhikaribishal/bookStoreBackend/models"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "POST") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	book := &models.Book{}

	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		helpers.Respond(w, helpers.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp, statusCode := book.Create()
	helpers.Respond(w, resp, statusCode)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "GET") {
		return
	}

	response, statusCode := getAllBooks()
	helpers.Respond(w, response, statusCode)
}

func GetBook(w http.ResponseWriter, r *http.Request, bookID int64) {
	if !helpers.EnsureMethod(w, r, "GET") {
		return
	}

	response, statusCode := getBook(bookID)
	helpers.Respond(w, response, statusCode)
}

func UpdateBook(w http.ResponseWriter, r *http.Request, bookID int64) {
	if !helpers.EnsureMethod(w, r, "PATCH") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		helpers.Respond(w, helpers.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	response, statusCode := updateBook(bookID, book)
	helpers.Respond(w, response, statusCode)
}

func DeleteBook(w http.ResponseWriter, r *http.Request, bookID int64) {
	if !helpers.EnsureMethod(w, r, "DELETE") {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	response, statusCode := deleteBook(bookID)
	helpers.Respond(w, response, statusCode)
}

func getAllBooks() (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Successfully connected!")

	var books []models.Book

	sqlStatement := `SELECT * FROM books`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return helpers.Message(false, "Unable to execute the query."), http.StatusBadRequest
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book

		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Publication, &book.PublishedDate, &book.ISBN)

		if err != nil {
			return helpers.Message(false, "Unable to scan the row."), http.StatusBadRequest
		}

		books = append(books, book)
	}

	response := helpers.Message(true, "Books fetched successfully")
	response["books"] = books
	return response, http.StatusOK
}

func getBook(bookID int64) (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	var book models.Book

	sqlStatement := `SELECT * FROM books WHERE id=$1`

	row := db.QueryRow(sqlStatement, bookID)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Publication, &book.PublishedDate, &book.ISBN)
	if err != nil && err != sql.ErrNoRows {
		return helpers.Message(false, "Connection error. Please retry"), http.StatusInternalServerError
	}

	response := helpers.Message(true, "Book fetched successfully")
	response["book"] = book
	return response, http.StatusOK
}

func updateBook(id int64, book models.Book) (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `UPDATE books SET title=$2, author=$3, price=$4, publication=$5, published_date=$6, isbn=$7 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, book.Title, book.Author, book.Price, book.Publication, book.PublishedDate, book.ISBN)
	if err != nil {
		return helpers.Message(false, "Unable to execute the query."), http.StatusInternalServerError
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helpers.Message(false, "Error while checking the affeced rows."), http.StatusInternalServerError
	}

	fmt.Printf("Total rows/record affected %v\n", rowsAffected)

	book.ID = id
	response := helpers.Message(true, "Book updated successfully")
	response["book"] = book

	return response, http.StatusOK
}

func deleteBook(id int64) (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM books WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		return helpers.Message(false, "Unable to execute the query."), http.StatusInternalServerError
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helpers.Message(false, fmt.Sprintf("Error while checking the affected rows. %v", err)), http.StatusInternalServerError
	}

	response := helpers.Message(true, fmt.Sprintf("Book deleted successfully. Rows affected: %v", rowsAffected))
	return response, http.StatusOK
}
