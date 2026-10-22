package config

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// InitDB opens and verifies a PostgreSQL connection using required environment
// variables: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, and DB_SSLMODE.
func InitDB() (*sql.DB, error) {
	config, err := databaseConfigFromEnv(os.Getenv)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.connectionString())
	if err != nil {
		return nil, fmt.Errorf("membuka koneksi database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("menguji koneksi database: %w", err)
	}

	return db, nil
}

type databaseConfig struct {
	host     string
	port     int
	user     string
	password string
	database string
	sslmode  string
}

func databaseConfigFromEnv(getenv func(string) string) (databaseConfig, error) {
	var config databaseConfig
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
		if strings.TrimSpace(getenv(key)) == "" {
			return databaseConfig{}, fmt.Errorf("environment variable %s wajib diisi", key)
		}
	}

	config.host = strings.TrimSpace(getenv("DB_HOST"))
	config.user = getenv("DB_USER")
	config.password = getenv("DB_PASSWORD")
	config.database = getenv("DB_NAME")
	config.sslmode = strings.TrimSpace(getenv("DB_SSLMODE"))

	port, err := strconv.Atoi(strings.TrimSpace(getenv("DB_PORT")))
	if err != nil || port < 1 || port > 65535 {
		return databaseConfig{}, fmt.Errorf("DB_PORT harus berupa angka antara 1 dan 65535")
	}
	config.port = port

	switch config.sslmode {
	case "disable", "require", "verify-ca", "verify-full", "allow", "prefer":
	default:
		return databaseConfig{}, fmt.Errorf("DB_SSLMODE tidak valid: %q", config.sslmode)
	}

	return config, nil
}

func (config databaseConfig) connectionString() string {
	values := url.Values{}
	values.Set("host", config.host)
	values.Set("port", strconv.Itoa(config.port))
	values.Set("user", config.user)
	values.Set("password", config.password)
	values.Set("dbname", config.database)
	values.Set("sslmode", config.sslmode)
	return values.Encode()
}
