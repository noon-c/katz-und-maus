package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", handleIndex)   // Home page
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
