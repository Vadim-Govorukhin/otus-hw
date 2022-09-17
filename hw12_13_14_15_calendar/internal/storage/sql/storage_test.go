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

//const databaseURL = "postgres://username:password@localhost:5432/database_name"
const databaseURL = "postgres://otus:otus@localhost:5432/calendar?sslmode=disable"

func TestStorage(t *testing.T) {
	t.Run("connect to db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		storageTempl := storage.New("sql", databaseURL)
		store := sqlstorage.New(storageTempl)

		err := store.Connect(context.Background())
		require.NoError(t, err)
	})
}
