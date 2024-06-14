package env

import (
	"log"
	"os"
	"strconv"
)

func GetDbType() string {
	value := os.Getenv("DB_TYPE")
	if value == "mysql" || value == "postgres" || value == "mssql" {
		return value
	} else {
		log.Fatalf("database type \"%s\" not found", value)
		return ""
	}
}

func GetDbHost() string {
	value := os.Getenv("DB_HOST")
	if value == "" {
		value = "localhost"
	}
	return value
}

func GetDbPort() int {
	value := os.Getenv("DB_PORT")
	if port, err := strconv.Atoi(value); err == nil {
		return port
	}
	return 3306
}

func GetDbUser() string {
	value := os.Getenv("DB_USER")
	if value == "" {
		value = "root"
	}
	return value
}

func GetDbPass() string {
	value := os.Getenv("DB_PASS")
	if value == "" {
		value = ""
	}
	return value
}

func GetDbName() string {
	value := os.Getenv("DB_NAME")
	if value == "" {
		value = "test"
	}
	return value
}

func GetDbMigration() bool {
	return os.Getenv("DB_MIGRATION") == "true"
}
