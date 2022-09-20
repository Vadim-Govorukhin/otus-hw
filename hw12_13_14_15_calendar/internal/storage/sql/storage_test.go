package sqlstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/stretchr/testify/require"
)

const (
	devDatabaseURL  = "postgres://otus:otus@localhost:5432/calendar?sslmode=disable"
	testDatabaseURL = "postgres://otus:otus@localhost:5432/calendar_test?sslmode=disable"
)

func teardown(s *Storage, tables []string) error {
	if len(tables) > 0 {
		//res, err := s.db.Exec("DELETE FROM ", strings.Join(tables, " ,"))
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

	storageTempl := &storage.Storage{Type: "sql", DatabaseURL: testDatabaseURL}
	store := New(storageTempl)

	ctx := context.Background()
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
		_ = setupTest(t)
	})

	t.Run("create events in test db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
		}()

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)
	})

	t.Run("lists events in test db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := setupTest(t)
		defer func() {
			err := teardown(store, []string{"events"})
			require.NoError(t, err)
		}()

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)

		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, list, storage.ListEvents{storage.TestEvent})

		err = store.Create(storage.TestEvent2)
		require.NoError(t, err)
		err = store.Create(storage.TestEvent3)
		require.NoError(t, err)

		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.ElementsMatch(t, list,
			storage.ListEvents{storage.TestEvent, storage.TestEvent2, storage.TestEvent3})

	})
}
