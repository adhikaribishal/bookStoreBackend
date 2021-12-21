package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adhikaribishal/bookStoreBackend/helpers"
	"github.com/adhikaribishal/bookStoreBackend/models"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "POST") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
