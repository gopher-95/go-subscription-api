package migrator

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) error {
	sourceURL := "file://migrations"

	log.Printf("Source URL: %s", sourceURL)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("ошибка создания мигратора: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Миграции не требуются, база данных актуальна")
	} else {
		log.Println("Миграции успешно применены")
	}

	return nil
}
