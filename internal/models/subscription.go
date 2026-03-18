package models

import "time"

// Subscription представляет подписку пользователя
type Subscription struct {
	ID          int        `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserID      string     `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date" db:"end_date"`
}

// UpdateCreateSubscriptionRequest структура для создания и обновления подписки
type UpdateCreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

// TotalCostRequest параметры запроса для подсчета стоимости
type TotalCostRequest struct {
	StartDate   string  `json:"start_date"`   // Начало периода (MM-YYYY)
	EndDate     string  `json:"end_date"`     // Конец периода (MM-YYYY)
	UserID      *string `json:"user_id"`      // Фильтр по пользователю
	ServiceName *string `json:"service_name"` // Фильтр по сервису
}

// TotalCostResponse ответ с суммарной стоимостью
type TotalCostResponse struct {
	TotalCost int    `json:"total_cost"`        // Общая стоимость в рублях
	Period    string `json:"period"`            // Период расчета
	UserID    string `json:"user_id,omitempty"` // ID пользователя (если был фильтр)
	Service   string `json:"service,omitempty"` // Название сервиса (если был фильтр)
}
