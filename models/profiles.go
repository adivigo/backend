package models

import (
	"context"
	"fmt"
	"latihan_gin/lib"
	"time"
)

type Profile struct {
	Id         int       `json:"id"`
	FirstName   *string    `json:"firstName" form:"first_name"`
	LastName   *string    `json:"lastName" form:"last_name"`
	PhoneNumber   string    `json:"phoneNumber" form:"phone_number"`
	Image   *string    `json:"image" form:"image"`
	Point *string `json:"point" form:"point,omitempty"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

func FindOneProfile(paramId int) Profile {
	var profile Profile

	conn := lib.DB()
	defer conn.Close(context.Background())

	err := conn.QueryRow(context.Background(), `
		SELECT id, first_name, last_name, phone_number, image, point, created_at, updated_at
		FROM users WHERE id = $1
	`, paramId).Scan(
		&profile.Id,
		&profile.FirstName,
		&profile.LastName,
		&profile.PhoneNumber,
		&profile.Image,
		&profile.Point,
		&profile.CreatedAt,
		&profile.UpdatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return profile
}