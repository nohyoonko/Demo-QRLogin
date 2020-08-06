package app

import (
	"crypto/rand"
	"demo/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	gin *gin.Engine
	db  model.DBHandler
}

type Verify struct {
	RandomNumber string `json:"randNum" binding:"required"`
	Credential   string `json:"credential" binding:"required"`
}

type QRCode struct {
	QRNumber string `json:"randNum" binding:"required"`
}

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

func (a *AppHandler) Run() {
	a.gin.Run()
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func (a *AppHandler) homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (a *AppHandler) createQRHandler(c *gin.Context) {
	for {
		randNum, err := generateRandomNumber(7)
		if err != nil {
			log.Println(err.Error())
		}

		//if there is no error(no duplication), loop finish
		if a.db.Create(randNum, "wait") == nil {
			c.JSON(http.StatusOK, gin.H{"randNum": randNum})
			break
		}
	}
	log.Println("Create Random Number for QR Code")
}

func (a *AppHandler) verifyHandler(c *gin.Context) {
	var v Verify
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//verify credential -> only "test" is ok
	if v.Credential != "test" {
		err := a.db.Update(v.RandomNumber, "fail")
		if err != nil {
			log.Println(err.Error())
		}
		c.JSON(http.StatusUnauthorized, gin.H{"verify": "fail"})
		log.Println("Verification Fail")
	} else {
		err := a.db.Update(v.RandomNumber, "success")
		if err != nil {
			log.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{"verify": "success"})
		log.Println("Verification Success")
	}
}

func (a *AppHandler) checkHandler(c *gin.Context) {
	var qr QRCode
	if err := c.ShouldBindJSON(&qr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//read QRNumber's status
	data, err := a.db.SelectOne(qr.QRNumber)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"QRstatus": data.Value}) //key(randNum), value(status)
	log.Println("Check Random Number's Status")
	err = a.db.Delete(qr.QRNumber)
	if err != nil {
		log.Println(err.Error())
	}
}

func (a *AppHandler) successHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", gin.H{})
}

func MakeHandler(filepath string) *AppHandler {
	r := gin.Default()
	r.Static("/js", "./js")
	r.Static("/css", "./css")
	r.LoadHTMLGlob("templates/*")

	a := &AppHandler{
		gin: r,
		db:  model.NewDBHandler(filepath), //DB open
	}

	r.GET("/", a.homeHandler)
	r.POST("/create", a.createQRHandler)
	r.POST("/mobile", a.verifyHandler)
	r.POST("/check", a.checkHandler)
	r.GET("/success", a.successHandler)

	return a
}
