package memorystorage

import (
	"sync"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu sync.RWMutex
	db map[storage.EventID]storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Create(e storage.Event) {
	s.mu.Lock()
	s.db[e.ID] = e
	s.mu.Unlock()
}

func (s *Storage) Update(eid storage.EventID, e storage.Event) {
	s.mu.Lock()
	s.db[eid] = e
	s.mu.Unlock()
}

func (s *Storage) Delete(eid storage.EventID) {
	s.mu.Lock()
	delete(s.db, eid)
	s.mu.Unlock()
}

func (s *Storage) ListEventsByDay(time.Time) {
	s.mu.RLock()
	// TODO
	s.mu.RUnlock()
}

func (s *Storage) ListEventsByWeek(time.Time) {
	s.mu.RLock()
	// TODO
	s.mu.RUnlock()
}

func (s *Storage) ListEventsByMonth(time.Time) {
	s.mu.RLock()
	// TODO
	s.mu.RUnlock()
}
