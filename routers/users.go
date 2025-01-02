package routers

import (
	"latihan_gin/controllers"
	"latihan_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.RouterGroup) {
	router.Use(middlewares.ValidateToken())
	router.GET("", controllers.GetAllUsers)
	router.GET("/:id", controllers.GetUserById)
	router.DELETE("/:id", controllers.DeleteUser)
	router.PATCH("/:id", controllers.EditUser)
	router.POST("", controllers.CreateUser)
}