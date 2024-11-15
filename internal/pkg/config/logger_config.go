package config

type LoggerConfig struct {
	Level string `yaml:"level" env-required:"true"`
	Path  string `yaml:"path"`
}
