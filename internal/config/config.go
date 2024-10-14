// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

// LoadConfig загружает конфигурацию из файла .env или переменных окружения
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "referral_db"),
		JWTSecret:  getEnv("JWT_SECRET", "your_jwt_secret"),
	}
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
