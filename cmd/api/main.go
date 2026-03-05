package main

import (
	"log"

	"github.com/gopher-95/go-subscription-api/internal/config"
	"github.com/gopher-95/go-subscription-api/internal/repository"
)

func main() {
	_, err := repository.NewDB(config.Load())
	if err != nil {
		log.Fatal("ошибка подключения к бд: ", err)
	}

	log.Println("Контейнер успешно запущен!")

}
