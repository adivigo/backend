package models

import (
	"context"
	"fmt"
	"latihan_gin/lib"
	"time"
)

type Profile struct {
	Id         int       `json:"id"`
	Email      string    `json:"email" form:"email"`
	Password      string    `json:"password" form:"password"`
	FirstName   *string    `json:"firstName" form:"first_name"`
	LastName   *string    `json:"lastName" form:"last_name"`
	PhoneNumber   string    `json:"phoneNumber" form:"phone_number"`
	Image   string    `json:"image" form:"image"`
	Role *string `json:"role" form:"role,omitempty"`
	Point *string `json:"point" form:"point,omitempty"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

func FindOneProfile(paramId int) Profile {
	var profile Profile

	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	err := conn.QueryRow(context.Background(), `
		SELECT id, email, first_name, last_name, phone_number, image, point, created_at, updated_at
		FROM users WHERE id = $1
	`, paramId).Scan(
		&profile.Id,
		&profile.Email,
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

func UpdateProfile(profile Profile) Profile {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())
	
	var updatedProfile Profile

	conn.QueryRow(context.Background(), `
		UPDATE users SET email=$1, password=$2, first_name=$4, last_name=$5, phone_number=$6, image=$7 WHERE id = $3
		RETURNING id, email, password, first_name, last_name, phone_number, image, created_at, updated_at
	`, profile.Email, profile.Password, profile.Id, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.Image).Scan(
		&updatedProfile.Id, 
		&updatedProfile.Email, 
		&updatedProfile.Password, 
		&updatedProfile.FirstName, 
		&updatedProfile.LastName, 
		&updatedProfile.PhoneNumber, 
		&updatedProfile.Image, 
		&updatedProfile.CreatedAt, 
		&updatedProfile.UpdatedAt,
	)
	return updatedProfile
}