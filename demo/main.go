package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

type QR struct {
	Num string
}

func postQRHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("public/mobile.html")
	if err != nil {
		log.Fatalln(err)
	}

	qrNum := r.FormValue("qr_num")
	fmt.Println("POST " + qrNum) //qr_num

	tmp.Execute(w, QR{qrNum})

	//http.ServeFile(w, r, "public/mobile.html")
	//http.Redirect(w, r, "/mobile.html", http.StatusMovedPermanently) //Get, 301
	//http.Redirect(w, r, "/mobile.html", http.StatusTemporaryRedirect) //Post, 307
}

func main() {
	mux := pat.New()
	mux.Post("/mobile", postQRHandler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
