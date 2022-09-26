package grpc

import (
	"fmt"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/api/stubs/eventer"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelToGRPC(e *model.Event) *eventer.Event {
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

func GRPCToModel(ge *eventer.Event) (*model.Event, error) {
	eid, err := GRPCToEventID(ge.Id)
	if err != nil {
		return &model.Event{}, err
	}
	e := &model.Event{
		ID:             *eid,
		Title:          ge.GetTitle(),
		StartDate:      ge.GetStartDate().AsTime(),
		EndDate:        ge.GetEndDate().AsTime(),
		Description:    ge.GetDescription(),
		UserID:         *GRPCToUserID(ge.UserId),
		NotifyUserTime: ge.GetNotifyUserTime(),
	}
	return e, nil
}

func GRPCToEventID(geid *eventer.EventID) (*model.EventID, error) {
	b, err := uuid.FromBytes([]byte(geid.GetValue()))
	if err != nil {
		err = fmt.Errorf("fail uuid.FromBytes with param %v: %w", geid.GetValue(), err)
	}
	return &b, err
}

func GRPCToUserID(guid *eventer.UserID) *model.UserID {
	uid := model.UserID(guid.GetValue())
	return &uid
}

func ListModelToListGRPC(elist []model.Event) []*eventer.Event {
	glist := make([]*eventer.Event, 0)
	for _, e := range elist {
		glist = append(glist, ModelToGRPC(&e))
	}
	return glist
}
