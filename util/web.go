package util

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var secret = make([]byte, 32)

func authHandler(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Authorization header not found"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Claims.(jwt.MapClaims)["exp"].(float64) < float64(time.Now().Unix()) {
			return nil, fmt.Errorf("token is expired")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("username", token.Claims.(jwt.MapClaims)["sub"])
	c.Next()
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == viper.GetString("username") && password == viper.GetString("password") {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"sub": username,
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
	records, err := GetRecords(200, page)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, records)
}

func readConfigHandler(c *gin.Context) {
	c.JSON(200, viper.AllSettings())
}

func writeConfigHandler(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if oldPassword, ok := updates["oldPassword"]; ok {
		if newPassword, ok := updates["newPassword"]; ok {
			if oldPassword != viper.GetString("password") {
				c.JSON(401, gin.H{"error": "旧密码不正确"})
				return
			}
			viper.Set("password", newPassword)
		}
		delete(updates, "oldPassword")
		delete(updates, "newPassword")
		delete(updates, "blacklist")
	}

	for key, value := range updates {
		viper.Set(key, value)
	}

	if err := viper.WriteConfig(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "ok")
}

func staticHandler(c *gin.Context) {
	path := c.Request.URL.Path
	if !strings.HasPrefix(path, "/api") {
		path = "dist" + path
		_, err := os.Stat(path)
		if err == nil {
			c.File(path)
		} else {
			c.File("dist/index.html")
		}
	}
}

func checkHandler(c *gin.Context) {
	var request struct {
		Token string   `json:"token"`
		IPs   []string `json:"ips"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	records, err := GetRecordsByIPs(request.IPs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, records)
}

func init() {
	rand.Read(secret)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/api/login", loginHandler)
	r.GET("/api/records", authHandler, recordsHandler)
	r.GET("/api/config", authHandler, readConfigHandler)
	r.POST("/api/config", authHandler, writeConfigHandler)
	r.POST("/api/check", checkHandler)
	r.NoRoute(staticHandler)

	addr := viper.GetString("web")
	fmt.Printf("Web server started on http://%s/\n", addr)
	go r.Run(addr)
}
