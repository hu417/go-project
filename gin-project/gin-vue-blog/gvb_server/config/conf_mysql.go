package config

// 配置信息参数
type Mysql struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DbName       string `yaml:"db-name"`
	UserName     string `yaml:"username"`
	Password     string `yaml:"password"`
	Params       string `yaml:"params" description:"字符集/时间/ssl"`
	MaxIdleConns int    `json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      string `yaml:"log_mode" description:"日志级别"`
}

// Connection
func (m *Mysql) Dsn() string {

	return m.UserName + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DbName + "?" + m.Params

}
