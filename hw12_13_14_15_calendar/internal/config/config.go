package config

type Config struct {
	Storage    *StorageConf
	Logger     *LoggerConf // Логирование
	HTTPServer *HTTPServerConf
	SessionKey string `toml:"session_key"` // Ключ для генерации сессий
}

type LoggerConf struct {
	Level string `toml:"level"` // Уровень логирования
	Path  string `toml:"path"`  // Путь до файла с логами
}

type StorageConf struct {
	Type        string `toml:"type"`         // Тип хранилища "memory" или "sql"
	DatabaseURL string `toml:"database_url"` // Адрес базы данных
}

type HTTPServerConf struct {
	Host    string `toml:"host"`
	Port    string `toml:"port"`
	LogPath string `toml:"log_path"`
}

func NewConfig() *Config {
	return &Config{}
}

// TODO
