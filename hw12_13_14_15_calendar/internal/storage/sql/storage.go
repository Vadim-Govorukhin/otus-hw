package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

var requests = map[string]string{
	"insert": `INSERT INTO events(event_id, title, start_date, end_date, descr, user_id, notify_user_time)
				VALUES(:event_id, :title, :start_date, :end_date, :descr, :user_id, :notify_user_time);`,
	"select_day": `SELECT * FROM events WHERE EXTRACT(DAY FROM start_date) = :start_date`,
}

type Storage struct { // TODO
	config        *storage.Storage
	db            *sqlx.DB
	preparedQuery map[string]*sqlx.NamedStmt
}

func New(config *storage.Storage) *Storage {
	fmt.Println("Create SQL Storage")
	return &Storage{
		config:        config,
		preparedQuery: make(map[string]*sqlx.NamedStmt)}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.config.DatabaseURL)

	if err != nil {
		fmt.Printf("failed to load driver: %v\n", err)
		return ErrorLoadDriver
	}
	fmt.Println(db.Stats())

	if err := db.Ping(); err != nil {
		fmt.Printf("failed to connect to db: %v\n", err)
		return ErrorConnectDB
	}
	s.db = db
	return nil
}

func (s *Storage) PreparedQueries(ctx context.Context) error {
	// создаем подготовленный запрос
	for key, val := range requests {
		stmt, err := s.db.PrepareNamed(val) // *sqlx.NamedStmt
		if err != nil {
			fmt.Printf("failed to prepare %s query '%v'\n error: %v", key, val, err)
			return err
		}
		s.preparedQuery[key] = stmt
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) {
	var err error
	for key, val := range s.preparedQuery {
		err = val.Close()
		if err != nil {
			fmt.Printf("failed to close prepared '%s' statement with err: %v", key, err)
		}
	}
	err = s.db.Close()
	if err != nil {
		fmt.Printf("failed to close db with err: %v", err)
	}
}

func (s *Storage) Create(e model.Event) error {
	query, ok := s.preparedQuery["insert"]
	if !ok {
		fmt.Printf("prepared query not found")
		return storage.ErrorPreparedQueryNotFound
	}
	_, err := query.Exec(e)
	if err != nil {
		fmt.Printf("failed to insert event %#v to db: error %v", e, err)
		return err
	}
	return nil
}

func (s *Storage) Update(eid model.EventID, e model.Event) error {
	// TODO
	return nil
}

func (s *Storage) Delete(eid model.EventID) {
	// TODO
}

func (s *Storage) ListEventsByDay(choosenDay time.Time) storage.ListEvents {
	listEvents := make(storage.ListEvents, 0) //
	rows, err := s.preparedQuery["select_day"].Queryx(choosenDay.Day())
	if err != nil {
		fmt.Printf("failed to select events by day %v: error %v", choosenDay.Day(), err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var tmp model.Event
		if err := rows.Scan(&tmp); err != nil {
			// ошибка сканирования
			return nil
		}
		// обрабатываем строку
		listEvents = append(listEvents, tmp)
	}
	if err := rows.Err(); err != nil {
		// ошибка при получении результатов
		return nil
	}

	return listEvents
}

func (s *Storage) ListEventsByWeek(choosenWeek time.Time) storage.ListEvents {
	// TODO
	return nil
}

func (s *Storage) ListEventsByMonth(choosenMonth time.Time) storage.ListEvents {
	// TODO
	return nil
}

func (s *Storage) ListAllEvents() storage.ListEvents {
	// TODO
	return nil
}

func (s *Storage) ListUserEvents(u model.UserID) storage.ListEvents {
	// TODO
	return nil
}
