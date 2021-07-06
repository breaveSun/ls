package config
type RedisConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstructure:"expire"`
}