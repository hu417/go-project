package config

type Server struct {
	HttpPort string `mapstructure:"http_port" json:"http_port" yaml:"http_port"`
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	RunMode  string `mapstructure:"run_mode" json:"run_mode" yaml:"run_mode"`
	Timeout  struct {
		Read  int `mapstructure:"read" json:"read" yaml:"read"`
		Write int `mapstructure:"write" json:"write" yaml:"write"`
	} `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}
