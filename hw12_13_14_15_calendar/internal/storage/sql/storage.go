package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	db          *sqlx.DB
	databaseURL string
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.databaseURL)

	if err != nil {
		fmt.Printf("failed to load driver: %v", err)
		return storage.ErrorLoadDriver
	}

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("failed to connect to db: %v", err)
		return storage.ErrorConnectDB
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Create(storage.Event) error {
	//query := `insert into events(owner, title, descr, start_date, end_date, )
	//		values($1, $2, $3, $4, $5)`

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
