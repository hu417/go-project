package config

// 配置信息参数
type Mysql struct {
	Host         string `mapstructure:"host" yaml:"host" json:"host"`
	Port         string `mapstructure:"port"  json:"port" yaml:"port"`
	DbName       string `mapstructure:"db-name" yaml:"db-name" json:"db-name"`
	UserName     string `mapstructure:"username" yaml:"username" json:"username"`
	Password     string `mapstructure:"password" yaml:"password" json:"password"`
	Params       string `mapstructure:"params" yaml:"params" json:"params" description:"字符集/时间/ssl"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      string `mapstructure:"log_mode" yaml:"log_mode" json:"log_mode" description:"日志级别"`
}

// Dsn 获取dsn
func (m *Mysql) Dsn() string {

	return m.UserName + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DbName + "?" + m.Params

}
