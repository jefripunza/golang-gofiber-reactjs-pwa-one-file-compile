package main

import (
	"backend/cron"
	"backend/server"
	"backend/server/connection"
	"backend/server/util"
	"backend/server/variable"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed dist/*
var embeddedFiles embed.FS

func main() {
	// var err error

	// ---------------------------------

	Env := util.Env{}
	Env.Load()

	// err = Env.SetTimezone()
	// if err != nil {
	// 	log.Fatalf("error on set timezone: %s", err.Error())
	// 	return
	// }

	Dir := util.Dir{}
	Dir.Make(variable.TempPath)

	// ---------------------------------

	MongoDB := connection.MongoDB{}
	MongoDB.Connect()

	SQL := connection.SQL{}
	SQL.Connect()

	RabbitMQ := connection.RabbitMQ{}
	RabbitMQ.Connect()

	Redis := connection.Redis{}
	Redis.Connect()

	// ---------------------------------

	cron.Start()
	server.Run(embeddedFiles)

	// ---------------------------------

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	time.Sleep(3 * time.Second)
	log.Println("🚦 Listen to Ctrl+C ...")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
