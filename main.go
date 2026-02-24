package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", handleIndex) // Home page
	r.GET("/", showLoginPage)   // Login page
	r.GET("/healthz", health)   // Healthcheck

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("listening on %s", port)

	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}

}
