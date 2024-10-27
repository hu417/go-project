package config

type JWT struct {
	Secret     string `yaml:"secret_key" mapstructure:"secret_key"` // 密钥
	Expiration int64 `yaml:"expiration" mapstructure:"expiration"` // 过期时间
}
