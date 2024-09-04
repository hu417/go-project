package config

type Jwt struct {
	SignKey  string `mapstructure:"sign_key" json:"sign_key" yaml:"sign_key"`
	Expires  int64  `mapstructure:"expires" json:"expires" yaml:"expires"`
}