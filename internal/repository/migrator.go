package repository

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) error {
	log.Println("запуск миграций базы данных...")
	log.Printf("папка с миграциями: file://migrations")

	sourceURL := "file://migrations"

	log.Printf("Source URL: %s", sourceURL)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Printf("ошибка создания мигратора: %v", err)
		return fmt.Errorf("ошибка создания мигратора: %w", err)
	}
	defer m.Close()

	log.Println("мигратор инициализирован, применяем миграции...")

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("ошибка применения миграций: %v", err)
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Миграции не требуются, база данных актуальна")
	} else {
		log.Println("Миграции успешно применены")
	}

	version, dirty, err := m.Version()
	if err == nil {
		log.Printf("текущая версия миграций: %d (грязная: %v)", version, dirty)
	}

	return nil
}
