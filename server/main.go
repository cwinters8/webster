package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cwinters8/webster"

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
			"Title":   "Webster",
			"Content": "Hello, World!",
		})
	})

	// make output dir if it doesn't exist
	if err := os.MkdirAll("tmp", 0755); err != nil {
		return fmt.Errorf("failed to create tmp directory: %w", err)
	}

	// takes submitted content from editor and saves to a file
	// TODO: this should be a specific route for writing to a temp file
	app.Post("/content", func(c *fiber.Ctx) error {
		var content webster.Content
		if err := c.BodyParser(&content); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "failed to parse text from editor:", err.Error())
		}
		if err := content.WriteTemp(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(201)
	})

	return app.Listen(fmt.Sprintf(":%d", port))
}

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
}
