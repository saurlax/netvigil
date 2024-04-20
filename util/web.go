package util

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

func init() {
	addr := viper.GetString("web")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		switch path {
		case "/api/records":
			recordsHandler(ctx)
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
