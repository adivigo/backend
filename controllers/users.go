package controllers

import (
	"fmt"
	"latihan_gin/models"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pilinux/argon2"
)

// Tickitz
// @Summary get all list users
// @Schemes
// @Description untuk mendapatkan list semua user
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param search query string false " "
// @Param sortBy query string false " "
// @Param sortOrder query string false " "
// @Param page query string false " "
// @Param limit query string false " "
// @Success 200 {object} Response{results=models.ListUsers}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	showUser := models.FindAllUsers(search, sortBy, sortOrder, page, pageLimit)
	total := models.CountUsers(search)

	totalPage := int(math.Ceil(float64(total) / float64(pageLimit)))

	nextPage := page + 1
	if nextPage > totalPage {
		nextPage = totalPage
	}

	prevPage := page - 1
	if prevPage <= 0 {
		prevPage = 1
	}

	c.JSON(200, Response{
		Success: true,
		Message: "List of users",
		PageInfo: PageInfo{
			CurrentPage: page,
			NextPage: nextPage,
			PrevPage: prevPage,
			TotalPage: totalPage,
			TotalData: total,
		},
		Results: showUser,
	})
}

// Tickitz
// @Summary get users by id
// @Schemes
// @Description untuk mendapatkan user
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} Response{results=models.User}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func GetUserById(ctx *gin.Context) {
	paramId, err := strconv.Atoi(ctx.Param("id")) 
	if err != nil {
		ctx.JSON(400, Response{
			Success: false,
			Message: "Wrong id format",
		})
		return
	}

	oneUser := models.FindOneUser(paramId)

	if oneUser == (models.User{}) {
		ctx.JSON(404, Response{
			Success: false,
			Message: "User not Found",
		})
		return
	}

	ctx.JSON(200, Response{
		Success: true,
		Message: "Detail user",
		Results: oneUser,
	})
}

// Tickitz
// @Summary delete user by id
// @Schemes
// @Description untuk menghapus user
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} Response{results=models.ListUsers}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /users/{id} [DELETE]
func DeleteUser(c *gin.Context){
	paramId, _ := strconv.Atoi(c.Param("id")) 
	user := models.FindOneUser(paramId)
    if user == (models.User{}) {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "user not found",
		})
		return
	}

	deleted := models.DeleteUser(paramId)

	c.JSON(http.StatusNotFound, Response{
		Success: true,
		Message: "delete user success",
		Results: deleted,
	})
}

// Tickitz
// @Summary update user
// @Schemes
// @Description untuk merubah data user
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path int true "user id"
// @Param email formData string true " "
// @Param password formData string true " "
// @Param firstName formData string true " "
// @Param lastName formData string true " "
// @Param phoneNumber formData string true " "
// @Param image formData file true " "
// @Success 200 {object} Response{results=models.ListUsers}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /users/{id} [PATCH]
func EditUser(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id")) 
	user := models.FindOneUser(paramId)

	if user == (models.User{}) {
		c.JSON(404, Response{
			Success: false,
			Message: "User not Found",
		})
		return
	}

	c.ShouldBind(&user)

	f, _ := c.MultipartForm()
	file, _ := c.FormFile("image")
	
	user.PhoneNumber = f.Value["phone_number"][0]

	if file.Filename != "" {
		if file.Size > MaxFileSize {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "File tidak boleh lebih dari 8MB",
			})
			return
		}
		filename := uuid.New().String()
		splittedFilename := strings.Split(file.Filename, ".")
		ext := splittedFilename[len(splittedFilename) - 1]
		if ext != "jpg" && ext != "png" {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "format file bukan png atau jpg",
			})
			return
		}
		storedFile := fmt.Sprintf("%s.%s", filename, ext)
		c.SaveUploadedFile(file, fmt.Sprintf("uploads/movie/%s", storedFile))
		user.Image = &storedFile
	}
	
	if !strings.Contains(user.Password, "$argon2i$v=19$m=65536,t=1,p=2$") {
		if user.Password != "" {
			user.Password, _ = argon2.CreateHash(user.Password, "", argon2.DefaultParams)
		}
	}

	updated := models.UpdateUser(user)
		fmt.Println(updated)

	c.JSON(200, Response{
		Success: true,
		Message: "Profile Updated",
	})
}

// Tickitz
// @Summary create user
// @Schemes
// @Description untuk membuat user baru
// @Tags Users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true " "
// @Param password formData string true " "
// @Success 200 {object} Response{results=models.ListUsers}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /users [POST]
func CreateUser(c *gin.Context) {
	var user models.User
	c.ShouldBind(&user)
	
	if user.Password != "" {
		user.Password, _ = argon2.CreateHash(user.Password, "", argon2.DefaultParams)
	}

	newUser := models.InsertUser(user)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "User Added Successfully",
		Results: newUser,
	})
}