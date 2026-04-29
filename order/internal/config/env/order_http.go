package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

// orderHTTPEnvConfig — сырые значения env для HTTP-сервиса order.
type orderHTTPEnvConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
}

// orderHTTPConfig — конфиг HTTP-слушателя с методом Address().
type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

// NewOrderHTTPConfig парсит переменные окружения и возвращает конфиг HTTP.
func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &orderHTTPConfig{raw: raw}, nil
}

// Address возвращает адрес bind в формате host:port.
func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
