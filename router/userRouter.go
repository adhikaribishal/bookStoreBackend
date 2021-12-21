package router

import (
	"net/http"
	"strconv"

	"github.com/adhikaribishal/bookStoreBackend/controllers"
	"github.com/adhikaribishal/bookStoreBackend/helpers"
)

type user struct {
	id int
}

func ServeUsers(w http.ResponseWriter, r *http.Request) {
	var head string

	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)

	switch head {
	case "":
		if r.Method == "GET" {
			controllers.GetAllUsers(w, r)
		} else {
			controllers.CreateUser(w, r)
		}
	case "authenticate":
		controllers.Authenticate(w, r)
	default:
		id, err := strconv.Atoi(head)
		if err != nil || id <= 0 {
			http.NotFound(w, r)
			return
		}
		user{id}.serveHTTP(w, r)
	}
}

func (h user) serveHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)

	switch head {
	case "":
		if r.Method == "GET" {
			controllers.GetUser(w, r, h.id)
		} else if r.Method == "PATCH" {
			controllers.UpdateUser(w, r, h.id)
		} else {
			controllers.DeleteUser(w, r, h.id)
		}
	default:
		http.NotFound(w, r)
	}
}
