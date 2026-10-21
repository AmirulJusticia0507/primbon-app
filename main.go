package main

import (
	"log"
	"primbon-app/config"
	"primbon-app/routes"
)

func main() {
	db := config.InitDB()
	app := routes.SetupRoutes(db)

	log.Println("🚀 Server running at http://localhost:8080")
	app.Listen(":8080")
}
