package main

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	//err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomNumber(n int) (string, error) {
	const numbers = "0123456789"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = numbers[b%byte(len(numbers))]
	}
	return string(bytes), nil
}

func main() {
	r := gin.Default()
	r.Static("/js", "./js")
	r.Static("/css", "./css")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.POST("/create", func(c *gin.Context) {
		randNum, err := generateRandomNumber(7)
		if err != nil {
			log.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"randNum": randNum,
		})
	})
	r.Run()
}
