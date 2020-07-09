package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

func postQRHandler(w http.ResponseWriter, r *http.Request) {
	qrcodeNumber := r.FormValue("qrnum")
	fmt.Fprintf(w, "qrcode: %s\n", qrcodeNumber)
}

func main() {
	mux := pat.New()
	mux.Post("/qrcode", postQRHandler)
	//mux.Get("/success", getQRHandler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
