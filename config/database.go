package config

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

func InitDB() *sql.DB {
    host := "localhost"
    port := 5432
    user := "postgres"
    password := "admin" // 🔁 ganti sesuai password PostgreSQL kamu
    dbname := "primbon_db"

    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("❌ Gagal konek ke database:", err)
    }

    log.Println("✅ Database connected")
    return db
}
