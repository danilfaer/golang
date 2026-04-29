package config

import (
	"fmt"
	"os"

	"github.com/danilfaer/golang/order/internal/config/env"
	"github.com/joho/godotenv"
)

// appConfig — единая точка доступа к загруженной конфигурации внутри пакета; снаружи — только AppConfig().
var appConfig *Config

// Config — агрегат всех подсистемных настроек (интерфейсы из env-парсеров).
type Config struct {
	Logger    LoggerConfig
	HTTP      OrderHTTPConfig
	Payment   OrderPaymentGRPCConfig
	Inventory OrderInventoryGRPCConfig
	Postgres  PosgresConfig
}

// AppConfig возвращает загруженный конфиг. Перед вызовом должен быть успешно выполнен Load().
func AppConfig() *Config {
	return appConfig
}

// Load поднимает .env через godotenv, затем парсит env в подконфиги.
// path — пути к файлам (как в godotenv.Load); если не передан ни один, библиотека по умолчанию читает «.env» в cwd.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	logCfg, err := env.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("logger config: %w", err)
	}
	httpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return fmt.Errorf("order http config: %w", err)
	}
	payCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return fmt.Errorf("payment grpc config: %w", err)
	}
	invCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return fmt.Errorf("inventory grpc config: %w", err)
	}
	pgCfg, err := env.NewPostgresConfig()
	if err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	appConfig = &Config{
		Logger:    logCfg,
		HTTP:      httpCfg,
		Payment:   payCfg,
		Inventory: invCfg,
		Postgres:  pgCfg,
	}
	return nil
}
