package main

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

func main() {
	mux := pat.New()
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", leftUserHandler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
