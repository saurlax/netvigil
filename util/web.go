package util

import (
	"crypto/rand"
	"fmt"
	"log"
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
	c.Next()
}

func loginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Username == viper.GetString("username") && req.Password == viper.GetString("password") {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"sub": req.Username,
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

func netstatsHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid limit number"})
		return
	}
	// Limit the number of records to prevent abuse
	if limit > 500 {
		c.JSON(400, gin.H{"error": "Limit number too large"})
		return
	}
	netstats, total, err := GetNetstats(limit, page)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{
		"netstats": netstats,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func threatsHandler(c *gin.Context) {
	records, err := GetThreats()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, records)
}

func deleteThreatsHandler(c *gin.Context) {
	ip := c.Param("ip")

	if err := DeleteThreatsByIP(ip); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	DeleteFireWallRule(ip)
}

func clientsHandler(c *gin.Context) {
	clients := GetClients()
	c.JSON(200, clients)
}

func createClientHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := CreateClient(req.Name); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "ok")
}

func deleteClientHandler(c *gin.Context) {
	apikey := c.Param("apikey")

	if err := DeleteClient(apikey); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "ok")
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

	for k, v := range updates {
		viper.Set(k, v)
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
	var req struct {
		APIKey string   `json:"apikey"`
		IPs    []string `json:"ips"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if !VerifyClient(req.APIKey) {
		c.JSON(401, gin.H{"error": "Invalid API key"})
		return
	}

	threats, err := GetThreatsByIPs(req.IPs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, threats)
}

func init() {
	addr := viper.GetString("web")
	// Disable web server if is not set
	if addr == "" {
		return
	}

	rand.Read(secret)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/api/login", loginHandler)
	r.GET("/api/netstats", netstatsHandler)
	r.GET("/api/threats", authHandler, threatsHandler)
	r.DELETE("/api/threats/:ip", authHandler, deleteThreatsHandler)
	r.GET("/api/clients", authHandler, clientsHandler)
	r.POST("/api/clients", authHandler, createClientHandler)
	r.DELETE("/api/clients/:apikey", authHandler, deleteClientHandler)
	r.GET("/api/config", authHandler, readConfigHandler)
	r.POST("/api/config", authHandler, writeConfigHandler)
	r.POST("/api/check", checkHandler)

	r.NoRoute(staticHandler)

	log.Printf("Web server started on http://%s/\n", addr)
	go r.Run(addr)
}
