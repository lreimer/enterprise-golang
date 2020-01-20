package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// configuration for static files and templates
	engine.LoadHTMLFiles("templates/index.html")
	engine.StaticFile("/favicon.ico", "favicon.ico")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "A Hitchhiker's Guide to Enterprise Microservices with Go",
		})
	})

	// run server on PORT
	engine.Run(port())
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
