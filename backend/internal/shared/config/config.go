package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment        string `env:"ENVIRONMENT"`
	SqliteDatabasePath string `env:"SQLITE_DATABASE_PATH"`
	MysqlDsn           string `env:"MYSQL_DSN"`

	GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT" envDefault:"5s"`
	HttpAddress     string        `env:"HTTP_ADDRESS"`
	AdminUsername   string        `env:"ADMIN_USERNAME"`
	AdminPassword   string        `env:"ADMIN_PASSWORD"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	return env.ParseAs[Config]()
}
