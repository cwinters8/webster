package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

const port = 8888

//go:embed views
var viewsFS embed.FS

func setup() error {
	fs := http.FS(viewsFS)
	engine := html.NewFileSystem(fs, ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(compress.New(), logger.New(), favicon.New(favicon.Config{
		FileSystem: fs,
		File:       "views/assets/images/favicon.ico",
	}))
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       fs,
		PathPrefix: "views/assets",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/index", fiber.Map{
			"Title":   "Craftly",
			"Content": "Hello, World!",
		})
	})
	return app.Listen(fmt.Sprintf(":%d", port))
}

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
}
