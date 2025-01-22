package routers

import (
	"latihan_gin/controllers"

	"github.com/gin-gonic/gin"
)

func MovieRouter(router *gin.RouterGroup) {
	router.GET("", controllers.GetAllMovies)
	router.GET("/:id", controllers.GetMovieById)
	router.POST("", controllers.CreateMovie)
	router.DELETE("/:id", controllers.DeleteMovie)
	router.PATCH("/:id", controllers.UpdateMovie)
	// router.PATCH("/:id", middlewares.ValidateToken(), controllers.UpdateMovie)
}