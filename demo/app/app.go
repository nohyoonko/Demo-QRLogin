package app

import (
	"crypto/rand"
	"demo/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QRHandler interface {
	homeHandler(c *gin.Context)     //메인 페이지 로드하기
	createQRHandler(c *gin.Context) //QR 코드에 담겨지는 random number 생성하기
	deleteQRHandler(c *gin.Context) //random number 삭제하기(메인 페이지에서 취소 버튼을 누르면 요청)
	verifyHandler(c *gin.Context)   //random number와 credential을 받아서 유효한지 확인하기
	checkHandler(c *gin.Context)    //random number에 대한 유효성 검사한 결과 확인하기
	successHandler(c *gin.Context)  //성공 페이지 로드하기
	Run()
	Close()
}

type appHandler struct {
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

func NewQRHandler(filepath string, env bool) QRHandler {
	return newAppHandler(filepath, env)
}

func newAppHandler(filepath string, env bool) QRHandler {
	r := gin.Default()
	r.Static("/js", "./js")
	r.Static("/css", "./css")
	r.LoadHTMLGlob("templates/*")

	a := &appHandler{
		gin: r,
		db:  model.NewDBHandler(filepath, env), //DB open
	}

	r.GET("/", a.homeHandler)
	r.POST("/create", a.createQRHandler)
	r.DELETE("/create", a.deleteQRHandler)
	r.POST("/mobile", a.verifyHandler)
	r.POST("/check", a.checkHandler)
	r.GET("/success", a.successHandler)

	return a
}

/* Generate Random Number (length: 7) */
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

func (a *appHandler) Run() {
	a.gin.Run()
}

func (a *appHandler) Close() {
	a.db.Close()
}

func (a *appHandler) homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (a *appHandler) createQRHandler(c *gin.Context) {
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

func (a *appHandler) deleteQRHandler(c *gin.Context) {
	var q QRCode
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.db.Delete(q.QRNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"delete": "fail"})
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})
	log.Println("Delete Random Number for QR Code")
}

func (a *appHandler) verifyHandler(c *gin.Context) {
	var v Verify
	if err := c.ShouldBindJSON(&v); err != nil {
		//if credential is empty string, response error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Credential is empty string")
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

func (a *appHandler) checkHandler(c *gin.Context) {
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

func (a *appHandler) successHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", gin.H{})
}
