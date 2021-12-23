package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/adhikaribishal/bookStoreBackend/database"
	"github.com/adhikaribishal/bookStoreBackend/helpers"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type key int64
type Token struct {
	UserId key
	jwt.StandardClaims
}

type User struct {
	ID        int64     `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Username  string    `json:"username" sql:"username"`
	FirstName string    `json:"first_name" sql:"first_name"`
	LastName  string    `json:"last_name" sql:"last_name"`
	CreatedAt time.Time `json:"createdat" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedat" sql:"updated_at"`
	Token     string    `json:"token" sql:"-"`
}

type UserRetrieve struct {
	ID        string `json:"id" sql:"id"`
	Email     string `json:"email" validate:"required" sql:"email"`
	Username  string `json:"username" sql:"username"`
	FirstName string `json:"first_name" sql:"first_name"`
	LastName  string `json:"last_name" sql:"last_name"`
}

func (user *User) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		return helpers.Message(false, "Email Address is required"), false
	}

	if len(user.Password) < 6 {
		return helpers.Message(false, "Password is required"), false
	}

	temp := &User{}

	db := database.CreateDatabseConnection()
	defer db.Close()

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

	sqlStatement := `SELECT email FROM users WHERE email=$1`

	row := db.QueryRow(sqlStatement, user.Email)

	err = row.Scan(&temp.Email)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Error: %v", err)
		return helpers.Message(false, "Connection error. Please retry"), false
	}
	fmt.Println("ASD")

	if temp.Email != "" {
		return helpers.Message(false, "Email address already in use by another user."), false
	}

	return helpers.Message(false, "Requirement passed"), true
}

func (user *User) Create() map[string]interface{} {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	db := database.CreateDatabseConnection()
	defer db.Close()

	log.Println("Succesfully connected!")

	sqlStatement := `INSERT INTO users
					(email, password, username, first_name, last_name)
					VALUES ($1,$2,$3,$4,$5)
					RETURNING id`

	var id int64
	err := db.QueryRow(sqlStatement, user.Email, user.Password, user.Username, user.FirstName, user.LastName).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	if id <= 0 {
		return helpers.Message(false, "Failed to create account, connection error.")
	}

	tk := &Token{UserId: key(id)}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	user.Token = tokenString

	user.Password = ""
	user.ID = id

	response := helpers.Message(true, "Account has been created")
	response["user"] = user
	return response
}

func Login(email, password string) map[string]interface{} {
	user := &User{}

	db := database.CreateDatabseConnection()
	defer db.Close()

	sqlStatement := `SELECT id, email, username, first_name, last_name FROM users WHERE email=$1`

	row := db.QueryRow(sqlStatement, email)

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName)

	if err != nil {
		if err == sql.ErrNoRows {
			return helpers.Message(false, "Email address not found")
		}
		return helpers.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return helpers.Message(false, "Invalid login credentials. Please try again")
	}

	user.Password = ""

	tk := &Token{UserId: key(user.ID)}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	user.Token = tokenString

	resp := helpers.Message(true, "Logged In")
	resp["user"] = user
	return resp
}
