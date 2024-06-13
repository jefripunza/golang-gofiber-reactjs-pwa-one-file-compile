package main

import (
	"backend/server"
	"embed"
)

//go:embed dist/*
var embeddedFiles embed.FS

func main() {
	server.Run(embeddedFiles)
}
