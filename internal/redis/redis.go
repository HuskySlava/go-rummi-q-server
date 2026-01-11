package redis

import (
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

var client *redis.Client

type config struct {
	addr     string
	password string
	db       int
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func loadConfig() config {

	// redis.Options.db expects int, env provides string
	dbStr := getEnv("REDIS_DB", "0")
	dbInt, err := strconv.Atoi(dbStr)
	if err != nil {
		dbInt = 0
	}

	return config{
		addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		password: getEnv("REDIS_PASSWORD", ""),
		db:       dbInt,
	}
}

func Init() error {

	if os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	cfg := loadConfig()

	client = redis.NewClient(&redis.Options{
		Addr:     cfg.addr,
		Password: cfg.password,
		DB:       cfg.db,
	})

	return nil
}
