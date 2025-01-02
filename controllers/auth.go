package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"latihan_gin/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/lpernett/godotenv"
	"github.com/pilinux/argon2"
)

func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Tickitz
// @Summary Register user
// @Schemes
// @Description untuk membuat user baru
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true " "
// @Param password formData string true " "
// @Success 200 {object} Response{results=models.ListUsers}
// @Router /auth/register [POST]
func Register(c *gin.Context) {
	var user models.User
	c.ShouldBind(&user)

	if user.Password != "" {
		user.Password, _ = argon2.CreateHash(user.Password, "", argon2.DefaultParams)
	}

	newUser  := models.InsertUser(user)
    if newUser.Id == 0 {
        c.JSON(http.StatusInternalServerError, Response{
            Success: false,
            Message: "Failed to Register",
        })
        return
    }

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Register Success",
	})
}

// Tickitz
// @Summary Login user
// @Schemes
// @Description untuk mengautentikasi user
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true " "
// @Param password formData string true " "
// @Success 200 {object} Response{results=models.ListUsers}
// @Router /auth/login [POST]
func Login(ctx *gin.Context) {
	var formUser  models.User

	if err := ctx.ShouldBind(&formUser ); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	foundUser := models.FindOneUserByMail(formUser.Email)

	if foundUser  == (models.User{}) {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
		})
		return
	}

	match, err := argon2.ComparePasswordAndHash(formUser.Password, "", foundUser.Password)
	if err != nil || !match {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Wrong email or password",
		})
		return
	}
	
	godotenv.Load()

	var JWT_SECRET []byte = []byte(GetMd5Hash(os.Getenv("JWT_SECRET")))

	signer, _ := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: JWT_SECRET}, (nil))
	baseInfo := jwt.Claims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	payload := struct {
		UserId int `json:"userId"`
	}{
		UserId: foundUser.Id,
	}

	token, _ := jwt.Signed(signer).Claims(baseInfo).Claims(payload).Serialize()

	tok := models.Token{
		Token: token,
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login success",
		Results: tok,
	})
}