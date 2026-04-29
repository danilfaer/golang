package config

import (
	"fmt"
	"os"

	"github.com/danilfaer/golang/inventory/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *Config

// Config — все подсистемные настройки inventory-сервиса.
type Config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
	Mongo  MongoConfig
}

// AppConfig возвращает конфиг после успешного Load().
func AppConfig() *Config {
	return appConfig
}

// Load читает .env (godotenv) и парсит переменные окружения.
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
	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return fmt.Errorf("mongo config: %w", err)
	}

	appConfig = &Config{
		Logger: logCfg,
		GRPC:   grpcCfg,
		Mongo:  mongoCfg,
	}
	return nil
}
