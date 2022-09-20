package storage

import (
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
