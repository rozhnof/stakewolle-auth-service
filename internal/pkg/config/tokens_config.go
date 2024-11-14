package config

type TokensConfig struct {
	AccessTokenTTL  string `yaml:"access_ttl"  env-required:"true"`
	RefreshTokenTTL string `yaml:"refresh_ttl" env-required:"true"`
}
