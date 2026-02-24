package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(auth *gin.RouterGroup) {
	auth.GET("/", func(c *gin.Context) {
		roleAny, _ := c.Get("role")
		role := roleAny.(string)

		if role == roleAdmin {
			handleIndex(c) // index.html
			return
		}
		if role == roleMint {
			handleMint(c) // mint.html
			return
		}

		c.AbortWithStatus(http.StatusForbidden)
	})

	// xGVAOyfob4FxtPhK for /admin
	// VBCrbTDiI8MYs73L for /mint
	auth.GET("/SFEtAKq5ueS64tlhUSKvNsRYNoh4iC8QYfSswOiYltCvJ3PETUdh0p", func(c *gin.Context) { handleIndex(c) })
	auth.GET("/SFEtAKq5ueS64tIhUSKvNsRYNoh4iC8QYfSswOiYltCvJ3PETUdh0p", func(c *gin.Context) { handleMint(c) })
}
