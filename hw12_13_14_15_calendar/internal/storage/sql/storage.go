package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Create(storage.Event) error {
	return nil
}

func (s *Storage) Update(storage.EventID, storage.Event) error {
	return nil
}

func (s *Storage) Delete(storage.EventID) {

}

func (s *Storage) ListEventsByDay(time.Time) storage.ListEvents {
	return nil
}

func (s *Storage) ListEventsByWeek(time.Time) storage.ListEvents {
	return nil
}

func (s *Storage) ListEventsByMonth(time.Time) storage.ListEvents {
	return nil
}

func (s *Storage) ListAllEvents() storage.ListEvents {
	return nil
}

func (s *Storage) ListUserEvents(storage.UserID) storage.ListEvents {
	return nil
}
