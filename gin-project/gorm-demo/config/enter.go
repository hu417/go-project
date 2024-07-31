package config

type Conf struct {
	System System `mapstructure:"system" yaml:"system" json:"system"`
	Mysql Mysql `mapstructure:"mysql" yaml:"mysql" json:"mysql"`
}
