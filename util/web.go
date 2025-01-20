package util

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	page, _ := strconv.Atoi(c.Param("page"))
	netstats, err := GetNetstats(viper.GetInt("page_size"), page)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, netstats)

}

func threatsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	records, err := GetThreats(viper.GetInt("page_size"), page)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, records)
}
func readConfigHandler(c *gin.Context) {
	settings := viper.AllSettings()

	if checkPeriod, ok := settings["check_period"].(string); ok {
		if strings.HasSuffix(checkPeriod, "s") {
			checkPeriodWithoutUnit := strings.TrimSuffix(checkPeriod, "s")
			if value, err := strconv.Atoi(checkPeriodWithoutUnit); err == nil {
				settings["check_period"] = value
			} else {
				fmt.Println("无法转换 check_period 为整数", err)
			}
		}
	}

	c.JSON(200, settings)
}

// 删除防火墙规则（自动适配 Windows / Linux）
func DelFireWall(ip string) {
	if getOS() == "windows" {
		in := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=netvigil_block_in_"+ip)
		out := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=netvigil_block_out_"+ip)

		if err := in.Run(); err != nil {
			log.Printf("Failed to delete inbound firewall rule for %s: %v\n", ip, err)
		} else {
			log.Printf("Inbound firewall rule deleted for %s\n", ip)
		}

		if err := out.Run(); err != nil {
			log.Printf("Failed to delete outbound firewall rule for %s: %v\n", ip, err)
		} else {
			log.Printf("Outbound firewall rule deleted for %s\n", ip)
		}
	} else { // Linux
		in := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
		out := exec.Command("iptables", "-D", "OUTPUT", "-d", ip, "-j", "DROP")

		if err := in.Run(); err != nil {
			log.Printf("Failed to delete inbound iptables rule for %s: %v\n", ip, err)
		} else {
			log.Printf("Inbound iptables rule deleted for %s\n", ip)
		}

		if err := out.Run(); err != nil {
			log.Printf("Failed to delete outbound iptables rule for %s: %v\n", ip, err)
		} else {
			log.Printf("Outbound iptables rule deleted for %s\n", ip)
		}
	}
}

func threatsOperationHandler(c *gin.Context) {
	var req struct {
		ID     int    `json:"id"`
		Action string `json:"action"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Action == "remove" {
		var ip string
		err := DB.QueryRow("SELECT ip FROM threats WHERE ROWID = ?", req.ID).Scan(&ip)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to find threat record"})
			fmt.Println("500 SELECT:", err)
			return
		}

		DelFireWall(ip)

		result, err := DB.Exec("DELETE FROM threats WHERE ROWID = ?", req.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true, "message": fmt.Sprintf("Threat with IP %s removed", ip), "result": result})
	} else {
		c.JSON(400, gin.H{"error": err.Error()})
	}
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
		Token string   `json:"token"`
		IPs   []string `json:"ips"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
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
	r.GET("/api/netstats", authHandler, netstatsHandler)
	r.GET("/api/threats", authHandler, threatsHandler)
	r.POST("/api/threats", authHandler, threatsOperationHandler)
	r.GET("/api/config", authHandler, readConfigHandler)
	r.POST("/api/config", authHandler, writeConfigHandler)
	r.POST("/api/check", checkHandler)

	r.NoRoute(staticHandler)

	log.Printf("Web server started on http://%s/\n", addr)
	go r.Run(addr)
}
