package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), secureHeaders())

	r.LoadHTMLGlob("templates/*")

	// r.GET("/home", handleIndex)     // Home page
	// r.GET("/", showLoginPage)       // Login page
	// r.GET("/logout", showLoginPage) // Logout
	// r.GET("/healthz", health) // Healthcheck

	// Public routes
	r.GET("/login", showLoginPage)
	r.POST("/login", doLogin)
	r.GET("/logout", doLogout)
	r.GET("/healthz", health)

	// Protected group
	auth := r.Group("/")
	auth.Use(ginAuthMiddleware())
	registerRoutes(auth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("listening on %s", port)

	er := r.Run(port)
	if er != nil {
		log.Fatal(er)
	}

}
