package cache

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() *redis.Client {
	// Load environment variables from .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Get environment variable
	server := os.Getenv("REDIS_URL")

	return redis.NewClient(&redis.Options{
		Addr: server,
	})
}
