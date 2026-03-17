package config

import (
	"fmt"
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
	DB_SSLMODE  string
	SERVER_PORT string
}

// функция загрузки env файла
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("не найден .env файл, используем значения по умолчанию")
	} else {
		log.Println("конфигурация загружена из env файла")
	}

	config := &Config{
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5433"),
		DB_USER:     getEnv("DB_USER", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),
		DB_NAME:     getEnv("DB_NAME", "subscription_db"),
		DB_SSLMODE:  getEnv("SSL_MODE", "disable"),
		SERVER_PORT: getEnv("SERVER_PORT", "8080"),
	}

	if config.DB_PASSWORD == "" {
		log.Println("не указан пароль к базе данных")
	}

	log.Printf("конфигурация загружена: порт=%s, БД=%s:%s/%s", config.SERVER_PORT, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	return config
}

// функция чтения env файла
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		if key != "DB_PASSWORD" {
			log.Printf("переменная %s не задана, используем значение по умолчанию: %s", key, defaultValue)
		}
		return defaultValue
	}

	if key != "DB_PASSWORD" {
		log.Printf("переменная %s=%s (из окружения)", key, value)
	}

	return value
}

func (cfg *Config) ConnectionStringToDB() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_SSLMODE)
}

func (cfg *Config) ConnectionStringToMigrator() string {

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME, cfg.DB_SSLMODE)
}
