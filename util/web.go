package util

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	r            *gin.Engine
	mySigningKey = []byte("sss")
)

func recordsHandler(c *gin.Context) {
	var err error
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Invalid page number"})
	}
	records, err := GetSortedRecords(c.Param("key"), 100, page)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error()})
	}
	c.JSON(200, gin.H{
		"records": records,
	})
}

func loginHandler(c *gin.Context) {
	Username := c.PostForm("username")
	Password := c.PostForm("password")

	if Username == "admin" && Password == "123456" {
		// c.Redirect(http.StatusMovedPermanently, "/")
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = Username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully Login!", "token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password", "username": Username, "password": Password})
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()
	r.Use(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		switch path {
		case "/api/records":
			recordsHandler(ctx)
		case "/api/login":
			loginHandler(ctx)
		default:
			path = "static" + path
			println(path)
			_, err := os.Stat(path)
			if err == nil {
				ctx.File(path)
			} else {
				ctx.File("static/index.html")
			}
		}
	})
}

func Run() {
	addr := viper.GetString("web")
	fmt.Printf("Web server started on http://%s/\n", addr)
	r.Run(addr)
}
