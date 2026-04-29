package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

// grpcListenEnvConfig — host/port для net.Listen gRPC-сервера (как в шаблоне compose: GRPC_HOST, GRPC_PORT).
type grpcListenEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type grpcListenConfig struct {
	raw grpcListenEnvConfig
}

// NewGRPCConfig парсит GRPC_HOST и GRPC_PORT.
func NewGRPCConfig() (*grpcListenConfig, error) {
	var raw grpcListenEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &grpcListenConfig{raw: raw}, nil
}

// Address возвращает host:port для Listen.
func (cfg *grpcListenConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
