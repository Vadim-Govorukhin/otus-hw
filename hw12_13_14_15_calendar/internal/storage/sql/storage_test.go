package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/stretchr/testify/require"
)

const devDatabaseURL = "postgres://otus:otus@localhost:5432/calendar?sslmode=disable"
const testDatabaseURL = "postgres://otus:otus@localhost:5432/calendar_test?sslmode=disable"

func teardown(s *sqlstorage.Storage, tables ...string) {
	if len(tables) > 0 {
		// TODO
	}
}

func TestStorage(t *testing.T) {
	t.Run("connect to dev db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		storageTempl := storage.New("sql", devDatabaseURL)
		store := sqlstorage.New(storageTempl)

		err := store.Connect(context.Background())
		require.NoError(t, err)
	})

	t.Run("connect to test db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		storageTempl := storage.New("sql", testDatabaseURL)
		store := sqlstorage.New(storageTempl)

		err := store.Connect(context.Background())
		require.NoError(t, err)
	})
}
