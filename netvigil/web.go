package netvigil

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

type Config struct {
	CaptureInterval time.Duration `toml:"capture_interval"`
	CheckInterval   time.Duration `toml:"check_interval"`
	// Username        string        `toml:"username"`
	// Password        string        `toml:"password"`
	Tix    []Tix   `toml:"tix"`
	Web    string  `toml:"web"`
	Routes []Route `toml:"routes"`
}

type Tix struct {
	Type      string   `toml:"type"`
	APIKey    string   `toml:"apikey"`
	Blacklist []string `toml:"blacklist"`
}

type Route struct {
	Method string `mapstructure:"method"`
	Path   string `mapsturcture:"path"`
}

var secret = make([]byte, 32)
var configFilePath = "config.toml"

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

func configHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		config, err := readConfig()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, config)
	} else if c.Request.Method == http.MethodPost {
		var newConfig Config
		if err := c.BindJSON(&newConfig); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON payload"})
			return
		}
		if err := writeConfig(newConfig); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Config updated successfully"})
	} else {
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

func readConfig() (Config, error) {
	var config Config
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	if err := toml.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, nil
}

func writeConfig(config Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return err
	}
	return nil
}

// authHandler 鉴权中间件
func authHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Next()
}

func init() {
	rand.Read(secret)
	addr := viper.GetString("web")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	routes := viper.Get("routes").([]interface{})
	for _, route := range routes {
		routeMap := route.(map[string]interface{})
		method := routeMap["method"].(string)
		path := routeMap["path"].(string)

		switch method {
		case "GET":
			r.GET(path, authHandler, configHandler)
		case "POST":
			r.POST(path, authHandler, configHandler)
		}
	}

	// 登录和记录路由
	r.POST("/api/login", loginHandler)
	r.GET("/api/records", recordsHandler)

	// 静态文件处理
	r.Use(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		switch path {
		case "/api/records", "/api/login", "/api/config":
			ctx.Next()
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
