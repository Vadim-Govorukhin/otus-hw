package storage

import (
	"errors"
	"time"
)

var (
	ErrorDayBusy = errors.New("данное время уже занято другим событием")
	// TODO
)

type EventID string

type Event struct {
	ID             EventID       // Уникальный идентификатор события
	Title          string        // Заголовок - коротий текст
	Date           time.Time     // Дата и время события;
	EndDate        time.Time     // дата и время окончания события;
	Description    string        // Описание события - длинный текст, опционально
	UserID         string        // ID пользователя, владельца события;
	NotifyUserTime time.Duration // За сколько времени высылать уведомление, опционально
}

type EventStorage interface {
	Create(Event)                // добавление события в хранилище;
	Update(EventID, Event)       // изменение события в хранилище;
	Delete(EventID)              // удаление события из хранилища;
	ListEventsByDay(time.Time)   // листинг событий на день
	ListEventsByWeek(time.Time)  // листинг событий на неделю
	ListEventsByMonth(time.Time) // листинг событий на день
}
