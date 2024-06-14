package env

import "os"

func GetRedisUrl() string {
	value := os.Getenv("REDIS_URL")
	if value == "" {
		value = "redis://localhost:6379"
	}
	return value
}
