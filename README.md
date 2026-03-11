# Subscription API Service

REST сервис для управления подписками. Тестовое задание Effective Mobile.

## Стек
Go 1.24, PostgreSQL 15, Docker, Chi, go-migrate

Миграции БД применяются автоматически при запуске.

## Запуск
```bash
git clone https://github.com/gopher-95/go-subscription-api
cd go-subscription-api
docker-compose up --build

Сервис: http://localhost:8080

## Методы API
API
Метод	URL	                            Описание
POST	/api/v1/subscriptions	        создать
GET	    /api/v1/subscriptions	        список
GET	    /api/v1/subscriptions/{id}	    получить
PUT	    /api/v1/subscriptions/{id}	    обновить
DELETE	/api/v1/subscriptions/{id}	    удалить

### Пример запроса
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 399,
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "start_date": "03-2026"
  }'

## Конфигурация
cp .env.example .env
отредактировать .env

Автор:
gopher-95
https://github.com/gopher-95