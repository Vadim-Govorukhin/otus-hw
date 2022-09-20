package memorystorage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var storageTempl = &storage.Storage{Type: "memory", DatabaseURL: ""}

func TestStorage(t *testing.T) {
	t.Run("create event", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New(storageTempl)

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)

		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, storage.ListEvents{storage.TestEvent}, list)
	})

	t.Run("Update and delete event", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New(storageTempl)

		store.Create(storage.TestEvent)
		err := store.Update(storage.TestEvent.ID, model.Event{})
		require.ErrorIs(t, err, storage.ErrorWrongUpdateID)

		tmpEvent := model.Event{ID: storage.TestEvent.ID}
		err = store.Update(storage.TestEvent.ID, tmpEvent)
		require.NoError(t, err)
		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, storage.ListEvents{tmpEvent}, list)

		err = store.Create(storage.TestEvent)
		require.ErrorIs(t, err, storage.ErrorEventIDBusy)
		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, storage.ListEvents{tmpEvent}, list)

		err = store.Create(storage.TestEvent2)
		require.NoError(t, err)
		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.ElementsMatch(t, list, storage.ListEvents{tmpEvent, storage.TestEvent2})
	})

	t.Run("check lists of events", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New(storageTempl)

		err := store.Create(storage.TestEvent)
		require.NoError(t, err)
		err = store.Create(storage.TestEvent2)
		require.NoError(t, err)
		err = store.Create(storage.TestEvent3)
		require.NoError(t, err)
		date := time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC)

		list, err := store.ListEventsByDay(date)
		require.NoError(t, err)
		require.ElementsMatch(t, list, storage.ListEvents{storage.TestEvent2, storage.TestEvent3})

		list, err = store.ListEventsByWeek(date)
		require.NoError(t, err)
		require.ElementsMatch(t, list, storage.ListEvents{storage.TestEvent2, storage.TestEvent})

		list, err = store.ListEventsByMonth(date)
		require.NoError(t, err)
		require.ElementsMatch(t, list, storage.ListEvents{storage.TestEvent2, storage.TestEvent})

		list, err = store.ListUserEvents(0)
		require.NoError(t, err)
		require.ElementsMatch(t, list, storage.ListEvents{storage.TestEvent3, storage.TestEvent})
	})
}
