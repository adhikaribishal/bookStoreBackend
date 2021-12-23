package router

import (
	"net/http"
	"strconv"

	"github.com/adhikaribishal/bookStoreBackend/controllers"
	"github.com/adhikaribishal/bookStoreBackend/helpers"
)

type book struct {
	id int64
}

func ServeBooks(w http.ResponseWriter, r *http.Request) {
	var head string

	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)

	switch head {
	case "":
		controllers.GetAllBooks(w, r)
	case "create":
		controllers.CreateBook(w, r)
	default:
		id, err := strconv.Atoi(head)
		if err != nil || id <= 0 {
			http.NotFound(w, r)
			return
		}
		book{int64(id)}.serveHTTP(w, r)
	}
}

func (h book) serveHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)

	switch head {
	case "":
		if r.Method == "GET" {
			controllers.GetBook(w, r, h.id)
		} else if r.Method == "PATCH" {
			controllers.UpdateBook(w, r, h.id)
		} else {
			controllers.DeleteBook(w, r, h.id)
		}
	default:
		http.NotFound(w, r)
	}
}
