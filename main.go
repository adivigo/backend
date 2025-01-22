package main

import (
	"latihan_gin/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title TixIT
// @version 1.0
// @description backend TixIT

// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	route := gin.New(func(e *gin.Engine) {
		e.RedirectTrailingSlash = false
	})

	route.MaxMultipartMemory = 2 << 20

	route.Static("/profiles/images", "uploads/images")
	route.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	routers.Routers(route)

	route.Run("localhost:8888")
}
