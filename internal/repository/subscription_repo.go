package repository

import (
	"database/sql"
	"fmt"

	"github.com/gopher-95/go-subscription-api/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// функция возвращает индентификатор последней добавленной записи
func (storage *Storage) Create(sub *models.Subscription) (int, error) {
	query := "INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5) RETURNING id"

	var id int

	err := storage.db.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления записи в бд: %w", err)
	}

	return id, nil
}

// фукнция возвращает информацию о подписке
func (storage *Storage) Get(id int) (*models.Subscription, error) {
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
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения подписки: %w", err)
	}

	return sub, nil
}

// функция возвраoает количество измененных строк
func (storage *Storage) Update(id int, sub *models.Subscription) (int, error) {
	query := "UPDATE subscriptions SET service_name = $1, price = $2, start_date = $3, end_date = $4 WHERE id = $5"

	res, err := storage.db.Exec(query, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, id)
	if err != nil {
		return 0, fmt.Errorf("ошибка обновления записи в бд: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества удаленных записей: %w", err)
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return int(rowsAffected), nil
}

func (storage *Storage) Delete(id int) (int, error) {
	query := "DELETE FROM subscriptions WHERE id = $1"

	res, err := storage.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("ошибка удаления строки: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества удаленных записей: %w", err)
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return int(rowsAffected), nil
}

func (storage *Storage) GetAll(limit int, offset int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	query := `SELECT id, service_name, price, user_id, start_date, end_date 
	          FROM subscriptions 
			  ORDER BY id
			  LIMIT $1 OFFSET $2`

	rows, err := storage.db.Query(query, limit, offset)
	if err != nil {
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

	return subscriptions, nil
}
