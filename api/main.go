package handler

import (
	"latihan_gin/routers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title TixIT
// @version 1.0
// @description backend TixIT

// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

// func main() {
// 	route := gin.New(func(e *gin.Engine) {
// 		e.RedirectTrailingSlash = false
// 	})

// 	route.MaxMultipartMemory = 2 << 20

// 	route.Static("/profiles/images", "uploads/images")
// 	route.Use(cors.New(cors.Config{
// 		AllowAllOrigins: true,
// 		AllowHeaders: []string{"Authorization", "Content-Type"},
// 		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
// 	}))

// 	routers.Routers(route)

// 	route.GET("/api/hello", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "Hello from Golang Backend!"})
// 	})

// 	route.Run()
// }

// Handler function (Vercel requires this)
func Handler(w http.ResponseWriter, r *http.Request) {
	route := gin.New()

	// Middleware & Config
	route.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowMethods:    []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	// Routes
	routers.Routers(route)
	route.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Golang Backend on Vercel!"})
	})

	// Serve HTTP request
	route.ServeHTTP(w, r)
}
