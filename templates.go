package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func renderHomePage(c *gin.Context) {
// 	c.Header("Content-Type", "text/html; charset=utf-8")
// 	c.String(http.StatusOK, `<!doctype html>
// 			<html lang="en">
// 				<head>
// 					<meta charset="utf-8">
// 					<title>GO</title>
// 				</head>
// 				<body>
// 					<h1>Running on GO</h1>
// 					<p>Server is up and running!</p>
// 				</body>
// 			</html>`)
// }

func health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Error": nil})
}
