package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// mongoEnvConfig — либо готовый MONGO_DB_URI, либо компоненты как в inventory.env.template после envsubst.
type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_DATABASE,required"`
	User     string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
}

type mongoConfig struct {
	raw      mongoEnvConfig
}

// NewMongoConfig собирает URI: приоритет у MONGO_DB_URI, иначе mongodb://user:pass@host:port/db?authSource=...
func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &mongoConfig{raw: raw}, nil
}

func (cfg *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.AuthDB,
	)
}

func (cfg *mongoConfig) Database() string {
	return cfg.raw.Database
}
