package routers

import (
	docs "latihan_gin/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
  	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	MovieRouter(router.Group("/movies"))
	AuthRouter(router.Group("/auth"))
	UserRouter(router.Group("/users"))
	ProfileRouter(router.Group("/profiles"))
	OrderRouter(router.Group("/orders"))
}