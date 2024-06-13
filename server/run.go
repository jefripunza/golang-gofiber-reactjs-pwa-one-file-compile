package server

import (
	"backend/server/connection"
	"backend/server/http"
	"backend/server/util"
	"backend/server/variable"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func cronStart() {

}

func Run(embeddedFiles embed.FS) {
	var err error

	Env := util.Env{}
	Env.Load()
	err = Env.SetTimezone()
	if err != nil {
		log.Fatalf("error on set timezone: %s", err.Error())
		return
	}

	Dir := util.Dir{}
	Dir.Make(variable.TempPath)

	// ---------------------------------

	MongoDB := connection.MongoDB{}
	MongoDB.Connect()

	RabbitMQ := connection.RabbitMQ{}
	RabbitMQ.Connect()

	// ---------------------------------

	go func() {
		cronStart()
	}()

	go func() {
		http.Server(embeddedFiles)
	}()

	// ---------------------------------

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	log.Println("🚦 Listen to Ctrl+C ...")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
