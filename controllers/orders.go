package controllers

import (
	"fmt"
	"latihan_gin/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Tickitz
// @Summary place order
// @Schemes
// @Description untuk memesan tiket
// @Tags Orders
// @Accept x-www-form-urlencoded
// @Produce json
// @Param cinemaId formData int false " "
// @Param movieId formData int false " "
// @Param paymentId formData int false " "
// @Param seats[] formData []string false " "
// @Success 200 {object} Response{results=models.TransactionResponse}
// @Failure 401 {object} Response401
// @Security ApiKeyAuth
// @Router /orders [post]
func PlaceOrder(ctx *gin.Context) {
    var order models.TransactionBody
	ctx.ShouldBind(&order)

	val, isAvailable := ctx.Get("userId")
    fmt.Println(val)
	userId := int(val.(float64))

    virtualId :=int(time.Now().UnixNano()/(1<<22))
    virtual := int(virtualId)
	
	var arrayseat []string
	for _, v := range order.Seats {
		splitSeats := strings.Split(v, ",")
		arrayseat = append(arrayseat, splitSeats...)
	}
	order.Seats = arrayseat 
	fmt.Println(order.Seats)
	
    totalPrice := len(order.Seats)*50000
    expiryDate := time.Now().Add(3 * 24 * time.Hour)
    status := "pending"

    orders := models.AddOrder(order, userId, virtual, totalPrice, expiryDate, status)

	if isAvailable {
		ctx.JSON(200, Response{
			Success: true,
			Message: "Detail order tiket",
            Results: orders,
		})
	}
}