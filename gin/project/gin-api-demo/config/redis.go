package config

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
