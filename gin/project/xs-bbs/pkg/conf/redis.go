package conf

// RedisConfig redis配置
type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
}