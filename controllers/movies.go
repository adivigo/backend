package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"latihan_gin/lib"
	"latihan_gin/models"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PageInfo struct {
	CurrentPage int `json:"currentPage,omitempty"`
	NextPage    int `json:"nextPage,omitempty"`
	PrevPage    int `json:"prevPage,omitempty"`
	TotalPage   int `json:"totalPage,omitempty"`
	TotalData   int `json:"totalData,omitempty"`
}

type Response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	PageInfo any    `json:"pageInfo,omitempty"`
	Results  any    `json:"results,omitempty"`
}

type Response401 struct {
	Success bool   `json:"success" default:"False"`
	Message string `json:"message" default:"Unauthorized"`
}

var MaxFileSize int64 = 2 << 20

// Tickitz
// @Summary get all list movies
// @Schemes
// @Description untuk mendapatkan list semua movie
// @Tags Movies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param search query string false " "
// @Param sortBy query string false " "
// @Param sortOrder query string false " "
// @Param page query string false " "
// @Param limit query string false " "
// @Success 200 {object} Response{results=models.ListMovies}
// @Failure 401 {object} Response401
// @Router /movies [get]
func GetAllMovies(c *gin.Context) {
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	var showMovie []models.MovieData
	var count int

	get := lib.Redis().Get(context.Background(), c.Request.RequestURI)
	getCount := lib.Redis().Get(context.Background(), fmt.Sprintf("count+%s", c.Request.RequestURI))
	if get.Val() != "" {
		rawData := []byte(get.Val())
		json.Unmarshal(rawData, &showMovie)
	} else {
		showMovie = models.FindAllMovies(search, sortBy, sortOrder, page, pageLimit)
		encoded, _ := json.Marshal(showMovie)
		lib.Redis().Set(context.Background(), c.Request.RequestURI, string(encoded), 0) // <- 60s
	}

	if getCount.Val() != "" {
		rawData := []byte(getCount.Val())
		json.Unmarshal(rawData, &count)
	} else {
		count = models.CountMovies(search)
		encoded, _ := json.Marshal(count)
		lib.Redis().Set(context.Background(), fmt.Sprintf("count+%s", c.Request.RequestURI), string(encoded), 0)
	}

	totalPage := int(math.Ceil(float64(count) / float64(pageLimit)))

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
		Message: "List of movies",
		PageInfo: PageInfo{
			CurrentPage: page,
			NextPage:    nextPage,
			PrevPage:    prevPage,
			TotalPage:   totalPage,
			TotalData:   count,
		},
		Results: showMovie,
	})
}

// Tickitz
// @Summary get list movie by id
// @Schemes
// @Description untuk mendapatkan list movie dari id
// @Tags Movies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path int true "movie id"
// @Success 200 {object} Response{results=models.Movie}
// @Failure 401 {object} Response401
// @Router /movies/{id} [get]
func GetMovieById(ctx *gin.Context) {
	paramId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, Response{
			Success: false,
			Message: "Wrong id format",
		})
		return
	}

	oneMovie := models.FindOneMovie(paramId)

	if oneMovie == (models.MovieData{}) {
		ctx.JSON(404, Response{
			Success: false,
			Message: "Movie not Found",
		})
		return
	}

	ctx.JSON(200, Response{
		Success: true,
		Message: "Detail movie",
		Results: oneMovie,
	})
}

// Tickitz
// @Summary create movie
// @Schemes
// @Description untuk membuat movie baru
// @Tags Movies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param title formData string true " "
// @Param image formData file true " "
// @Param banner formData file true " "
// @Param tag formData string true " "
// @Param releaseDate formData string true " "
// @Param duration formData string true " "
// @Param synopsis formData string true " "
// @Success 200 {object} Response{results=models.Movie}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /movies [POST]
func CreateMovie(c *gin.Context) {
	var form models.MovieBody

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid input",
		})
		fmt.Println(err)
		return
	}

	releaseDateStr := c.PostForm("release_date")
	releaseDate, err := time.Parse(time.DateOnly, releaseDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Format tanggal tidak valid",
		})
		return
	}

	savedImage, err := handleFileUpload(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	form.Image = savedImage

	savedBanner, err := handleFileUpload(c, "banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	form.Banner = savedBanner

	newMovie := models.InsertMovie(form, releaseDate)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Movie Added Successfully",
		Results: newMovie,
	})
}

func handleFileUpload(c *gin.Context, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("file not found: %s", fieldName)
	}

	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file tidak boleh lebih dari 2MB")
	}

	splittedFilename := strings.Split(file.Filename, ".")
	ext := splittedFilename[len(splittedFilename)-1]
	log.Println(ext)
	if ext != "jpg" && ext != "png" && ext != "jpeg" {
		return "", fmt.Errorf("format file bukan png/jpg/jpeg")
	}

	filename := uuid.New().String()
	storedFile := fmt.Sprintf("%s.%s", filename, ext)
	if err := c.SaveUploadedFile(file, fmt.Sprintf("uploads/images/%s", storedFile)); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return storedFile, nil
}

// Tickitz
// @Summary delete movie by id
// @Schemes
// @Description untuk menghapus movie
// @Tags Movies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path int true "movie id"
// @Success 200 {object} Response{results=models.Movie}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /movies/{id} [DELETE]
func DeleteMovie(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id"))
	movie := models.FindOneMovie(paramId)
	if movie == (models.MovieData{}) {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "movie not found",
		})
		return
	}

	deleted := models.RemoveMovie(paramId)

	c.JSON(http.StatusNotFound, Response{
		Success: true,
		Message: "delete movie success",
		Results: deleted,
	})
}

func UpdateMovie(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id"))
	movie := models.FindOneMovie(paramId)

	if movie == (models.MovieData{}) {
		c.JSON(404, Response{
			Success: false,
			Message: "movie not Found",
		})
		return
	}

	c.ShouldBind(&movie)

	savedImage, err := handleFileUpload(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	movie.Image = savedImage

	savedBanner, err := handleFileUpload(c, "banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	movie.Banner = savedBanner

	updated := models.UpdateMovie(movie)
		// fmt.Println(updated)

	c.JSON(200, Response{
		Success: true,
		Message: "Detail movie",
		Results: updated,
	})
}

func GetAllCinemas(c *gin.Context) {

	var showCinema []models.Cinemas
	showCinema = models.FindAllCinemas()

	c.JSON(200, Response{
		Success: true,
		Message: "List of cinema",
		Results: showCinema,
	})
}

func GetAllSeats(c *gin.Context) {

	var showSeats []models.Seats
	showSeats = models.FindAllSeats()

	c.JSON(200, Response{
		Success: true,
		Message: "List of seats",
		Results: showSeats,
	})
}

func GetCinemaById(ctx *gin.Context) {
	paramId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, Response{
			Success: false,
			Message: "Wrong id format",
		})
		return
	}

	oneCinema := models.FindOneCinema(paramId)

	if oneCinema == (models.Cinemas{}) {
		ctx.JSON(404, Response{
			Success: false,
			Message: "Cinema not Found",
		})
		return
	}

	ctx.JSON(200, Response{
		Success: true,
		Message: "Detail cinema",
		Results: oneCinema,
	})
}