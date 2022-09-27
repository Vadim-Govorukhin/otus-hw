package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	basestorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/base"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf := config.NewConfig()
	_, err := toml.DecodeFile(configFile, conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Read config: %+v\n", conf)

	logg := logger.New(conf.Logger)
	logg.Infof("Create logger: %T\n", logg)

	storage, err := basestorage.InitStorage(conf.Storage, logg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logg.Infof("Create storage: %T\n", storage)

	calendar := app.New(logg, storage)
	httpserver := internalhttp.NewServer(logg, calendar, conf.HTTPServer)
	grpcserver := internalgrpc.NewServer(logg, calendar, conf.GRPCServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpserver.Stop(ctx); err != nil {
			logg.Errorf("failed to stop http server: %s", err.Error())
		}

		if err := grpcserver.Stop(ctx); err != nil {
			logg.Errorf("failed to stop grpc server: %s", err.Error())
		}
		if err := storage.Close(ctx); err != nil {
			logg.Errorf("failed to close storage: %s", err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := httpserver.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	if err := grpcserver.Start(ctx); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
