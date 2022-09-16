package storage

import "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"

type Storage struct { // TODO
	Store       string // Тип хранилища "memory" или "sql"
	DatabaseURL string // Адрес базы данных
}

func New(conf config.StorageConf) *Storage {
	return &Storage{Store: conf.Store,
		DatabaseURL: conf.DatabaseURL}
}
