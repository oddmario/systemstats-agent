package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oddmario/systemstats-agent/routes"
)

func initRoutes(engine *gin.Engine) {
	engine.Any("/robots.txt", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")

		c.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	engine.GET("/stats", routes.Stats)
}
