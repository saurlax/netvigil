package netvigil

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var signKey = make([]byte, 32)

func recordsHandler(c *gin.Context) {
	var err error
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 0
	}
	records, err := GetSortedRecords(c.Param("key"), 100, page)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error()})
	}
	c.JSON(200, records)
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == viper.GetString("username") && password == viper.GetString("password") {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["expire"] = time.Now().Add(time.Hour * 72).Unix()

		tokenString, err := token.SignedString(signKey)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": tokenString})
	} else {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
	}
}

func init() {
	rand.Read(signKey)
	addr := viper.GetString("web")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		switch path {
		case "/api/records":
			recordsHandler(ctx)
		case "/api/login":
			loginHandler(ctx)
		default:
			path = "static" + path
			_, err := os.Stat(path)
			if err == nil {
				ctx.File(path)
			} else {
				ctx.File("static/index.html")
			}
		}
	})
	fmt.Printf("Web server started on http://%s/\n", addr)
	go r.Run(addr)
}
