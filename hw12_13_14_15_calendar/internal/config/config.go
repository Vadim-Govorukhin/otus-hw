package config

type Config struct {
	BindAddr    string     `toml:"bind_addr"`    // Адрес (порт), на котором запускаем веб сервер
	Store       string     `toml:"store"`        // Тип хранилища "memory" или "sql"
	DatabaseURL string     `toml:"database_url"` // Адрес базы данных
	SessionKey  string     `toml:"session_key"`  // Ключ для генерации сессий
	Logger      LoggerConf // Логирование
	// TODO
}

type LoggerConf struct {
	Level string `toml:"level"` // Уровень логирования
	// TODO
}

func NewConfig() *Config {
	return &Config{}
}

// TODO
