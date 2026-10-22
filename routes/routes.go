package routes

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"primbon-app/controllers"
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

	app.Get("/", controllers.Home)
	app.Post("/ramal", func(c *fiber.Ctx) error {
		return controllers.RamalJodoh(c, db)
	})
	app.Get("/hari-baik", controllers.HariBaik)
	app.Get("/health", func(c *fiber.Ctx) error {
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
