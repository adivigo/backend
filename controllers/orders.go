package controllers

import (
	"fmt"
	"latihan_gin/models"
	"net/http/httputil"
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
	err := ctx.ShouldBind(&order)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(400, Response{
			Success: false,
			Message: "order gagal",
		})
		return
	}
	dump,_:=httputil.DumpRequest(ctx.Request,true)
	fmt.Printf("dump: %q\n", dump)
	fmt.Println("ini order",order)

	val, isAvailable := ctx.Get("userId")
    fmt.Println("user id:",val)
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

func GetAllPayment(c *gin.Context) {

	var showPayment []models.PaymentResponse
	showPayment = models.FindAllPayment()

	c.JSON(200, Response{
		Success: true,
		Message: "List of payment",
		Results: showPayment,
	})
}