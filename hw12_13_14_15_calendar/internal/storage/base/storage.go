package basestorage

import (
	"context"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func InitStorage(conf *config.StorageConf, logg *logger.Logger) (store storage.EventStorage, err error) {
	storageTempl := storage.New(conf)

	logg.Infof("create %s storage", conf.Type)
	switch conf.Type {
	case "memory":
		return memorystorage.New(storageTempl), nil
	case "sql":
		store := sqlstorage.New(storageTempl)
		ctx := context.TODO()

		logg.Infof("connect to db by url: %s", conf.DatabaseURL)
		err := store.Connect(ctx) // logg
		if err != nil {
			logg.Errorf("failed to connect to sql db: %s\n", err)
			return nil, sqlstorage.ErrorConnectDB
		}

		err = store.PreparedQueries(ctx)
		if err != nil {
			logg.Errorf("failed to prepare queries: %s\n", err)
			return nil, sqlstorage.ErrorPrepareQuery
		}
		return store, nil
	}
	return nil, storage.ErrorWrongTypeStorage
}
