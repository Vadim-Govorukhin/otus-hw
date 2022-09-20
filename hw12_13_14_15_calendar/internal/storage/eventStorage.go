package storage

import (
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
)

type EventStorage interface {
	Create(model.Event) error                        // Добавление события в хранилище
	Update(model.EventID, model.Event) error         // Изменение события в хранилище
	Delete(model.EventID)                            // Удаление события из хранилища
	ListEventsByDay(time.Time) (ListEvents, error)   // Листинг событий на день
	ListEventsByWeek(time.Time) (ListEvents, error)  // Листинг событий на неделю
	ListEventsByMonth(time.Time) (ListEvents, error) // Листинг событий на день
	ListAllEvents() (ListEvents, error)              // Листинг всех событий
	ListUserEvents(model.UserID) (ListEvents, error) // Листинг всех событий юзера
}

type Notification struct {
	EventID model.EventID // ID события
	Title   string        //Заголовок события
	Date    time.Time     // Дата события
	User    model.UserID  // Пользователь, которому отправлять
}

type ListEvents []model.Event
