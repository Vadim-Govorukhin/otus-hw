package memorystorage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
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

var storageTempl = &storage.Storage{Type: "memory", DatabaseURL: ""}

func TestStorage(t *testing.T) {
	t.Run("create event", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New(storageTempl)

		err := store.Create(testEvent)
		require.NoError(t, err)

		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, storage.ListEvents{testEvent}, list)
	})

	t.Run("Update and delete event", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New(storageTempl)

		store.Create(testEvent)
		err := store.Update(testEvent.ID, model.Event{})
		require.ErrorIs(t, err, storage.ErrorWrongUpdateID)

		tmpEvent := model.Event{ID: testEvent.ID}
		err = store.Update(testEvent.ID, tmpEvent)
		require.NoError(t, err)
		list, err := store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, storage.ListEvents{tmpEvent}, list)

		err = store.Create(testEvent)
		require.ErrorIs(t, err, storage.ErrorEventIDBusy)
		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.Equal(t, storage.ListEvents{tmpEvent}, list)

		err = store.Create(testEvent2)
		require.NoError(t, err)
		list, err = store.ListAllEvents()
		require.NoError(t, err)
		require.True(t, list.Equal(storage.ListEvents{tmpEvent, testEvent2}))
	})

	t.Run("check lists of events", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New(storageTempl)

		err := store.Create(testEvent)
		require.NoError(t, err)
		err = store.Create(testEvent2)
		require.NoError(t, err)
		err = store.Create(testEvent3)
		require.NoError(t, err)
		date := time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC)

		list, err := store.ListEventsByDay(date)
		require.NoError(t, err)
		require.True(t, list.Equal(storage.ListEvents{testEvent2, testEvent3}))

		list, err = store.ListEventsByWeek(date)
		require.NoError(t, err)
		require.True(t, list.Equal(storage.ListEvents{testEvent2, testEvent}))

		list, err = store.ListEventsByMonth(date)
		require.NoError(t, err)
		require.True(t, list.Equal(storage.ListEvents{testEvent2, testEvent}))

		list, err = store.ListUserEvents(0)
		require.NoError(t, err)
		require.True(t, list.Equal(storage.ListEvents{testEvent3, testEvent}))
	})
}
