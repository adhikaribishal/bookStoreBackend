package router

import (
	"fmt"
	"net/http"

	"github.com/adhikaribishal/bookStoreBackend/helpers"
)

var Serve = helpers.NoTrailingSlash(serve)

func serve(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
	switch head {
	case "":
		serveHome(w, r)
	case "api":
		serveApi(w, r)
	default:
		http.NotFound(w, r)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home\n")
}

func serveApi(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
	switch head {
	case "users":
		ServeUsers(w, r)
	default:
		http.NotFound(w, r)
	}
}
