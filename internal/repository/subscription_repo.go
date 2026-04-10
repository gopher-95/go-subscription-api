package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gopher-95/go-subscription-api/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	log.Println("инициализация хранилища подписок")
	return &Storage{db: db}
}

// функция возвращает индентификатор последней добавленной записи
func (storage *Storage) Create(sub *models.Subscription) (int, error) {
	log.Printf("SQL INSERT: service=%s, price=%d, user=%s",
		sub.ServiceName, sub.Price, sub.UserID)

	query := "INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5) RETURNING id"

	var id int
	err := storage.db.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		log.Printf("ошибка создания подписки: %v", err)
		return 0, fmt.Errorf("ошибка добавления записи в бд: %w", err)
	}
	log.Printf("подписка создана: id=%d", id)
	return id, nil
}

// фукнция возвращает информацию о подписке
func (storage *Storage) Get(id int) (*models.Subscription, error) {
	log.Printf("SQL SELECT by id=%d", id)

	sub := &models.Subscription{}

	query := "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1"

	err := storage.db.QueryRow(query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate)

	if err == sql.ErrNoRows {
		log.Printf("подписка id=%d не найден", id)
		return nil, nil
	}
	if err != nil {
		log.Printf("ошибка получения подписки id=%d: %v", id, err)
		return nil, fmt.Errorf("ошибка получения подписки: %w", err)
	}

	log.Printf("подписка id=%d получена", id)

	return sub, nil
}

// функция возвращает количество измененных строк
func (storage *Storage) Update(id int, sub *models.Subscription) (int, error) {
	log.Printf("SQL UPDATE id=%d: service=%s, price=%d", id, sub.ServiceName, sub.Price)

	query := "UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6"

	res, err := storage.db.Exec(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, id)
	if err != nil {
		log.Printf("ошибка обновления id=%d: %v", id, err)
		return 0, fmt.Errorf("ошибка обновления записи в бд: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества удаленных записей: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("подписка id=%d не найдена для обновления", id)
		return 0, sql.ErrNoRows
	}

	log.Printf("подписка id=%d обновлена (затронуто строк: %d)",
		id, rowsAffected)

	return int(rowsAffected), nil
}

func (storage *Storage) Delete(id int) (int, error) {

	log.Printf("SQL DELETE id=%d", id)

	query := "DELETE FROM subscriptions WHERE id = $1"

	res, err := storage.db.Exec(query, id)
	if err != nil {
		log.Printf("ошибка удаления id=%d: %v", id, err)
		return 0, fmt.Errorf("ошибка удаления строки: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества удаленных записей: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("подписка id=%d не найдена для удаления", id)
		return 0, sql.ErrNoRows
	}

	log.Printf("подписка id=%d удалена", id)

	return int(rowsAffected), nil
}

func (storage *Storage) GetAll(limit int, offset int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	log.Printf("SQL SELECT all: limit=%d, offset=%d", limit, offset)

	query := `SELECT id, service_name, price, user_id, start_date, end_date 
	          FROM subscriptions 
			  ORDER BY id
			  LIMIT $1 OFFSET $2`

	rows, err := storage.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("ошибка выполнения запроса: %v", err)
		return nil, fmt.Errorf("ошибка выполнения запроса getall: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var subscription models.Subscription
		var endDateNull sql.NullTime

		err := rows.Scan(&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&endDateNull)
		if err != nil {
			log.Printf("оибка сканирования строки: %v", err)
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		if endDateNull.Valid {
			subscription.EndDate = &endDateNull.Time
		} else {
			subscription.EndDate = nil
		}

		subscriptions = append(subscriptions, subscription)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	log.Printf("получено подписок: %d (limit=%d, offset=%d)",
		len(subscriptions), limit, offset)

	return subscriptions, nil
}

func (storage *Storage) GetSubscriptionsForPeriod(startDate, endDate time.Time, userID, serviceName *string) ([]models.Subscription, error) {
	log.Printf("🔍 SQL SELECT for period: %s - %s",
		startDate.Format("2006-01"), endDate.Format("2006-01"))

	var subscriptions []models.Subscription

	query := `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions
        WHERE start_date <= $1  -- подписка началась не позже конца периода
        AND (end_date IS NULL OR end_date >= $2)  -- подписка не закончилась до начала периода
    `

	args := []interface{}{endDate, startDate}
	argCount := 2

	if userID != nil {
		argCount++
		query += ` AND user_id = $` + fmt.Sprint(argCount)
		args = append(args, *userID)
		log.Printf("фильтр по user_id: %s", *userID)
	}

	if serviceName != nil {
		argCount++
		query += ` AND service_name = $` + fmt.Sprint(argCount)
		args = append(args, *serviceName)
		log.Printf("фильтр по service: %s", *serviceName)
	}

	rows, err := storage.db.Query(query, args...)
	if err != nil {
		log.Printf("ошибка выполнения запроса за период: %v", err)
		return nil, fmt.Errorf("ошибка выполения запроса за период: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription
		var endDateNull sql.NullTime

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&endDateNull,
		)

		if err != nil {
			log.Printf("ошибка сканирования строки: %v", err)
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		if endDateNull.Valid {
			sub.EndDate = &endDateNull.Time
		}

		subscriptions = append(subscriptions, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	log.Printf("получено подписок за период: %d", len(subscriptions))

	return subscriptions, nil
}
