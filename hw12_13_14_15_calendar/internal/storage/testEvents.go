package storage

import (
	"encoding/json"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

var (
	TestEvent = model.Event{
		ID:             uuid.New(),
		Title:          "Test0",
		StartDate:      time.Date(2022, time.September, 15, 1, 2, 3, 0, time.Local),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.Local),
		Description:    "A Testing Event0",
		UserID:         0,
		NotifyUserTime: (24 * time.Hour).Seconds(),
	}
	TestEvent2 = model.Event{
		ID:             uuid.New(),
		Title:          "Test1",
		StartDate:      time.Date(2022, time.September, 16, 1, 2, 3, 0, time.Local),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.Local),
		Description:    "A Testing Event1",
		UserID:         1,
		NotifyUserTime: (24 * time.Hour).Seconds(),
	}
	TestEvent3 = model.Event{
		ID:             uuid.New(),
		Title:          "Test2",
		StartDate:      time.Date(2022, time.November, 16, 1, 2, 3, 0, time.Local),
		EndDate:        time.Date(2022, time.December, 15, 1, 2, 3, 0, time.Local),
		Description:    "A Testing Event2",
		UserID:         0,
		NotifyUserTime: (24 * time.Hour).Seconds(),
	}
)

var (
	TestEventJson  []byte
	TestEvent2Json []byte
	TestEvent3Json []byte
)

type TestEventIDRespose struct {
	ID model.EventID `db:"event_id" json:"event_id,omitempty"`
}

var (
	TestEventIDJson  []byte
	TestEvent3IDJson []byte
	TestEvent2IDJson []byte
)

func init() {
	TestEventJson, _ = json.Marshal(TestEvent)
	TestEvent2Json, _ = json.Marshal(TestEvent2)
	TestEvent3Json, _ = json.Marshal(TestEvent3)

	TestEventID := TestEventIDRespose{
		ID: TestEvent.ID,
	}
	TestEvent2ID := TestEventIDRespose{
		ID: TestEvent2.ID,
	}
	TestEvent3ID := TestEventIDRespose{
		ID: TestEvent3.ID,
	}

	TestEventIDJson, _ = json.Marshal(TestEventID)
	TestEvent2IDJson, _ = json.Marshal(TestEvent2ID)
	TestEvent3IDJson, _ = json.Marshal(TestEvent3ID)
}
