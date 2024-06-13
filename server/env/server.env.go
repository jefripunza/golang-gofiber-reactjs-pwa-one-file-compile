package env

import (
	"os"
)

func GetServerName() string {
	value := os.Getenv("SERVER_NAME")
	if value == "" {
		value = "P34C3_KHYREIN"
	}
	return value
}

func GetServerPort() string {
	value := os.Getenv("SERVER_PORT")
	if value == "" {
		value = "3003"
	}
	return value
}

func GetSecretKey() string {
	value := os.Getenv("SECRET_KEY")
	if value == "" {
		value = "your_secret_key"
	}
	return value
}
