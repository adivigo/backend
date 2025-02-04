package models

import (
	"context"
	"fmt"
	"latihan_gin/lib"
	"time"

	"github.com/jackc/pgx/v5"
)

type Transactions struct {
	TransactionBody
	TransactionResponse
	Barcode    string   `json:"barcode"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	
}

type TransactionBody struct {
	Id         int       `json:"id"`
	CinemaId   int      `json:"cinemaId" form:"cinemaId"`
	MovieId    int      `json:"movieId" form:"movieId"`
	Seats     []string `json:"seats" form:"seats" binding:"-"` // Important: "seats[]"
	PaymentId  int      `json:"paymentId" form:"paymentId"`
	UserId     int      `json:"userId"`
}

type TransactionResponse struct {
	VirtualId  string
	TotalPrice int      `json:"totalPrice"`
	ExpiryDate time.Time
	Status string
}

type PaymentResponse struct {
	Id  int `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
}

func AddOrder(order TransactionBody, userId, virtual, totalPrice int, expiryDate time.Time, status string) TransactionResponse{
	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	var newlyCreated TransactionResponse

	err := conn.QueryRow(context.Background(), `
		INSERT INTO orders (cinema_id, movie_id, seats, payment_id, user_id, virtual_id, total_price, expiry_date, status) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING virtual_id, total_price, expiry_date, status
	`, order.CinemaId, order.MovieId, order.Seats, order.PaymentId, userId, virtual, totalPrice, expiryDate, status).Scan(
		&newlyCreated.VirtualId, 
		&newlyCreated.TotalPrice, 
		&newlyCreated.ExpiryDate, 
		&newlyCreated.Status,
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	return newlyCreated
}

func FindAllPayment() []PaymentResponse {
	conn, erro := lib.DB()
	if erro != nil {
		fmt.Println(erro)
	}
	defer conn.Close(context.Background())

	var rows pgx.Rows
	var err error

	rows, _ = conn.Query(context.Background(), `
			select payments.id, payments.name, payments.image from payments order by id asc
		`)

	payments, err := pgx.CollectRows(rows, pgx.RowToStructByName[PaymentResponse])
	if err != nil {
		fmt.Println(err)
	}

	return payments
}