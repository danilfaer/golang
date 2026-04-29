package config

// LoggerConfig — настройки логгера (совпадает по смыслу с order).
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// GRPCConfig — адрес bind gRPC-сервера inventory (host:port).
type GRPCConfig interface {
	Address() string
}

// MongoConfig — подключение к MongoDB и имя БД для коллекций.
type MongoConfig interface {
	URI() string
	Database() string
}
