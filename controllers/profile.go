package controllers

import (
	"fmt"
	"latihan_gin/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pilinux/argon2"
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

func EditProfile(c *gin.Context) {
	val, _ := c.Get("userId")
	userId := int(val.(float64))
	
	profile := models.FindOneProfile(userId)
	c.ShouldBind(&profile)

	// f, _ := c.MultipartForm()
	file, _ := c.FormFile("image")

	profile.PhoneNumber = c.PostForm("phone_number")
	var oldImage string

	if file != nil && file.Filename != "" {
		oldImage = profile.Image

		if file.Size > MaxFileSize {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "File tidak boleh lebih dari 8MB",
			})
			return
		}
		filename := uuid.New().String()
		splittedFilename := strings.Split(file.Filename, ".")
		ext := splittedFilename[len(splittedFilename)-1]
		if ext != "jpg" && ext != "png" {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "format file bukan png atau jpg",
			})
			return
		}
		storedFile := fmt.Sprintf("%s.%s", filename, ext)
		c.SaveUploadedFile(file, fmt.Sprintf("uploads/images/%s", storedFile))
		profile.Image = storedFile

		if oldImage != "" {
			oldFilePath := fmt.Sprintf("uploads/images/%s", oldImage)
			os.Remove(oldFilePath)
		}
	} else {
		profile.Image = models.FindOneProfile(userId).Image
	}

	if !strings.Contains(profile.Password, "$argon2i$v=19$m=65536,t=1,p=2$") {
		if profile.Password != "" {
			profile.Password, _ = argon2.CreateHash(profile.Password, "", argon2.DefaultParams)
		}
	}

	updated := models.UpdateProfile(profile)

		c.JSON(200, Response{
			Success: true,
			Message: "Profile Updated",
			Results: updated,
		})
	
}