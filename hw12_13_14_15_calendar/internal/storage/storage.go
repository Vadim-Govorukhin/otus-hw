package storage

type Storage struct { // TODO
	Store       string // Тип хранилища "memory" или "sql"
	DatabaseURL string // Адрес базы данных
}

func New(store, databaseURL string) *Storage {
	return &Storage{Store: store,
		DatabaseURL: databaseURL}
}
