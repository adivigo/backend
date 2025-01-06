package controllers

import (
	"latihan_gin/models"

	"github.com/gin-gonic/gin"
)

// Tickitz
// @Summary get profile by id
// @Schemes
// @Description untuk mendapatkan profile
// @Tags Profiles
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} Response{results=models.User}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /profiles [get]
func GetProfileById(ctx *gin.Context) {
	val, isAvailable := ctx.Get("userId")
	userId := int(val.(float64))

	profile := models.FindOneProfile(userId)

	if isAvailable {
		ctx.JSON(200, Response{
			Success: true,
			Message: "Detail user",
			Results: profile,
		})
	}
}