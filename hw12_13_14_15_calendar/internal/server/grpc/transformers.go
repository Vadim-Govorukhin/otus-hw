package internalgrpc

import (
	"fmt"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/api/stubs/eventer"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventToGRPC(e *model.Event) *eventer.Event {
	ev := &eventer.Event{
		Id:             EventIDToGRPC(&e.ID),
		Title:          e.Title,
		StartDate:      timestamppb.New(e.StartDate),
		EndDate:        timestamppb.New(e.EndDate),
		Description:    &e.Description,
		UserId:         UserIDToGRPC(&e.UserID),
		NotifyUserTime: &e.NotifyUserTime,
	}
	return ev
}

func EventIDToGRPC(eid *model.EventID) *eventer.EventID {
	return &eventer.EventID{Value: eid.String()}
}

func UserIDToGRPC(uid *model.UserID) *eventer.UserID {
	return &eventer.UserID{Value: int32(*uid)}
}

func GRPCToEvent(ge *eventer.Event) (*model.Event, error) {
	eid, err := GRPCToEventID(ge.Id)
	if err != nil {
		return &model.Event{}, err
	}

	e := &model.Event{
		ID:             *eid,
		Title:          ge.GetTitle(),
		StartDate:      ge.GetStartDate().AsTime().In(time.Local),
		EndDate:        ge.GetEndDate().AsTime().In(time.Local),
		Description:    ge.GetDescription(),
		UserID:         *GRPCToUserID(ge.UserId),
		NotifyUserTime: ge.GetNotifyUserTime(),
	}
	return e, nil
}

func GRPCToEventID(geid *eventer.EventID) (*model.EventID, error) {
	b, err := uuid.Parse(geid.GetValue())
	if err != nil {
		err = fmt.Errorf("fail uuid.Parse with param %v: %w", geid.GetValue(), err)
	}
	return &b, err
}

func GRPCToUserID(guid *eventer.UserID) *model.UserID {
	uid := model.UserID(guid.GetValue())
	return &uid
}

func ListModelToListGRPC(elist []model.Event) []*eventer.Event {
	if len(elist) == 0 {
		return nil
	}
	glist := make([]*eventer.Event, 0)
	for _, e := range elist {
		glist = append(glist, EventToGRPC(&e))
	}
	return glist
}
