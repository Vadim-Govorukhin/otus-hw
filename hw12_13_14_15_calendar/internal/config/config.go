package config

type Config struct {
	Storage    *StorageConf
	Logger     *LoggerConf // Логирование
	SessionKey string      `toml:"session_key"` // Ключ для генерации сессий
	BindAddr   string      `toml:"bind_addr"`   // Адрес (порт), на котором запускаем веб сервер
	// TODO
}

type LoggerConf struct {
	Level string `toml:"level"` // Уровень логирования
	// TODO
}

type StorageConf struct {
	Type        string `toml:"type"`         // Тип хранилища "memory" или "sql"
	DatabaseURL string `toml:"database_url"` // Адрес базы данных
}

func NewConfig() *Config {
	return &Config{}
}

// TODO
