package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"primbon-app/controllers"
)

func SetupRoutes(db *sql.DB) *fiber.App {
	app := fiber.New()

	app.Get("/", controllers.Home)
	app.Post("/ramal", func(c *fiber.Ctx) error {
		return controllers.RamalJodoh(c, db)
	})
	app.Get("/hari-baik", controllers.HariBaik)

	return app
}
