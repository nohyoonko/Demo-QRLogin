package main

import (
	"demo/app"
	"log"
)

func main() {
	const env bool = true //file 삭제 권한 여부를 제어하는 변수, true면 삭제 가능

	r := app.NewQRHandler("dbPath", env)
	defer r.Close()

	log.Println("Started App")
	r.Run()
}
