package storage

import (
	"time"
)

type EventID string
type UserID string

type ListEvents []Event

type Event struct {
	ID             EventID       // Уникальный идентификатор события
	Title          string        // Заголовок - коротий текст
	Date           time.Time     // Дата и время события;
	EndDate        time.Time     // дата и время окончания события;
	Description    string        // Описание события - длинный текст, опционально
	UserID         UserID        // ID пользователя, владельца события;
	NotifyUserTime time.Duration // За сколько времени высылать уведомление, опционально
}

type EventStorage interface {
	Create(Event) error                     // Добавление события в хранилище
	Update(EventID, Event) error            // Изменение события в хранилище
	Delete(EventID)                         // Удаление события из хранилища
	ListEventsByDay(time.Time) ListEvents   // Листинг событий на день
	ListEventsByWeek(time.Time) ListEvents  // Листинг событий на неделю
	ListEventsByMonth(time.Time) ListEvents // Листинг событий на день
	ListAllEvents() ListEvents              // Листинг всех событий
	ListUserEvents(UserID) ListEvents       // Листинг всех событий юзера
}

type Notification struct {
	EventID EventID   // ID события
	Title   string    //Заголовок события
	Date    time.Time // Дата события
	User    UserID    // Пользователь, которому отправлять
}
