package main

import (
	"log"
	"primbon-app/config"
	"primbon-app/routes"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("❌ Gagal menginisialisasi database:", err)
	}
	defer db.Close()

	app := routes.SetupRoutes(db)

	log.Println("🚀 Server running at http://localhost:8080")
	app.Listen(":8080")
}
