package middlewares

import (
	"latihan_gin/controllers"
	"latihan_gin/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckInput() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		var formUser models.User
		ctx.ShouldBind(&formUser)

		if !strings.Contains(formUser.Email , "@") || len(formUser.Password) < 4  {
			ctx.JSON(http.StatusBadRequest, controllers.Response{
				Success: false,
				Message: "Email atau password tidak valid",
			})
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}