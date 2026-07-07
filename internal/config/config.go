package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvDev  Environment = "dev"
	EnvProd Environment = "prod"
)

func (e Environment) IsDev() bool { return e == EnvDev }

type Config struct {
	Environment        Environment `env:"ENVIRONMENT" envDefault:"dev"`
	SqliteDatabasePath string      `env:"SQLITE_DATABASE_PATH" envDefault:"./store.db"`

	HostAddress       string        `env:"HOST_ADDRESS" envDefault:"1.2.3.4"`
	GracefulTimeout   time.Duration `env:"GRACEFUL_TIMEOUT" envDefault:"5s"`
	HttpAddress       string        `env:"HTTP_ADDRESS" envDefault:":8080"`
	AdminUsername     string        `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword     string        `env:"ADMIN_PASSWORD" envDefault:"password"`
	AdminPasswordHash string        `env:"ADMIN_PASSWORD_HASH"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	return env.ParseAs[Config]()
}
