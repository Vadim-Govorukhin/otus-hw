package app

import (
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	storage storage.EventStorage
	log     *logger.Logger
}

func New(logger *logger.Logger, storage storage.EventStorage) *App {
	return &App{log: logger,
		storage: storage}
}

func (a *App) Create(e model.Event) (uuid.UUID, error) {
	a.log.Infof("create event %v", e)
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	err := a.storage.Create(e)
	if err != nil {
		a.log.Errorf("can't create event: %e", err)
	}
	return e.ID, err
}

func (a *App) Update(eid model.EventID, e model.Event) error {
	a.log.Infof("update event with id=%s by event %v", eid, e)
	err := a.storage.Update(eid, e)
	if err != nil {
		a.log.Errorf("can't update event: %e", err)
	}
	return err
}

func (a *App) Delete(eid model.EventID) {
	a.log.Infof("delete event with id=%s", eid)
	a.storage.Delete(eid)
}

func (a *App) ListEventsByDay(date time.Time) ([]model.Event, error) {
	return a.storage.ListEventsByDay(date)
}

func (a *App) ListEventsByWeek(date time.Time) ([]model.Event, error) {
	return a.storage.ListEventsByWeek(date)
}

func (a *App) ListEventsByMonth(date time.Time) ([]model.Event, error) {
	return a.storage.ListEventsByMonth(date)
}

func (a *App) ListAllEvents() ([]model.Event, error) {
	return a.storage.ListAllEvents()
}

func (a *App) ListUserEvents(uid model.UserID) ([]model.Event, error) {
	return a.storage.ListUserEvents(uid)
}
