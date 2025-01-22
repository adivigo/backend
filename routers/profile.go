package routers

import (
	"latihan_gin/controllers"
	"latihan_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func ProfileRouter(router *gin.RouterGroup) {
	router.Use(middlewares.ValidateToken())
	router.GET("", controllers.GetProfileById)
	router.PATCH("", controllers.EditProfile)
}