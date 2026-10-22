package routes

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"primbon-app/controllers"
	"primbon-app/internal/assets"
)

func SetupRoutes(db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Primbon App",
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if fiberErr, ok := err.(*fiber.Error); ok {
				code = fiberErr.Code
			}
			if code >= fiber.StatusInternalServerError {
				log.Printf("permintaan %s %s gagal: %v", c.Method(), c.Path(), err)
				return c.Status(code).SendString("Terjadi kesalahan pada server. Silakan coba lagi nanti.")
			}
			return c.Status(code).SendString("Permintaan tidak dapat diproses.")
		},
	})

	staticFiles, err := fs.Sub(assets.Files, "static")
	if err != nil {
		log.Printf("gagal memuat static assets: %v", err)
	} else {
		app.Get("/static/*", func(c *fiber.Ctx) error {
			filePath := strings.TrimPrefix(c.Path(), "/static/")
			data, err := fs.ReadFile(staticFiles, filePath)
			if err != nil {
				return fiber.ErrNotFound
			}
			if strings.HasSuffix(filePath, ".svg") {
				c.Set("Content-Type", "image/svg+xml")
			}
			return c.Send(data)
		})
	}
	app.Get("/", controllers.Home)
	app.Post("/ramal", func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("Kalkulator sementara tidak tersedia karena database belum dikonfigurasi.")
		}
		return controllers.RamalJodoh(c, db)
	})
	app.Get("/hari-baik", controllers.HariBaik)
	app.Post("/hari-baik", controllers.HariBaik)
	app.Post("/api/ai/chat", func(c *fiber.Ctx) error {
		return controllers.AIChat(c, os.Getenv("AI_API_KEY"), os.Getenv("AI_MODEL"))
	})
	app.Get("/favicon.svg", func(c *fiber.Ctx) error {
		data, err := fs.ReadFile(assets.Files, "static/favicon.svg")
		if err != nil {
			return fiber.ErrNotFound
		}
		c.Set("Content-Type", "image/svg+xml")
		c.Set("Cache-Control", "public, max-age=86400")
		return c.Send(data)
	})
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.Redirect("/favicon.svg?v=2", fiber.StatusTemporaryRedirect)
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "database belum dikonfigurasi"})
		}
		if err := db.PingContext(c.UserContext()); err != nil {
			log.Printf("health check database gagal: %v", err)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "tidak sehat"})
		}
		return c.JSON(fiber.Map{"status": "sehat"})
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Halaman tidak ditemukan.")
	})

	return app
}
