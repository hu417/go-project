package config

type Config struct {
	App   App   `mapstructure:"app" json:"app" yaml:"app"`
	Log   Log   `mapstructure:"log" json:"log" yaml:"log"`
	Mysql MySql `mapstructure:"database" json:"database" yaml:"database"`
	Jwt   Jwt   `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
}
