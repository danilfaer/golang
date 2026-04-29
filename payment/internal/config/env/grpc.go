package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

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

func (cfg *grpcListenConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
