package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	// internalhttp "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
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
	fmt.Printf("Create logger: %+v\n", logg)

	store, err := initStorage(conf.Storage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Create storage: %+v\n", store)

	var _ = store
	//storage := memorystorage.New()
	/*
		calendar := app.New(logg, storage)

		server := internalhttp.NewServer(logg, calendar)

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http server: " + err.Error())
			}
		}()

		logg.Info("calendar is running...")

		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	*/
}

// Перенести в storage куда-то
func initStorage(conf *config.StorageConf) (store storage.EventStorage, err error) {

	storageTempl := storage.New(conf)

	switch conf.Type {
	case "memory":
		return memorystorage.New(storageTempl), nil
	case "sql":
		return sqlstorage.New(storageTempl), nil
	}
	return nil, errors.New("wrong type of storage")
}
