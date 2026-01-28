package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", runningPage)
	r.GET("/healthz", health) // Healthcheck

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started: http://localhost:8080")
	log.Printf("listening on %s", port)
}

func runningPage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!doctype html>
			<html lang="en">
				<head>
					<meta charset="utf-8">
					<title>running...</title>
				</head>
				<body>
					<h1>running...</h1>
				</body>
			</html>`)
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
