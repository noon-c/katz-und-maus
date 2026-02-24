package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	hashMint := os.Getenv("MINT_PASS_HASH")

	hash := os.Getenv("ADMIN_PASS_HASH")
	pass := os.Getenv("MASTER_PASSWORD")

	fmt.Println(hash)
	fmt.Println(pass)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	fmt.Println("match:", err == nil, "err:", err)
	fmt.Println("hashMint:", hashMint)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), secureHeaders())

	r.LoadHTMLGlob("templates/*")

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
