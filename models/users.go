package models

import (
	"context"
	"fmt"
	"latihan_gin/lib"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" form:"email"`
	Password   string    `json:"password" form:"password"`
	FirstName   *string    `json:"firstName" form:"first_name"`
	LastName   *string    `json:"lastName" form:"last_name"`
	PhoneNumber   string    `json:"phoneNumber" form:"phone_number"`
	Image   *string    `json:"image" form:"image"`
	Role *string `json:"role" form:"role,omitempty"`
	Point *string `json:"point" form:"point,omitempty"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type Token struct {
	Token string `json:"token" form:"token"`
}

type ListUsers []User

func FindAllUsers(search, sortBy, sortOrder string, page, pageLimit int) ListUsers {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	var rows pgx.Rows
	var err error

	offset := (page - 1) * pageLimit

	if search != "" {
		rows, err = conn.Query(context.Background(), `
			SELECT id, password, email, first_name, last_name, phone_number, image, role, point, created_at, updated_at
			FROM users
			WHERE email ILIKE $1
			ORDER BY ` + sortBy + ` ` + sortOrder + `
			LIMIT $2 OFFSET $3
		`, "%"+search+"%", pageLimit, offset)
	} else {
		rows, err = conn.Query(context.Background(), `
			SELECT id, password, email, first_name, last_name, phone_number, image, role, point, created_at, updated_at
			FROM users
			ORDER BY ` + sortBy + ` ` + sortOrder + `
			LIMIT $1 OFFSET $2
		`, pageLimit, offset)
	}

	if err != nil {
		fmt.Println(err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		fmt.Println(err)
	}

	return users
}

func FindOneUser(paramId int) User {
	var user User

	conn, _ := lib.DB()
	defer conn.Close(context.Background())


	conn.QueryRow(context.Background(), `
		SELECT id, email, password, first_name, last_name, phone_number, image, created_at, updated_at
		FROM users WHERE id = $1
	`, paramId).Scan(
		&user.Id, 
		&user.Email, 
		&user.Password, 
		&user.FirstName, 
		&user.LastName, 
		&user.PhoneNumber, 
		&user.Image, 
		&user.CreatedAt, 
		&user.UpdatedAt)
	return user
}

func FindOneUserByMail(email string) User {
	var user User

	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	conn.QueryRow(context.Background(), `
		SELECT id, email, password, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	
	return user
}

func InsertUser(user User) User {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())
	
	var newlyCreated User

	conn.QueryRow(context.Background(), `
		INSERT INTO users (email, password) VALUES
		($1, $2)
		RETURNING id, email, password, first_name, last_name, phone_number, image, role, point, created_at, updated_at
	`, user.Email, user.Password).Scan(
		&newlyCreated.Id, 
		&newlyCreated.Email,
		&newlyCreated.Password, 
		&newlyCreated.FirstName, 
		&newlyCreated.LastName, 
		&newlyCreated.PhoneNumber, 
		&newlyCreated.Image, 
		&newlyCreated.Role, 
		&newlyCreated.Point, 
		&newlyCreated.CreatedAt, 
		&newlyCreated.UpdatedAt,
	)
	return newlyCreated
}

func UpdateUser(user User) User {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())
	
	var updatedUser User

	conn.QueryRow(context.Background(), `
		UPDATE users SET email=$1, password=$2, first_name=$4, last_name=$5, phone_number=$6, image=$7 WHERE id = $3
		RETURNING id, email, password, first_name, last_name, phone_number, image, created_at, updated_at
	`, user.Email, user.Password, user.Id, user.FirstName, user.LastName, user.PhoneNumber, user.Image).Scan(
		&updatedUser.Id, 
		&updatedUser.Email, 
		&updatedUser.Password, 
		&updatedUser.FirstName, 
		&updatedUser.LastName, 
		&updatedUser.PhoneNumber, 
		&updatedUser.Image, 
		&updatedUser.CreatedAt, 
		&updatedUser.UpdatedAt,
	)
	return updatedUser
}

func DeleteUser(id int) User{
	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	deletedUser := User{}

	conn.QueryRow(context.Background(),`
		DELETE FROM users
		WHERE id = $1
		RETURNING id, email, password, created_at, updated_at
	`, id).Scan(
		&deletedUser.Id,
		&deletedUser.Email,
		&deletedUser.Password,
		&deletedUser.CreatedAt,
		&deletedUser.UpdatedAt,
	)
	return deletedUser
}

func CountUsers(search string) int {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	var total int

	search = fmt.Sprintf("%%%s%%", search)

	conn.QueryRow(context.Background(), `
		SELECT COUNT(id) FROM users
		WHERE email ILIKE $1
	`, search).Scan(&total)

	return total
}