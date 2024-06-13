package env

import "os"

func GetRabbitUrl() string {
	value := os.Getenv("RABBIT_URL")
	if value == "" {
		value = "amqp://guest:guest@localhost:5672/"
	}
	return value
}
