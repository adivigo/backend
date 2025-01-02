package main

import (
	"latihan_gin/routers"

	"github.com/gin-gonic/gin"
)

// @title TixIT
// @version 1.0
// @description this is sample server

// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	route := gin.New(func (e *gin.Engine)  {
		e.RedirectTrailingSlash = false
	})

	route.MaxMultipartMemory = 8 << 20

	routers.Routers(route)

	route.Run("localhost:8888")
}