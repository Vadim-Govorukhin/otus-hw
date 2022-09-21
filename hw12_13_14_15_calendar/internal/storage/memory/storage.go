package memorystorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

var _ storage.EventStorage = New(nil)

type Storage struct {
	config *storage.Storage
	mu     sync.RWMutex
	db     map[model.EventID]model.Event
}

func New(config *storage.Storage) *Storage {
	fmt.Println("Create memory Storage")
	return &Storage{
		config: config,
		db:     make(map[model.EventID]model.Event, 0),
	}
}

func (s *Storage) Create(e model.Event) (err error) {
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

func (s *Storage) Update(eid model.EventID, e model.Event) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if eid != e.ID {
		fmt.Printf("Нельзя обновить событие с ID %s на событие с ID %s\n", eid, e.ID)
		return storage.ErrorWrongUpdateID
	}
	s.db[eid] = e
	return
}

func (s *Storage) Delete(eid model.EventID) {
	s.mu.Lock()
	delete(s.db, eid)
	s.mu.Unlock()
}

func (s *Storage) ListEventsByDay(choosenDay time.Time) ([]model.Event, error) {
	listEvents := make([]model.Event, 0) //
	day := choosenDay.Day()
	s.mu.RLock()
	for _, val := range s.db {
		if val.StartDate.Day() == day {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents, nil
}

func (s *Storage) ListEventsByWeek(choosenWeek time.Time) ([]model.Event, error) {
	listEvents := make([]model.Event, 0) //
	year, week := choosenWeek.ISOWeek()
	s.mu.RLock()
	var vWeek, vYear int
	for _, val := range s.db {
		vYear, vWeek = val.StartDate.ISOWeek()
		if (vWeek == week) && (vYear == year) {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents, nil
}

func (s *Storage) ListEventsByMonth(choosenMonth time.Time) ([]model.Event, error) {
	listEvents := make([]model.Event, 0) //
	month := choosenMonth.Month()
	s.mu.RLock()
	for _, val := range s.db {
		if val.StartDate.Month() == month {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents, nil
}

func (s *Storage) ListAllEvents() ([]model.Event, error) {
	listEvents := make([]model.Event, len(s.db))
	i := 0
	for _, val := range s.db {
		listEvents[i] = val
		i++
	}
	return listEvents, nil
}

func (s *Storage) ListUserEvents(userID model.UserID) ([]model.Event, error) {
	listEvents := make([]model.Event, 0) //
	s.mu.RLock()
	for _, val := range s.db {
		if val.UserID == userID {
			listEvents = append(listEvents, val)
		}
	}
	s.mu.RUnlock()
	return listEvents, nil
}
