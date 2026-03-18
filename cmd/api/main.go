package main

import (
	"log/slog"
	"os"

	"github.com/gopher-95/go-subscription-api/internal/config"
	"github.com/gopher-95/go-subscription-api/internal/handlers"
	"github.com/gopher-95/go-subscription-api/internal/repository"
	"github.com/gopher-95/go-subscription-api/internal/server"
	"github.com/gopher-95/go-subscription-api/internal/service"

	_ "github.com/gopher-95/go-subscription-api/docs"
)

// @title Subscription API Service
// @version 1.0.0
// @description REST сервис для управления подписками пользователей

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @tag.name Subscriptions
// @tag.description Операции с подписками пользователей

func main() {
	cfg := config.Load()

	slog.Info("Starting application...")

	slog.Info("Загружаем конфигурацию")

	err := repository.RunMigrations(cfg.ConnectionStringToMigrator())
	if err != nil {
		slog.Error("Migration failed: password is incorrect")
	}
	slog.Info("migrations completed")

	db, err := repository.NewDB(cfg.ConnectionStringToDB())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("database connected")

	storage := repository.NewStorage(db)
	subscriptionService := service.NewService(storage)
	subscriptionHandeler := handlers.NewSubscriptionHandler(subscriptionService)

	router := server.Router(subscriptionHandeler)

	srv := server.NewServer(cfg.SERVER_PORT, router)
	err = srv.Run()
	if err != nil {
		slog.Error("server failed", "error", err)
	}

}
