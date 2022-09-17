package storage

import (
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
)

type EventStorage interface {
	Create(model.Event) error                // Добавление события в хранилище
	Update(model.EventID, model.Event) error // Изменение события в хранилище
	Delete(model.EventID)                    // Удаление события из хранилища
	ListEventsByDay(time.Time) ListEvents    // Листинг событий на день
	ListEventsByWeek(time.Time) ListEvents   // Листинг событий на неделю
	ListEventsByMonth(time.Time) ListEvents  // Листинг событий на день
	ListAllEvents() ListEvents               // Листинг всех событий
	ListUserEvents(model.UserID) ListEvents  // Листинг всех событий юзера
}

type Notification struct {
	EventID model.EventID // ID события
	Title   string        //Заголовок события
	Date    time.Time     // Дата события
	User    model.UserID  // Пользователь, которому отправлять
}

type ListEvents []model.Event

func (l *ListEvents) Equal(tgt ListEvents) bool {
	if len(*l) != len(tgt) {
		return false
	}

	lMap := make(map[model.EventID]model.Event, len(*l))
	for _, item := range *l {
		lMap[item.ID] = item
	}

	tgtMap := make(map[model.EventID]model.Event, len(*l))
	for _, item := range tgt {
		tgtMap[item.ID] = item
	}

	for key, val := range lMap {
		tgtVal, ok := tgtMap[key]
		if !ok || tgtVal != val {
			return false
		}
	}
	return true

}
