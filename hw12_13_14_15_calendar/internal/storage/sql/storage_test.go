package sqlstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" // Postgres Driver.
	"github.com/stretchr/testify/require"
)

const (
	devDatabaseURL  = "postgres://otus:otus@localhost:5432/calendar?sslmode=disable"
	testDatabaseURL = "postgres://otus:otus@localhost:5432/calendar_test?sslmode=disable"
)

func teardown(s *Storage, tables []string) error {
	if len(tables) > 0 {
		// res, err := s.db.Exec("DELETE FROM ", strings.Join(tables, " ,"))
		res, err := s.db.Exec("DELETE FROM events;")
		if err != nil {
			fmt.Printf("can't delete contents from %s\n with error: %v\n", tables, err)
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			fmt.Printf("can't count affected rows: %v\n", err)
			return err
		}
		fmt.Printf("deleted %v rows\n", n)
	}
	return nil
}

func setupTest(t *testing.T) *Storage {
	t.Helper()
	fmt.Printf("====== start test %s =====\n", t.Name())

	ctx := context.Background()

	storageTempl := &storage.Storage{Type: "sql", DatabaseURL: testDatabaseURL}
	store := New(storageTempl)

	err := store.Connect(ctx)
	require.NoError(t, err)

	err = store.PreparedQueries(ctx)
	require.NoError(t, err)

	return store
}

func TestStorage(t *testing.T) {
	t.Run("connect to dev db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())

		storageTempl := &storage.Storage{Type: "sql", DatabaseURL: devDatabaseURL}
		store := New(storageTempl)
		defer store.Close(context.Background())

		err := store.Connect(context.Background())
		require.NoError(t, err)
	})

	t.Run("connect to test db", func(t *testing.T) {
		store := setupTest(t)
		defer store.Close(context.Background())
	})

	t.Run("create events in test db", func(t *testing.T) {
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
			store.Close(context.Background())
		}()

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)
	})

	t.Run("get by id test db", func(t *testing.T) {
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
			store.Close(context.Background())
		}()

		_, err := store.GetEventByid(storage.TestEvent.ID)
		require.ErrorIs(t, err, storage.ErrorWrongID)

		err = store.Create(storage.TestEvent)
		require.NoError(t, err)

		e, err := store.GetEventByid(storage.TestEvent.ID)
		require.NoError(t, err)
		require.Equal(t, storage.TestEvent, e)
	})

	t.Run("Update and delete event", func(t *testing.T) {
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
			store.Close(context.Background())
		}()

		store.Create(storage.TestEvent)
		tmpEvent := storage.TestEvent3
		err := store.Update(storage.TestEvent.ID, tmpEvent)
		require.NoError(t, err)

		tmpEvent.ID = storage.TestEvent.ID
		err = store.Update(storage.TestEvent.ID, tmpEvent)
		require.NoError(t, err)
		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, []model.Event{tmpEvent}, list)

		err = store.Create(storage.TestEvent)
		require.Error(t, err)
		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, []model.Event{tmpEvent}, list)

		err = store.Create(storage.TestEvent2)
		require.NoError(t, err)
		store.Delete(storage.TestEvent.ID)
		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.ElementsMatch(t, list, []model.Event{storage.TestEvent2})
	})

	t.Run("lists events in test db", func(t *testing.T) {
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
			store.Close(context.Background())
		}()

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)

		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, list, []model.Event{storage.TestEvent})
	})

	t.Run("lists events in test db", func(t *testing.T) {
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
			store.Close(context.Background())
		}()

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)
		err = store.Create(storage.TestEvent2)
		require.NoError(t, err)
		err = store.Create(storage.TestEvent3)
		require.NoError(t, err)

		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.ElementsMatch(t, list,
			[]model.Event{storage.TestEvent, storage.TestEvent2, storage.TestEvent3})

		date := time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC)

		list, err = store.ListEventsByDay(date)
		require.NoError(t, err)
		require.ElementsMatch(t, list, []model.Event{storage.TestEvent2, storage.TestEvent3})

		list, err = store.ListEventsByWeek(date)
		require.NoError(t, err)
		require.ElementsMatch(t, list, []model.Event{storage.TestEvent2, storage.TestEvent})

		list, err = store.ListEventsByMonth(date)
		require.NoError(t, err)
		require.ElementsMatch(t, list, []model.Event{storage.TestEvent2, storage.TestEvent})

		list, err = store.ListUserEvents(0)
		require.NoError(t, err)
		require.ElementsMatch(t, list, []model.Event{storage.TestEvent3, storage.TestEvent})
	})
}
