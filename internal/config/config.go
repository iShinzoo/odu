package config

import "os"

type Config struct {
	DBUrl string
}

func LoadConfig() *Config {
	return &Config{
		DBUrl: getEnv("DATABASE_URL", "postgres://krsna:krsna123@127.0.0.1:5432/orderdb?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
