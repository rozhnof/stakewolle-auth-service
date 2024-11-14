package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Mode     string       `yaml:"mode"    env-required:"true"`
	Server   ServerConfig `yaml:"server"  env-required:"true"`
	Logger   LoggerConfig `yaml:"logging" env-required:"true"`
	Tokens   TokensConfig `yaml:"tokens"  env-required:"true"`
	Postgres PostgresConfig
	Redis    RedisConfig
}

func NewConfig(configPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
