package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"primbon-app/config"
	"primbon-app/routes"
)

func main() {
	logger := log.New(os.Stderr, "primbon-app ", log.LstdFlags|log.LUTC)
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Fatalf("gagal membaca file .env: %v", err)
	}
	db, err := config.InitDB()
	if err != nil {
		logger.Fatalf("gagal menginisialisasi database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Printf("gagal menutup koneksi database: %v", err)
		}
	}()

	app := routes.SetupRoutes(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatalf("gagal membuka port %s: %v", port, err)
	}
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("server berjalan pada port %s", port)
		serverErrors <- app.Listener(listener)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	select {
	case sig := <-signals:
		logger.Printf("menerima sinyal %s; mematikan server", sig)
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Printf("server berhenti dengan error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Printf("shutdown server tidak selesai dengan tertib: %v", err)
	}
}
