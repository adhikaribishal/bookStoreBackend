package router

import (
	"net/http"
	"strconv"

	"github.com/adhikaribishal/bookStoreBackend/controllers"
	"github.com/adhikaribishal/bookStoreBackend/helpers"
)

type book struct {
	ID int
}

func ServeBooks(w http.ResponseWriter, r *http.Request) {
	var head string

	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)

	switch head {
	// case "":
	// 	controllers.GetAllBooks(w, r)
	case "create":
		controllers.CreateBook(w, r)
	default:
		id, err := strconv.Atoi(head)
		if err != nil || id <= 0 {
			http.NotFound(w, r)
			return
		}
		user{id}.serveHTTP(w, r)
	}
}
