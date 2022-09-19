package model

import (
	"time"

	"github.com/google/uuid"
)

type EventID = uuid.UUID
type UserID = int

type Event struct {
	ID             EventID       `db:"event_id"`         // Уникальный идентификатор события
	Title          string        `db:"title"`            // Заголовок - коротий текст
	StartDate      time.Time     `db:"start_date"`       // Дата и время события;
	EndDate        time.Time     `db:"end_date"`         // дата и время окончания события;
	Description    string        `db:"descr"`            // Описание события - длинный текст, опционально
	UserID         UserID        `db:"user_id"`          // ID пользователя, владельца события;
	NotifyUserTime time.Duration `db:"notify_user_time"` // За сколько времени высылать уведомление, опционально
}
