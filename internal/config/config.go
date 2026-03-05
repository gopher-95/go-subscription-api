package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	SSL_MODE    string
	SERVER_PORT string
}

// функция загрузки env файла
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("не найден .env файл, используем значения по умолчанию")
	}

	config := &Config{
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5432"),
		DB_USER:     getEnv("DB_USER", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),
		DB_NAME:     getEnv("DB_NAME", "subscription_db"),
		SSL_MODE:    getEnv("SSL_MODE", "disable"),
		SERVER_PORT: getEnv("SERVER_PORT", "8080"),
	}

	if config.DB_PASSWORD == "" {
		log.Println("не указан пароль к базе данных")
	}

	return config
}

// функция чтения env файла
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
