package grpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ginzap "github.com/akath19/gin-zap"
	"github.com/gin-gonic/gin"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	jsontime "github.com/liamylian/jsontime/v2/v2"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/api/stubs/eventer"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var calendarApp *app.App
var _ server.Server = &Server{}

type Server struct {
	eventer.UnimplementedCalendarServer
	address string
	server  *grpc.Server
}

func (s *Server) Start(ctx context.Context) error {
	jsontime.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
	fmt.Print("Start gRPC listener")

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("fail to listen: %w", err)
	}

	fmt.Print("regislet gRPC service")
	eventer.RegisterCalendarServer(s.server, *s) //
	err = s.server.Serve(lis)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func configureLoggerGin(logg *logger.Logger, logPath string) gin.HandlerFunc {
	curDir, err := os.Getwd()
	if err != nil {
		logg.DPanicf("can't get working dir %w", err)
	}

	logFile, err := os.OpenFile(filepath.Join(curDir, logPath),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		logg.DPanicf("can't open log file %w", err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	ginLogg := ginzap.Logger(3*time.Second, logg.Desugar())
	return ginLogg
}

func NewServer(logg *logger.Logger, app *app.App, conf *config.GRPCServerConf) *Server {
	calendarApp = app
	address := net.JoinHostPort(conf.Host, conf.Port)
	ginLogg := configureLoggerGin(logg, conf.LogPath)
	var _ = ginLogg
	return &Server{
		address: address,
		server:  grpc.NewServer(grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(logg.Desugar())))}
}

func (s Server) CreateEvent(ctx context.Context, ge *eventer.Event) (*eventer.EventID, error) {
	e, err := GRPCToModel(ge)
	if err != nil {
		return &eventer.EventID{}, err
	}
	eid, err := calendarApp.Create(*e)
	if err != nil {
		return nil, err
	}
	return EventIDToGRPC(&eid), err
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
	list, err := calendarApp.ListAllEvents()
	resp := &eventer.EventResponse{Event: ListModelToListGRPC(list)}
	return resp, err
}

func (s Server) ListAllEventByUser(ctx context.Context, uid *eventer.UserID) (*eventer.EventResponse, error) {
	panic("not implemented") // TODO: Implement
}
