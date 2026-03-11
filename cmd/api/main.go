package main

import (
	"log"

	"github.com/gopher-95/go-subscription-api/internal/config"
	"github.com/gopher-95/go-subscription-api/internal/handlers"
	"github.com/gopher-95/go-subscription-api/internal/repository"
	"github.com/gopher-95/go-subscription-api/internal/server"
	"github.com/gopher-95/go-subscription-api/internal/service"
)

func main() {
	cfg := config.Load()

	err := repository.RunMigrations(cfg.ConnectionStringToMigrator())
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

	storage := repository.NewStorage(db)
	subscriptionService := service.NewService(storage)
	subscriptionHandeler := handlers.NewSubscriptionHandler(subscriptionService)

	router := server.Router(subscriptionHandeler)

	srv := server.NewServer(cfg.SERVER_PORT, router)
	err = srv.Run()
	if err != nil {
		log.Fatal("ошибка запуска сервера:", err)
	}

}
