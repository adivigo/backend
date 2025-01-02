package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/lpernett/godotenv"
)

func DB() *pgx.Conn {
	godotenv.Load()
	urlconnect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PASSWORD"),
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_PORT"),
	os.Getenv("POSTGRES_DB"))
	connString := urlconnect
	conn, _ := pgx.Connect(context.Background(), connString)
	return conn
}