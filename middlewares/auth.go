package middlewares

import (
	"latihan_gin/controllers"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/joho/godotenv"
)

func ValidateToken() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		head := ctx.GetHeader("Authorization")
		splitheader := strings.Split(head, " ")
		if len(splitheader) != 2 || splitheader[0] != "Bearer" {
				ctx.JSON(http.StatusUnauthorized, controllers.Response{
					Success: false,
					Message: "User Unauthorized",
				})
				ctx.Abort()
				return
			}
			token := splitheader[1]
			tok, _ := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.HS256})
			
			out := make(map[string]any)
			
		godotenv.Load()
		var JWT_SECRET []byte = []byte(controllers.GetMd5Hash(os.Getenv("JWT_SECRET")))
			
		erro := tok.Claims(JWT_SECRET, &out)

		ctx.Set("userId", out["userId"].(float64))

		if erro != nil {
			ctx.JSON(http.StatusUnauthorized, controllers.Response{
				Success: false,
				Message: "User Unathorized",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}