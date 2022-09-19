package sqlstorage

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/stretchr/testify/require"
)

const (
	devDatabaseURL  = "postgres://otus:otus@localhost:5432/calendar?sslmode=disable"
	testDatabaseURL = "postgres://otus:otus@localhost:5432/calendar_test?sslmode=disable"
)

var (
	testEvent = model.Event{
		ID:             uuid.New(),
		Title:          "Test0",
		StartDate:      time.Date(2022, time.September, 15, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event0",
		UserID:         0,
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
	testEvent2 = model.Event{
		ID:             uuid.New(),
		Title:          "Test1",
		StartDate:      time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event1",
		UserID:         1,
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
	testEvent3 = model.Event{
		ID:             uuid.New(),
		Title:          "Test2",
		StartDate:      time.Date(2022, time.November, 16, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event2",
		UserID:         0,
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
)

func teardown(s *Storage, tables []string) error {
	if len(tables) > 0 {
		res, err := s.db.Exec("DELETE FROM ", strings.Join(tables, " ,"))
		if err != nil {
			fmt.Printf("can't delete contents from %s\n with error: %v\n", tables, err)
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			fmt.Printf("can't count affected rows: %v\n", err)
			return err
		}
		fmt.Printf("deleted %v rows", n)
	}
	return nil
}

func TestStorage(t *testing.T) {
	t.Run("connect to dev db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())

		storageTempl := &storage.Storage{Type: "sql", DatabaseURL: devDatabaseURL}
		store := New(storageTempl)

		err := store.Connect(context.Background())
		require.NoError(t, err)
	})

	t.Run("connect to test db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		storageTempl := &storage.Storage{Type: "sql", DatabaseURL: testDatabaseURL}
		store := New(storageTempl)

		err := store.Connect(context.Background())
		require.NoError(t, err)
	})

	t.Run("prepare queries", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		storageTempl := &storage.Storage{Type: "sql", DatabaseURL: testDatabaseURL}
		store := New(storageTempl)

		ctx := context.Background()
		err := store.Connect(ctx)
		require.NoError(t, err)

		err = store.PreparedQueries(ctx)
		require.NoError(t, err)
	})

	t.Run("create events in test db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		storageTempl := &storage.Storage{Type: "sql", DatabaseURL: testDatabaseURL}
		store := New(storageTempl)

		ctx := context.Background()
		err := store.Connect(ctx)
		require.NoError(t, err)

		err = store.PreparedQueries(ctx)
		require.NoError(t, err)

		err = store.Create(testEvent)
		require.NoError(t, err)
		err = store.Create(testEvent2)
		require.NoError(t, err)
		err = store.Create(testEvent3)
		require.NoError(t, err)

		err = teardown(store, []string{"events"})
		require.NoError(t, err)
	})
}
