package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	EventID = uuid.UUID
	UserID  = int
)

type Event struct {
	ID             EventID   `db:"event_id" json:"event_id,omitempty"`                      // Уникальный идентификатор события
	Title          string    `db:"title" json:"title"`                                      // Заголовок - коротий текст
	StartDate      time.Time `db:"start_date" json:"start_date" time_format:"sql_datetime"` // Дата и время события;
	EndDate        time.Time `db:"end_date" json:"end_date" time_format:"sql_datetime"`     // дата и время окончания события;
	Description    string    `db:"descr" json:"descr,omitempty"`                            // Описание события - длинный текст, опционально
	UserID         UserID    `db:"user_id" json:"user_id"`                                  // ID пользователя, владельца события;
	NotifyUserTime float64   `db:"notify_user_time" json:"notify_user_time,omitempty"`      // За сколько времени в секундах высылать уведомление, опционально
}
