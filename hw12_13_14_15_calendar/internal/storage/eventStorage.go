package storage

import (
	"context"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
)

type EventStorage interface {
	Create(model.Event) error                           // Добавление события в хранилище
	Update(model.EventID, model.Event) error            // Изменение события в хранилище
	Delete(model.EventID)                               // Удаление события из хранилища
	GetEventByid(model.EventID) (model.Event, error)    // Получить событие по его id
	ListEventsByDay(time.Time) ([]model.Event, error)   // Листинг событий на день
	ListEventsByWeek(time.Time) ([]model.Event, error)  // Листинг событий на неделю
	ListEventsByMonth(time.Time) ([]model.Event, error) // Листинг событий на месяц
	ListAllEvents() ([]model.Event, error)              // Листинг всех событий
	ListUserEvents(model.UserID) ([]model.Event, error) // Листинг всех событий юзера
	Close(context.Context) error                        // Закрытие хранилища
}

type Notification struct {
	EventID model.EventID // ID события
	Title   string        // Заголовок события
	Date    time.Time     // Дата события
	User    model.UserID  // Пользователь, которому отправлять
}
