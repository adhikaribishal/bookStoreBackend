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

	user, err := getUser(userID)
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(users)
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

	updatedRows := updateUser(int64(userID), user)

	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	res := response{
		ID:      int64(userID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	if !helpers.EnsureMethod(w, r, "DELETE") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	deletedRows := deleteUser(int64(userID))

	msg := fmt.Sprintf("User deleted successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(userID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertUser(user models.User) int64 {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(
						id SERIAL PRIMARY KEY,
						email text NOT NULL,
						password text NOT NULL,
						username text,
						first_name text,
						last_name text,
						created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
						UNIQUE(email, username)
						)`)
	if err != nil {
		log.Fatal(err)
	}

	sqlStatement := `INSERT INTO users
					(email, password, username, first_name, last_name)
					VALUES ($1,$2,$3,$4,$5)
					RETURNING id`

	var id int64
	err = db.QueryRow(sqlStatement, user.Email, user.Password, user.Username, user.FirstName, user.LastName).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	return id
}

func getUser(userID int) (models.UserRetrieve, error) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	var user models.UserRetrieve

	sqlStatement := `SELECT id, email, username, first_name, last_name FROM users WHERE id=$1`

	row := db.QueryRow(sqlStatement, userID)

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName)

	switch err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return user, err
}

func getAllUsers() ([]models.UserRetrieve, error) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Successfully connected!")

	var users []models.UserRetrieve

	sqlStatement := `SELECT id, email, username, first_name, last_name FROM users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserRetrieve

		err = rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)
	}

	return users, err
}

func updateUser(id int64, user models.User) int64 {
	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `UPDATE users SET email=$2, first_name=$3, last_name=$4 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, user.Email, user.FirstName, user.LastName)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affeced rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v\n", rowsAffected)

	return rowsAffected
}

func deleteUser(id int64) int64 {
	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM users WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)

	return rowsAffected
}
