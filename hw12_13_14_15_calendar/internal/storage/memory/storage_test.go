package memorystorage_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	testEvent = storage.Event{
		ID:             storage.EventID("0"),
		Title:          "Test0",
		Date:           time.Date(2022, time.September, 15, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event0",
		UserID:         storage.UserID("user0"),
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
	testEvent2 = storage.Event{
		ID:             storage.EventID("1"),
		Title:          "Test1",
		Date:           time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event1",
		UserID:         storage.UserID("user1"),
		NotifyUserTime: time.Duration(24 * time.Hour),
	}
	testEvent3 = storage.Event{
		ID:             storage.EventID("2"),
		Title:          "Test2",
		Date:           time.Date(2022, time.November, 16, 1, 2, 3, 0, time.UTC),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.UTC),
		Description:    "A Testing Event2",
		UserID:         storage.UserID("user0"),
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
		err := store.Update(storage.EventID("0"), storage.Event{})
		require.ErrorIs(t, err, storage.ErrorWrongUpdateID)

		tmpEvent := storage.Event{ID: storage.EventID("0")}
		err = store.Update(storage.EventID("0"), tmpEvent)
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
		require.True(t, reflect.DeepEqual(storage.ListEvents{tmpEvent, testEvent2}, list))
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
		fmt.Println(reflect.DeepEqual(storage.ListEvents{testEvent2, testEvent3}, list))
		fmt.Println(reflect.DeepEqual(storage.ListEvents{testEvent3, testEvent2}, list))
		require.Equal(t, storage.ListEvents{testEvent2, testEvent3}, list)
		require.True(t, reflect.DeepEqual(storage.ListEvents{testEvent2, testEvent3}, list))

		list = store.ListEventsByWeek(date)
		require.True(t, reflect.DeepEqual(storage.ListEvents{testEvent, testEvent2}, list))

		list = store.ListEventsByMonth(date)
		require.True(t, reflect.DeepEqual(storage.ListEvents{testEvent, testEvent2}, list))

		list = store.ListUserEvents("user0")
		require.True(t, reflect.DeepEqual(storage.ListEvents{testEvent, testEvent3}, list))

	})
}
