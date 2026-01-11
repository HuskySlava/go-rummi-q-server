package redis

import (
	"github.com/joho/godotenv"
	"os"
)

type config struct {
	addr     string
	password string
	db       string
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func loadConfig() config {
	return config{
		addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		password: getEnv("REDIS_PASSWORD", ""),
		db:       getEnv("REDIS_DB", "0"),
	}
}

func Init() error {

	if os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}
	
	return nil
}
