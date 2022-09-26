package grpc

import (
	"context"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/api/stubs/eventer"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	eventer.UnimplementedCalendarServer
}

func (s Server) CreateEvent(ctx context.Context, e *eventer.Event) (*eventer.EventID, error) {
	panic("not implemented") // TODO: Implement
}

func (s Server) UprateEvent(ctx context.Context, _ *eventer.UpdateEventRequest) (*eventer.EventID, error) {
	panic("not implemented") // TODO: Implement
}

func (s Server) DeleteEvent(ctx context.Context, eid *eventer.EventID) (*eventer.EventID, error) {
	panic("not implemented") // TODO: Implement
}

func (s Server) GetEventByID(ctx context.Context, eid *eventer.EventID) (*eventer.Event, error) {
	panic("not implemented") // TODO: Implement
}

func (s Server) ListEventByDay(ctx context.Context, date *timestamppb.Timestamp) (*eventer.EventResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s Server) ListEventByMonth(ctx context.Context, date *timestamppb.Timestamp) (*eventer.EventResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s Server) ListAllEvent(ctx context.Context, _ *emptypb.Empty) (*eventer.EventResponse, error) {
	var err error
	e := &eventer.EventResponse{}
	return e, err
}

func (s Server) ListAllEventByUser(ctx context.Context, uid *eventer.UserID) (*eventer.EventResponse, error) {
	panic("not implemented") // TODO: Implement
}

func NewServer(logg *logger.Logger, app *app.App, conf *config.HTTPServerConf) eventer.CalendarServer {
	return Server{}
}
