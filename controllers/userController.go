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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, "POST") {
		return
	}

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertId := insertUser(user)

	res := response{
		ID:      insertId,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetUser(w http.ResponseWriter, r *http.Request, userID int) {
	if !helpers.EnsureMethod(w, r, "GET") {
		return
	}

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Orgin", "*")

	user, err := getUser(userID)
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
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

	fmt.Fprintf(w, "Update User: %d\n", userID)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	if !helpers.EnsureMethod(w, r, "DELETE") {
		return
	}

	fmt.Fprintf(w, "Delete User: %d\n", userID)
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
						updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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

func getUser(userID int) (models.User, error) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	var user models.User

	sqlStatement := `SELECT * FROM users WHERE id=$1`

	row := db.QueryRow(sqlStatement, userID)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt)

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

func getAllUsers() ([]models.User, error) {
	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Successfully connected!")

	var users []models.User

	sqlStatement := `SELECT * FROM users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.FirstName, &user.LastName,
			&user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)
	}

	return users, err
}
