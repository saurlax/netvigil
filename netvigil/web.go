package netvigil

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var secret = make([]byte, 32)

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == viper.GetString("username") && password == viper.GetString("password") {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
			"username": username,
		})
		tokenString, err := token.SignedString(secret)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": tokenString})
	} else {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
	}
}

func recordsHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 0
	}
	records, err := GetRecords(1000, page)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error()})
	}
	c.JSON(200, records)
}

func init() {
	rand.Read(secret)
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
			path = "dist" + path
			_, err := os.Stat(path)
			if err == nil {
				ctx.File(path)
			} else {
				ctx.File("dist/index.html")
			}
		}
	})
	fmt.Printf("Web server started on http://%s/\n", addr)
	go r.Run(addr)
}
