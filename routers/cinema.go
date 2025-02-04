package routers

import (
	"latihan_gin/controllers"

	"github.com/gin-gonic/gin"
)

func CinemaRouter(router *gin.RouterGroup) {
	router.GET("", controllers.GetAllCinemas)
	router.GET("/:id", controllers.GetCinemaById)
}