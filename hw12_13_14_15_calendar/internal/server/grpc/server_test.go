package internalgrpc_test

import (
	"context"
	"net"
	"testing"

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
)

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

func TestSayHello(t *testing.T) {
	logg, conn := testStart()
	defer conn.Close()
	ctx := context.Background()

	client := eventer.NewCalendarClient(conn)
	var empty *emptypb.Empty

	t.Run("create event", func(t *testing.T) {
		logg.Infof("============== start test %s ==========", t.Name())

		/*
			respID, err := client.CreateEvent(ctx, internalgrpc.ModelToGRPC(&storage.TestEvent))
			require.NoError(t, err)
			require.Equal(t, respID.GetValue(), storage.TestEvent.ID.String())
		*/
		resp, err := client.ListAllEvent(ctx, empty)
		require.NoError(t, err)
		expList := []model.Event{storage.TestEvent}
		require.Equal(t, resp.GetEvent(), internalgrpc.ListModelToListGRPC(expList))
	})
}
