package sqlstorage

import (
	"context"
	"fmt"
	"time"

	//nolint:gci
	// Postgres Driver.
	_ "github.com/jackc/pgx/stdlib"

	"github.com/jmoiron/sqlx"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

var requests = map[string]string{
	"insert": `INSERT INTO events(event_id, title, start_date, end_date, descr, user_id, notify_user_time)
				VALUES(:event_id, :title, :start_date, :end_date, :descr, :user_id, :notify_user_time);`,
	"update": `UPDATE events SET (event_id, title, start_date, end_date, descr, user_id, notify_user_time)=
				(:event_id, :title, :start_date, :end_date, :descr, :user_id, :notify_user_time)
				WHERE event_id=:event_id;`,
	"delete":     `DELETE FROM events WHERE event_id=:eid;`,
	"select_id":  `SELECT * FROM events WHERE event_id=:eid;`,
	"select_day": `SELECT * FROM events WHERE EXTRACT(DAY FROM start_date) = :start_date;`,
	"select_week": `SELECT * FROM events WHERE EXTRACT(WEEK FROM start_date) = :start_date_week AND
					EXTRACT(YEAR FROM start_date) = :start_date_year;`,
	"select_month": `SELECT * FROM events WHERE EXTRACT(MONTH FROM start_date) = :start_date;`,
	"select_all":   `SELECT * FROM events;`,
	"select_user":  `SELECT * FROM events WHERE user_id = :user_id;`,
}

type Storage struct {
	config        *storage.Storage
	db            *sqlx.DB
	preparedQuery map[string]*sqlx.NamedStmt
}

var _ storage.EventStorage = New(nil)

func New(config *storage.Storage) *Storage {
	return &Storage{
		config:        config,
		preparedQuery: make(map[string]*sqlx.NamedStmt),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) PreparedQueries(ctx context.Context) error {
	for key, val := range requests {
		stmt, err := s.db.PrepareNamed(val) // *sqlx.NamedStmt
		if err != nil {
			return fmt.Errorf("failed to prepare %s query '%v'\n error: %w", key, val, err)
		}
		s.preparedQuery[key] = stmt
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	var err error
	for key, val := range s.preparedQuery {
		err = val.Close()
		if err != nil {
			return fmt.Errorf("failed to close prepared '%s' statement with err: %w", key, err)
		}
	}
	err = s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db with err: %w", err)
	}
	return nil
}

func (s *Storage) Create(e model.Event) error {
	query, ok := s.preparedQuery["insert"]
	if !ok {
		fmt.Printf("prepared query not found")
		return storage.ErrorPreparedQueryNotFound
	}
	_, err := query.Exec(e)
	if err != nil {
		fmt.Printf("failed to insert event %#v to db: error %v\n", e, err)
		return err
	}
	return nil
}

func (s *Storage) Update(eid model.EventID, e model.Event) error {
	e.ID = eid
	queryUpdate, ok := s.preparedQuery["update"]
	if !ok {
		fmt.Printf("prepared query not found")
		return storage.ErrorPreparedQueryNotFound
	}

	_, err := s.GetEventByid(eid)
	if err != nil {
		fmt.Printf("failed to select event by id %s: error %v", eid, err)
		return err
	}

	_, err = queryUpdate.Exec(e)
	if err != nil {
		fmt.Printf("failed to update event %#v to db: error %v", e, err)
		return err
	}
	return nil
}

func (s *Storage) Delete(eid model.EventID) error {
	m := map[string]interface{}{"eid": eid}
	query, ok := s.preparedQuery["delete"]
	if !ok {
		return storage.ErrorPreparedQueryNotFound
	}
	_, err := query.Exec(m)
	return err
}

func (s *Storage) GetEventByid(eid model.EventID) (model.Event, error) {
	m := map[string]interface{}{"eid": eid}
	listEvents, err := s.listEventsByQuery("select_id", m)
	if err != nil {
		return model.Event{}, err
	}
	if len(listEvents) == 0 {
		return model.Event{}, storage.ErrorWrongID
	}
	return listEvents[0], nil
}

func (s *Storage) ListEventsByDay(choosenDay time.Time) ([]model.Event, error) {
	m := map[string]interface{}{"start_date": choosenDay.Day()}
	listEvents, err := s.listEventsByQuery("select_day", m)
	if err != nil {
		return nil, err
	}
	return listEvents, nil
}

func (s *Storage) ListEventsByWeek(choosenWeek time.Time) ([]model.Event, error) {
	year, week := choosenWeek.ISOWeek()
	param := map[string]interface{}{
		"start_date_year": fmt.Sprint(year),
		"start_date_week": fmt.Sprint(week),
	}
	listEvents, err := s.listEventsByQuery("select_week", param)
	if err != nil {
		return nil, err
	}
	return listEvents, nil
}

func (s *Storage) ListEventsByMonth(choosenMonth time.Time) ([]model.Event, error) {
	m := map[string]interface{}{"start_date": choosenMonth.Month()}
	listEvents, err := s.listEventsByQuery("select_month", m)
	if err != nil {
		return nil, err
	}
	return listEvents, nil
}

func (s *Storage) ListAllEvents() ([]model.Event, error) {
	listEvents, err := s.listEventsByQuery("select_all", nil)
	if err != nil {
		return nil, err
	}
	return listEvents, nil
}

func (s *Storage) ListUserEvents(u model.UserID) ([]model.Event, error) {
	m := map[string]interface{}{"user_id": u}
	listEvents, err := s.listEventsByQuery("select_user", m)
	if err != nil {
		return nil, err
	}
	return listEvents, nil
}

func (s *Storage) listEventsByQuery(queryKey string, param interface{}) ([]model.Event, error) {
	listEvents := make([]model.Event, 0) //

	var err error
	var rows *sqlx.Rows
	// var query interface{}
	if param != nil {
		query, ok := s.preparedQuery[queryKey]
		if !ok {
			return nil, storage.ErrorPreparedQueryNotFound
		}

		rows, err = query.Queryx(param)
	} else {
		query := requests[queryKey]
		rows, err = s.db.Queryx(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to select events by '%v' with param %v: error %w", queryKey, param, err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp model.Event
		if err := rows.Scan(&tmp.ID, &tmp.Title, &tmp.StartDate, &tmp.EndDate,
			&tmp.Description, &tmp.UserID, &tmp.NotifyUserTime); err != nil {
			return nil, err
		}
		listEvents = append(listEvents, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listEvents, nil
}

func (s *Storage) ClearAll() error {
	res, err := s.db.Exec("DELETE FROM events;")
	if err != nil {
		fmt.Printf("can't delete contents from events with error: %v\n", err)
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("can't count affected rows: %v\n", err)
		return err
	}
	fmt.Printf("deleted %v rows\n", n)
	return nil
}
