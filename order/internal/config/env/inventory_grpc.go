package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

// inventoryEnvConfig — сырые значения env для gRPC-клиента Inventory.
type inventoryEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

// inventoryGRPCConfig — конфиг gRPC к Inventory с методом Address().
type inventoryGRPCConfig struct {
	raw inventoryEnvConfig
}

// NewInventoryGRPCConfig парсит переменные окружения и возвращает конфиг клиента.
func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &inventoryGRPCConfig{raw: raw}, nil
}

// Address возвращает адрес в формате host:port (пригодно для dial).
func (cfg *inventoryGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
