package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	SqliteDatabasePath string
	MysqlDsn           string

	GracefulTimeout time.Duration
	HttpAddress     string
	AdminUsername   string
	AdminPassword   string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	gracefulTimeout, err := strconv.Atoi(os.Getenv("GRACEFUL_TIMEOUT"))
	if err != nil {
		gracefulTimeout = 5
	}

	return Config{
		SqliteDatabasePath: os.Getenv("SQLITE_DATABASE_PATH"),
		MysqlDsn:           os.Getenv("MYSQL_DSN"),

		GracefulTimeout: time.Duration(gracefulTimeout) * time.Second,
		HttpAddress:     os.Getenv("HTTP_ADDRESS"),
		AdminUsername:   os.Getenv("ADMIN_USERNAME"),
		AdminPassword:   os.Getenv("ADMIN_PASSWORD"),
	}

}
