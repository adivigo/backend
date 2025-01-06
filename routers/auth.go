package routers

import (
	"latihan_gin/controllers"
	"latihan_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup) {
	router.POST("/register", middlewares.CheckInputRegister(), controllers.Register)
	router.POST("/login", middlewares.CheckInputLogin(), controllers.Login)
}