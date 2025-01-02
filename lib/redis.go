package lib

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
	godotenv.Load()
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	return rdb
}