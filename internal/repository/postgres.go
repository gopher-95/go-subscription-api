package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewDB(connectionString string) (*sql.DB, error) {
	log.Println("подключение к PostgreSQL...")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("ошибка открытия соединения: %v", err)
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	err = db.Ping()
	if err != nil {
		log.Printf("ошибка подключения к БД: %v", err)
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	log.Println("Подключение к БД успешно!")

	return db, nil

}
