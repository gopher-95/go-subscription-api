package repository

import (
	"database/sql"
	"fmt"

	"github.com/gopher-95/go-subscription-api/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) GetById(id int64) (*domain.Subscription, error) {
	sub := &domain.Subscription{}

	query := "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1"

	err := r.db.QueryRow(query, id).Scan(
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

func (r *Repository) DeleteById(id int64) error {
	query := "DELETE FROM subscriptions WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления строки: %w", err)
	}

	number, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных записей: %w", err)
	}

	if number == 0 {
		return sql.ErrNoRows
	}

	return nil
}
