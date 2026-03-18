package config

import "github.com/joho/godotenv"

type Config struct {
}

func Load() Config {
	_ = godotenv.Load()
	return Config{}
}
