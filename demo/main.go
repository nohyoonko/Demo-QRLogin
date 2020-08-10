package main

import (
	"demo/app"
	"log"
)

func main() {
	r := app.MakeHandler("dbPath")
	defer r.Close()

	log.Println("Started App")
	r.Run()
}
