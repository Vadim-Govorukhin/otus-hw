package sqlstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/stretchr/testify/require"
)

const devDatabaseURL = "postgres://otus:otus@localhost:5432/calendar?sslmode=disable"
const testDatabaseURL = "postgres://otus:otus@localhost:5432/calendar_test?sslmode=disable"

func teardown(s *Storage, tables ...string) {
	if len(tables) > 0 {
		// TODO
	}
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
}
