package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		roleAny, _ := c.Get("role")
		role := roleAny.(string)

		switch role {
		case roleAdmin:
			handleIndex(c)
		case roleMint:
			handleMint(c)
		default:
			c.AbortWithStatus(http.StatusForbidden)
		}
	})

	rg.GET("/admin", func(c *gin.Context) {
		roleAny, _ := c.Get("role")
		if roleAny.(string) != roleAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		handleIndex(c)
	})

	rg.GET("/mint", func(c *gin.Context) {
		roleAny, _ := c.Get("role")
		if roleAny.(string) != roleMint {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		handleMint(c)
	})
}
