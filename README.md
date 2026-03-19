# Subscription API Service

REST сервис для управления подписками пользователей.  

## 🛠️Стек технологий
- Go 1.25
- PostgreSQL 15
- Docker / Docker Compose
- Chi Router
- Go-migrate
- Swagger

## 🚀Быстрый старт
1. Клонировать репозиторй 
```bash
git clone https://github.com/gopher-95/go-subscription-api
```
2. Перейти в проект 
```bash
cd go-subscription-api
```
3. Собрать образ
```bash
docker-compose up --build
```
4. После успешного запуска сервис будет доступен:
   - **API**: `http://localhost:8080`
   - **Swagger документация**: `http://localhost:8080/swagger/index.html`

## 🌐API Endpoints

| Метод | Endpoint | Описание | Параметры |
|------|---------|----------|----------|
| POST | `/api/v1/subscriptions` | Создать новую подписку | body: JSON с данными подписки |
| GET | `/api/v1/subscriptions` | Получить список подписок | `limit`, `offset` (query params) |
| GET | `/api/v1/subscriptions/{id}` | Получить подписку по ID | `id` (path param) |
| PUT | `/api/v1/subscriptions/{id}` | Обновить существующую подписку | `id` (path), body (JSON) |
| DELETE | `/api/v1/subscriptions/{id}` | Удалить подписку | `id` (path param) |
| GET | `/api/v1/subscriptions/total-cost` | Рассчитать стоимость за период | `start_date`, `end_date`, `user_id`, `service_name` |

## 📝Примеры запросов

### Создание подписки (POST)

**Запрос:**
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 399,
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "start_date": "03-2026"
  }'
```

**Ответ:**
```json
{
    "message": "подписка успешно создана",
    "id": 1
}
```

---

### Получение подписки (GET)

**По ID:**
```bash
curl http://localhost:8080/api/v1/subscriptions/1
```

**Ответ:**
```json
{
    "message": "подписка успешно получена",
    "subscription": {
        "id": 1,
        "service_name": "Yandex Plus",
        "price": 399,
        "user_id": "123e4567-e89b-12d3-a456-426614174000",
        "start_date": "2026-03-01T00:00:00Z",
        "end_date": null
    }
}
```

**Список всех подписок (с пагинацией):**
```bash
curl "http://localhost:8080/api/v1/subscriptions?limit=5&offset=0"
```

---

### Обновление подписки (PUT)

```bash
curl -X PUT http://localhost:8080/api/v1/subscriptions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus Premium",
    "price": 599,
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "start_date": "03-2026",
    "end_date": "03-2027"
  }'
```

**Ответ:**
```json
{
    "message": "подписка успешно обновлена",
    "id": 1,
    "rows_affected": 1
}
```

---

### Удаление подписки (DELETE)

```bash
curl -X DELETE http://localhost:8080/api/v1/subscriptions/1
```

**Ответ:**
```json
{
    "message": "подписка успешно удалена",
    "deleted_count": 1
}
```

---

### Подсчет стоимости за период (GET)

**Без фильтров:**
```bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2026&end_date=12-2026"
```

**С фильтром по пользователю:**
```bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2026&end_date=12-2026&user_id=123e4567-e89b-12d3-a456-426614174000"
```

**Ответ:**
```json
{
    "total_cost": 4794,
    "period": "01-2026 - 12-2026",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

## 🗄️Миграции базы данных

Миграции применяются автоматически при запуске контейнера.

```bash
# файлы миграций находятся в папке
migrations/
├── 000001_create_subscriptions_table.up.sql
├── 000001_create_subscriptions_table.down.sql
├── 000002_add_indexes.up.sql
└── 000002_add_indexes.down.sql
```
## 📊 Логирование

В проекте настроено структурированное логирование на всех уровнях приложения:

| Компонент | События |
|-----------|---------|
| **Config** | загрузка конфигурации, чтение переменных окружения |
| **Repository** | SQL запросы, ошибки БД, время выполнения операций |
| **Service** | валидация данных, бизнес-логика, расчет стоимости |
| **Handlers** | входящие HTTP запросы, исходящие ответы, ошибки |
| **Main** | запуск сервера, миграции, подключение к БД |

## 📄Переменные окружения

Файл env.example
```env
# база данных 
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscription_db
DB_SSLMODE=disable

# сервер
SERVER_PORT=8080
```
### Настройка
1. Скопировать пример конфигурации
```bash
cp .env.example .env
```
2. Отредактировать файл под свои параметры
```bash
nano .env
```

## 📚 Swagger документация
1. После запуска сервера документация доступна по адресу:
```text
http://localhost:8080/swagger/index.html
```
2. Что задокументировано

| Endpoint                          | Method  | Description              |
|----------------------------------|--------|--------------------------|
| `/api/v1/subscriptions`          | POST   | Create subscription      |
| `/api/v1/subscriptions`          | GET    | Get subscriptions list   |
| `/api/v1/subscriptions/{id}`     | GET    | Get subscription by ID   |
| `/api/v1/subscriptions/{id}`     | PUT    | Update subscription      |
| `/api/v1/subscriptions/{id}`     | DELETE | Delete subscription      |
| `/api/v1/subscriptions/total-cost` | GET  | Get total cost for period |

3. Генерация
```bash
swag init -g cmd/api/main.go
```
4. Файлы 
```text
docs/
├── swagger.json
└── swagger.yaml
```

## 👤 Автор

🐙 **GitHub:** [gopher-95](https://github.com/gopher-95)
📧 **Email:** gopher-95@yandex.ru




