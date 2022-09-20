package storage

import "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"

type Storage struct { // TODO
	Type        string // Тип хранилища "memory" или "sql"
	DatabaseURL string // Адрес базы данных
}

func New(storeConf *config.StorageConf) *Storage {
	return &Storage{
		Type:        storeConf.Type,
		DatabaseURL: storeConf.DatabaseURL,
	}
}
