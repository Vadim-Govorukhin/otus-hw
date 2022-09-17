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

var (
	testEvent = model.Event{
		ID:             model.EventID("0"),
		Title:          "Test0",
		StartDate:      time.Date(2022, time.September, 15, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event0",
		UserID:         model.UserID("user0"),
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
	testEvent2 = model.Event{
		ID:             model.EventID("1"),
		Title:          "Test1",
		StartDate:      time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event1",
		UserID:         model.UserID("user1"),
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
	testEvent3 = model.Event{
		ID:             model.EventID("2"),
		Title:          "Test2",
		StartDate:      time.Date(2022, time.November, 16, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event2",
		UserID:         model.UserID("user0"),
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
)

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
		err := store.Update(model.EventID("0"), model.Event{})
		require.ErrorIs(t, err, storage.ErrorWrongUpdateID)

		tmpEvent := model.Event{ID: model.EventID("0")}
		err = store.Update(model.EventID("0"), tmpEvent)
		require.NoError(t, err)
		list := store.ListAllEvents()
		require.Equal(t, storage.ListEvents{tmpEvent}, list)

		err = store.Create(testEvent)
		require.ErrorIs(t, err, storage.ErrorEventIDBusy)
		list = store.ListAllEvents()
		require.Equal(t, storage.ListEvents{tmpEvent}, list)

		err = store.Create(testEvent2)
		require.NoError(t, err)
		list = store.ListAllEvents()
		require.True(t, list.Equal(storage.ListEvents{tmpEvent, testEvent2}))
	})

	t.Run("check lists of events", func(t *testing.T) {
		fmt.Printf("====== start test %s =====\n", t.Name())
		store := memorystorage.New()

		err := store.Create(testEvent)
		require.NoError(t, err)
		err = store.Create(testEvent2)
		require.NoError(t, err)
		err = store.Create(testEvent3)
		require.NoError(t, err)
		date := time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC)

		list := store.ListEventsByDay(date)
		require.True(t, list.Equal(storage.ListEvents{testEvent2, testEvent3}))

		list = store.ListEventsByWeek(date)
		require.True(t, list.Equal(storage.ListEvents{testEvent2, testEvent}))

		list = store.ListEventsByMonth(date)
		require.True(t, list.Equal(storage.ListEvents{testEvent2, testEvent}))

		list = store.ListUserEvents("user0")
		require.True(t, list.Equal(storage.ListEvents{testEvent3, testEvent}))
	})
}
