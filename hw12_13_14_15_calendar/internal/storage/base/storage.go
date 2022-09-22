package basestorage

import (
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func InitStorage(conf *config.StorageConf, logg *logger.Logger) (store storage.EventStorage, err error) {
	storageTempl := storage.New(conf)

	switch conf.Type {
	case "memory":
		return memorystorage.New(storageTempl), nil
	case "sql":
		return sqlstorage.New(storageTempl), nil
	}
	return nil, storage.ErrorWrongTypeStorage
}
