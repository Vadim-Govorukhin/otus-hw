package sqlstorage_test

import (
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("connect to db", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())

		var databaseURL string
		db, err := sqlx.Open("pgx", databaseURL)
		require.NoError(t, err)

		err = db.Ping()
		require.NoError(t, err)
	})
}
