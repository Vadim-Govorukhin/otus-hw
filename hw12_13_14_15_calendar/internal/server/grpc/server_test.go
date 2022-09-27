package internalgrpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/api/stubs/eventer"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	internalgrpc "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	basestorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/base"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GRPCListsEqual(t *testing.T, exp, act []*eventer.Event) {
	require.Equal(t, len(exp), len(act))

	for _, elExp := range exp {
		for _, elAct := range exp {
			if elExp.GetId() == elAct.GetId() {
				require.Equal(t, elExp, elAct)
			}
		}
	}
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var calendarApp *app.App

func testStart() (*zap.SugaredLogger, *grpc.ClientConn) {
	conf := config.NewConfig()
	toml.DecodeFile("../../../configs/calendar_config.toml", &conf)

	logger, _ := zap.NewDevelopment()
	logg := logger.Sugar()

	conf.Storage.Type = "memory" // tests based on memory storage
	storage, _ := basestorage.InitStorage(conf.Storage, logg)
	calendarApp = app.New(logg, storage)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	service := internalgrpc.NewServer(logg, calendarApp, conf.GRPCServer)
	eventer.RegisterCalendarServer(s, service)
	go func() {
		if err := s.Serve(lis); err != nil {
			logg.Fatalf("Server exited with error: %v", err)
		}
	}()
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logg.Fatalf("Server exited with error: %v", err)
	}
	return logg, conn

}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGRPCService(t *testing.T) {
	logg, conn := testStart()
	defer conn.Close()
	ctx := context.Background()

	client := eventer.NewCalendarClient(conn)

	empty := new(emptypb.Empty)
	date := time.Date(2022, time.September, 16, 1, 2, 3, 0, time.UTC)

	t.Run("get empty event list", func(t *testing.T) {
		logg.Infof("============== start test %s ==========", t.Name())

		resp, err := client.ListAllEvent(ctx, empty)
		require.NoError(t, err)
		expList := []model.Event{}
		require.Equal(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())
	})

	t.Run("create event", func(t *testing.T) {
		logg.Infof("============== start test %s ==========", t.Name())

		respID, err := client.CreateEvent(ctx, internalgrpc.EventToGRPC(&storage.TestEvent))
		require.NoError(t, err)
		require.Equal(t, storage.TestEvent.ID.String(), respID.GetValue())

		respID, err = client.CreateEvent(ctx, internalgrpc.EventToGRPC(&storage.TestEvent2))
		require.NoError(t, err)
		require.Equal(t, storage.TestEvent2.ID.String(), respID.GetValue())

		respID, err = client.CreateEvent(ctx, internalgrpc.EventToGRPC(&storage.TestEvent3))
		require.NoError(t, err)
		require.Equal(t, storage.TestEvent3.ID.String(), respID.GetValue())

		resp, err := client.ListAllEvent(ctx, empty)
		require.NoError(t, err)
		expList := []model.Event{storage.TestEvent2, storage.TestEvent, storage.TestEvent3}
		GRPCListsEqual(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())
	})

	t.Run("get list of events", func(t *testing.T) {
		logg.Infof("============== start test %s ==========", t.Name())

		resp, err := client.ListEventByDay(ctx, timestamppb.New(date))
		require.NoError(t, err)
		expList := []model.Event{storage.TestEvent3, storage.TestEvent2}
		GRPCListsEqual(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())

		resp, err = client.ListEventByMonth(ctx, timestamppb.New(date))
		require.NoError(t, err)
		expList = []model.Event{storage.TestEvent2, storage.TestEvent}
		GRPCListsEqual(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())

		var uid = 0
		resp, err = client.ListAllEventByUser(ctx, internalgrpc.UserIDToGRPC(&uid))
		require.NoError(t, err)
		expList = []model.Event{storage.TestEvent, storage.TestEvent3}
		GRPCListsEqual(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())
	})

	t.Run("Update and delete event", func(t *testing.T) {
		logg.Infof("============== start test %s ==========", t.Name())

		respID, err := client.DeleteEvent(ctx, internalgrpc.EventIDToGRPC(&storage.TestEvent.ID))
		require.NoError(t, err)
		require.Equal(t, storage.TestEvent.ID.String(), respID.GetValue())

		resp, err := client.ListAllEvent(ctx, empty)
		require.NoError(t, err)
		expList := []model.Event{storage.TestEvent2, storage.TestEvent3}
		GRPCListsEqual(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())

		updateReq := eventer.UpdateEventRequest{EventId: internalgrpc.EventIDToGRPC(&storage.TestEvent2.ID),
			Event: internalgrpc.EventToGRPC(&storage.TmpEvent)}
		respID, err = client.UprateEvent(ctx, &updateReq)
		require.NoError(t, err)
		require.Equal(t, storage.TmpEvent.ID.String(), respID.GetValue())

		resp, err = client.ListAllEvent(ctx, empty)
		require.NoError(t, err)
		expList = []model.Event{storage.TmpEvent, storage.TestEvent3}
		GRPCListsEqual(t, internalgrpc.ListModelToListGRPC(expList), resp.GetEvent())
	})
}
