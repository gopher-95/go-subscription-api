package main

import (
	"log"

	"github.com/gopher-95/go-subscription-api/internal/config"
	"github.com/gopher-95/go-subscription-api/internal/migrator"
	"github.com/gopher-95/go-subscription-api/internal/repository"
)

func main() {
	cfg := config.Load()

	err := migrator.RunMigrations(cfg.ConnectionStringToMigrator())
	if err != nil {
		log.Fatal("не удалось создать миграции для бд: ", err)
	}

	log.Println("Миграции выполнены")

	db, err := repository.NewDB(cfg.ConnectionStringToDB())
	if err != nil {
		log.Println("не удалось запустить бд")
	}
	defer db.Close()

	log.Println("База данных готова к работе")

}
