package routers

import (
	"latihan_gin/controllers"
	"latihan_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.RouterGroup) {
	router.Use(middlewares.ValidateToken())
	router.POST("", controllers.PlaceOrder)
	router.GET("/payment", controllers.GetAllPayment)
}