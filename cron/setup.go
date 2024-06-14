package cron

import (
	cronjob "backend/cron/job"

	"github.com/robfig/cron/v3"
)

func Start() {
	go func() {
		cronJob := cron.New()

		cronJob.AddFunc("@every 5s", cronjob.RevokeTokenExpired)

		cronJob.Start()
	}()
}
