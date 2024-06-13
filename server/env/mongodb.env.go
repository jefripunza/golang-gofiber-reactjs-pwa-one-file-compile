package env

import "os"

func GetMongoUrl() string {
	value := os.Getenv("MONGO_URL")
	if value == "" {
		value = "mongodb://localhost:27017"
	}
	return value
}

func GetMongoName() string {
	value := os.Getenv("MONGO_NAME")
	if value == "" {
		value = "test"
	}
	return value
}
