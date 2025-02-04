package routers

import (
	"latihan_gin/controllers"

	"github.com/gin-gonic/gin"
)

func SeatRouter(router *gin.RouterGroup) {
	router.GET("", controllers.GetAllSeats)
}