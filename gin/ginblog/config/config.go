package config

type Config struct {
	Server Server `mapstructure:"server" yaml:"server" json:"server"`
	Mysql  Mysql  `mapstructure:"mysql" yaml:"mysql" json:"mysql"`
	Jwt    Jwt    `mapstructure:"jwt" yaml:"jwt" json:"jwt"`
}
