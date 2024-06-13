package http

import (
	"backend/server/env"
	"embed"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Server(embeddedFiles embed.FS) *fiber.App {
	var err error

	port := env.GetServerPort()
	server_name := env.GetServerName()

	app := fiber.New(fiber.Config{
		ServerHeader:          server_name,
		DisableStartupMessage: true,
		CaseSensitive:         true,
		BodyLimit:             10 * 1024 * 1024, // 10 MB / max file size
	})

	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept,X-Browser-ID",
	}))
	app.Use(requestid.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	staticFiles, _ := fs.Sub(embeddedFiles, "dist")
	fileServer := http.FS(staticFiles)
	app.Use("/", func(c *fiber.Ctx) error {
		if c.Path() == "/" || !strings.HasPrefix(c.Path(), "/api") {
			reqPath := c.Path()
			if reqPath == "/" {
				reqPath = "/index.html"
			}

			file, err := fileServer.Open(reqPath)
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			// Determine the content type
			ext := filepath.Ext(reqPath)
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			c.Type(ext)
			c.Response().Header.Set("Content-Type", mimeType)

			return c.SendStream(file, int(fileInfo.Size()))
		}
		return c.Next()
	})

	app.Use(logger.New())

	Module(app)
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "endpoint not found!",
		})
	})

	log.Printf("✅ Server \"%s\" started on port http://localhost:%s\n", server_name, port)
	if err = app.Listen("127.0.0.1:" + port); err != nil {
		log.Fatalln("error start server:", err)
	}

	return app
}
