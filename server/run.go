package server

import (
	"backend/server/http"
	"embed"
)

func Run(embeddedFiles embed.FS) {
	go func() {
		http.Server(embeddedFiles)
	}()

}
