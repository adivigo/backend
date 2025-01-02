package routers

import (
	"latihan_gin/controllers"
	"latihan_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup) {
	router.Use(middlewares.CheckInput())
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
}