package memorystorage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var testEvent = storage.Event{
	ID:             storage.EventID("0"),
	Title:          "Test",
	Date:           time.Date(2022, time.September, 15, 1, 2, 3, 0, time.UTC),
	EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
	Description:    "A Testing Event",
	UserID:         "user0",
	NotifyUserTime: time.Duration(24 * time.Hour),
}

func TestStorage(t *testing.T) {
	t.Run("create event", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New()

		err := store.Create(testEvent)
		require.NoError(t, err)

		list := store.ListAllEvents()
		require.Equal(t, storage.ListEvents{testEvent}, list)
	})

	t.Run("Update and delete event", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New()

		store.Create(testEvent)
		store.Update(storage.EventID("0"), storage.Event{})
		list := store.ListAllEvents()
		require.Equal(t, storage.ListEvents{storage.Event{}}, list)

		err := store.Create(testEvent)
		require.ErrorIs(t, err, storage.ErrorEventIDBusy)
		list = store.ListAllEvents()
		require.Equal(t, storage.ListEvents{storage.Event{}}, list)

		testEvent.ID = storage.EventID("1")
		err = store.Create(testEvent)
		require.NoError(t, err)
		list = store.ListAllEvents()
		require.Equal(t, storage.ListEvents{storage.Event{}, testEvent}, list)
	})
}
