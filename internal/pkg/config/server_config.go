package config

type ServerConfig struct {
	Address string `yaml:"address" env-required:"true"`
}
