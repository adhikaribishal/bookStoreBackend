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

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "POST") {
		return
	}

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		helpers.Respond(w, helpers.Message(false, "Invalid Request"), http.StatusBadRequest)
		return
	}

	resp := models.Login(user.Email, user.Password)
	helpers.Respond(w, resp, http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "POST") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		helpers.Respond(w, helpers.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp := user.Create()
	helpers.Respond(w, resp, http.StatusCreated)
}

func GetUser(w http.ResponseWriter, r *http.Request, userID int) {
	if !helpers.EnsureMethod(w, r, "GET") {
		return
	}

	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Orgin", "*")

	response, statusCode := getUser(userID)

	helpers.Respond(w, response, statusCode)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "GET") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response, statusCode := getAllUsers()
	helpers.Respond(w, response, statusCode)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, userID int) {
	if !helpers.EnsureMethod(w, r, "PATCH") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unble to decode the request body. %v", err)
	}

	response, statusCode := updateUser(int64(userID), user)
	helpers.Respond(w, response, statusCode)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	if !helpers.EnsureMethod(w, r, "DELETE") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	response, statusCode := deleteUser(int64(userID))
	helpers.Respond(w, response, statusCode)
}

func getUser(userID int) (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	var user models.UserRetrieve

	sqlStatement := `SELECT id, email, username, first_name, last_name FROM users WHERE id=$1`

	row := db.QueryRow(sqlStatement, userID)

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName)
	if err != nil && err != sql.ErrNoRows {
		return helpers.Message(false, "Connection error. Please retry"), http.StatusInternalServerError
	}

	response := helpers.Message(true, "User fetched succesfully")
	response["user"] = user
	return response, http.StatusOK
}

func getAllUsers() (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Successfully connected!")

	var users []models.UserRetrieve

	sqlStatement := `SELECT id, email, username, first_name, last_name FROM users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return helpers.Message(false, "Unable to execute the query."), http.StatusBadRequest
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserRetrieve

		err = rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName)

		if err != nil {
			return helpers.Message(false, "Unable to scan the row."), http.StatusBadRequest
		}

		users = append(users, user)
	}

	response := helpers.Message(true, "Users fetched successfully")
	response["users"] = users
	return response, http.StatusOK
}

func updateUser(id int64, user models.User) (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `UPDATE users SET email=$2, first_name=$3, last_name=$4 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return helpers.Message(false, "Unable to execute the query."), http.StatusInternalServerError
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helpers.Message(false, "Error while checking the affeced rows."), http.StatusInternalServerError
	}

	fmt.Printf("Total rows/record affected %v\n", rowsAffected)

	response := helpers.Message(true, "User updated successfully")
	response["user"] = user

	return response, http.StatusOK
}

func deleteUser(id int64) (map[string]interface{}, int) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM users WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		return helpers.Message(false, "Unable to execute the query."), http.StatusInternalServerError
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helpers.Message(false, fmt.Sprintf("Error while checking the affected rows. %v", err)), http.StatusInternalServerError
	}

	response := helpers.Message(true, fmt.Sprintf("User deleted successfully. Rows affected: %v", rowsAffected))
	return response, http.StatusOK
}
