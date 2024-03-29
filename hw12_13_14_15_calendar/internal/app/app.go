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
	return &App{
		log:     logger,
		storage: storage,
	}
}

func (a *App) Create(e model.Event) (uuid.UUID, error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	a.log.Infof("create event %v", e)
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
		a.log.Errorf("can't update event: %s", err)
	}
	return err
}

func (a *App) Delete(eid model.EventID) error {
	a.log.Infof("delete event with id=%s", eid)
	err := a.storage.Delete(eid)
	if err != nil {
		a.log.Errorf("can't delete event: %s", err)
	}
	return err
}

func (a *App) GetEventByid(eid model.EventID) (model.Event, error) {
	a.log.Infof("get event by id=%s", eid)
	event, err := a.storage.GetEventByid(eid)
	if err != nil {
		a.log.Errorf("can't get event: %s", err)
	}
	return event, err
}

func (a *App) ListEventsByDay(date time.Time) ([]model.Event, error) {
	a.log.Infof("select events by day %v", date.Day())
	listEvents, err := a.storage.ListEventsByDay(date)
	a.log.Infof("selected %v events", len(listEvents))
	if err != nil {
		a.log.Errorf("can't select events: %s", err)
	}
	return listEvents, err
}

func (a *App) ListEventsByWeek(date time.Time) ([]model.Event, error) {
	a.log.Info("select events by week")
	listEvents, err := a.storage.ListEventsByWeek(date)
	a.log.Infof("selected %v events", len(listEvents))
	if err != nil {
		a.log.Errorf("can't select events: %s", err)
	}
	return listEvents, err
}

func (a *App) ListEventsByMonth(date time.Time) ([]model.Event, error) {
	a.log.Infof("select events by month %v", date.Month())
	listEvents, err := a.storage.ListEventsByMonth(date)
	a.log.Infof("selected %v events", len(listEvents))
	if err != nil {
		a.log.Errorf("can't select events: %s", err)
	}
	return listEvents, err
}

func (a *App) ListAllEvents() ([]model.Event, error) {
	a.log.Info("select all events")
	listEvents, err := a.storage.ListAllEvents()
	a.log.Infof("selected %v events", len(listEvents))
	if err != nil {
		a.log.Errorf("can't select events: %s", err)
	}
	return listEvents, err
}

func (a *App) ListUserEvents(uid model.UserID) ([]model.Event, error) {
	a.log.Infof("select events by user %v", uid)
	listEvents, err := a.storage.ListUserEvents(uid)
	a.log.Infof("selected %v events", len(listEvents))
	if err != nil {
		a.log.Errorf("can't select events: %s", err)
	}
	return listEvents, err
}
