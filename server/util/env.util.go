package util

import (
	"backend/server/env"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Env struct{}

func (ref Env) Load() {
	pwd := env.GetPwd()
	envFilePath := filepath.Join(pwd, ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Println("file .env not found")
	}
}

func (ref Env) SetTimezone() error {
	timezone := "Asia/Jakarta"
	os.Setenv("TZ", timezone)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return err
	}
	time.Local = loc
	return nil
}
