package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	list, err := MapList()
	if err != nil {
		log.Fatalf("MapList err: %v", err.Error())
		return
	}

	listJson, err := json.Marshal(list)
	if err != nil {
		log.Fatalf("MapList err: %v", err.Error())
		return
	}

	if err := os.MkdirAll("./result/", 0755); err != nil {
		log.Fatalf("os.MkdirAll err: %v", err.Error())
		return
	}

	r.StaticFile("/static/js/jquery.js", "./static/js/jquery-3.6.3.min.js")
	r.StaticFile("/static/js/bootstrap.js", "./static/js/bootstrap.min.js")
	r.StaticFile("/static/css/bootstrap.css", "./static/css/bootstrap.min.css")

	r.Static("./result/", "./result/")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"list": string(listJson),
		})
	})

	r.POST("/voice", func(c *gin.Context) {
		var params VoiceParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		fileName, err := Voice(&params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"file_name": fileName,
		})
	})

	r.Run(":8080")
}
