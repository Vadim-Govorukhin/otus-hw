package memorystorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

var _ storage.EventStorage = New()

type Storage struct {
	mu sync.RWMutex
	db map[storage.EventID]storage.Event
}

func New() *Storage {
	return &Storage{db: make(map[storage.EventID]storage.Event, 0)}
}

func (s *Storage) Create(e storage.Event) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("Добавляем событие с id=%s\n", e.ID)
	if _, ok := s.db[e.ID]; ok {
		fmt.Printf("Не удалось добавить событие с ID=%s\n", e.ID)
		return storage.ErrorEventIDBusy
	}
	s.db[e.ID] = e
	return
}

func (s *Storage) Update(eid storage.EventID, e storage.Event) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if eid != e.ID {
		fmt.Printf("Нельзя обновить событие с ID %s на событие с ID %s\n", eid, e.ID)
		return storage.ErrorWrongUpdateID
	}
	s.db[eid] = e
	return
}

func (s *Storage) Delete(eid storage.EventID) {
	s.mu.Lock()
	delete(s.db, eid)
	s.mu.Unlock()
}

func (s *Storage) ListEventsByDay(choosenDay time.Time) storage.ListEvents {
	listEvents := make(storage.ListEvents, 0) //
	day := choosenDay.Day()
	s.mu.RLock()
	for _, val := range s.db {
		if val.Date.Day() == day {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents
}

func (s *Storage) ListEventsByWeek(choosenWeek time.Time) storage.ListEvents {
	listEvents := make(storage.ListEvents, 0) //
	year, week := choosenWeek.ISOWeek()
	s.mu.RLock()
	var vWeek, vYear int
	for _, val := range s.db {
		vYear, vWeek = val.Date.ISOWeek()
		if (vWeek == week) && (vYear == year) {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents
}

func (s *Storage) ListEventsByMonth(choosenMonth time.Time) storage.ListEvents {
	listEvents := make(storage.ListEvents, 0) //
	month := choosenMonth.Month()
	s.mu.RLock()
	for _, val := range s.db {
		if val.Date.Month() == month {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents
}

func (s *Storage) ListAllEvents() storage.ListEvents {
	listEvents := make(storage.ListEvents, len(s.db))
	i := 0
	for _, val := range s.db {
		listEvents[i] = val
		i++
	}
	return listEvents
}

func (s *Storage) ListUserEvents(userID storage.UserID) storage.ListEvents {
	listEvents := make(storage.ListEvents, 0) //
	s.mu.RLock()
	for _, val := range s.db {
		if val.UserID == userID {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents
}
