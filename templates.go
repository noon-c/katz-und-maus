package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Error": nil})
}

func handleMint(c *gin.Context) {
	c.HTML(http.StatusOK, "mint.html", gin.H{"Error": nil})
}

func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"Error": nil})
}
