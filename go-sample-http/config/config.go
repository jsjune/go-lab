package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	GRPCPort string
	DBPath   string
	GinMode  string
}

func Load() *Config {
	// .env 파일이 있으면 로드 (없어도 무시)
	godotenv.Load()

	return &Config{
		Port:     getEnv("APP_PORT", "8080"),
		GRPCPort: getEnv("APP_GRPC_PORT", "9090"),
		DBPath:   getEnv("APP_DB_PATH", "./app.db"),
		GinMode:  getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
