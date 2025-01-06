package lib

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func DB() (*pgx.Conn, error) {
	godotenv.Load()
	
	config, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(context.Background(), config.ConnString())

	if err != nil {
		return nil, err
	}

	return conn, nil

	// urlconnect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
	// os.Getenv("POSTGRES_USER"),
	// os.Getenv("POSTGRES_PASSWORD"),
	// os.Getenv("POSTGRES_HOST"),
	// os.Getenv("POSTGRES_PORT"),
	// os.Getenv("POSTGRES_DB"))
	// connString := urlconnect
	// return conn

}