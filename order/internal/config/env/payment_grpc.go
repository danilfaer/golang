package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

// paymentEnvConfig — сырые значения env для gRPC-клиента Payment.
type paymentEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

// paymentGRPCConfig — конфиг gRPC к Payment с методом Address().
type paymentGRPCConfig struct {
	raw paymentEnvConfig
}

// NewPaymentGRPCConfig парсит переменные окружения и возвращает конфиг клиента.
func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	var raw paymentEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &paymentGRPCConfig{raw: raw}, nil
}

// Address возвращает адрес в формате host:port (пригодно для dial).
func (cfg *paymentGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
