package config

import (
	"fmt"
	"os"

	"github.com/danilfaer/golang/payment/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *Config

// Config — настройки payment gRPC-сервиса.
type Config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
}

// AppConfig возвращает конфиг после успешного Load().
func AppConfig() *Config {
	return appConfig
}

// Load читает .env и парсит переменные окружения.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	logCfg, err := env.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("logger config: %w", err)
	}
	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return fmt.Errorf("grpc config: %w", err)
	}

	appConfig = &Config{
		Logger: logCfg,
		GRPC:   grpcCfg,
	}
	return nil
}
