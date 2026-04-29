package config

// LoggerConfig — настройки логгера (как в order/inventory).
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// GRPCConfig — адрес bind gRPC-сервера payment.
type GRPCConfig interface {
	Address() string
}
