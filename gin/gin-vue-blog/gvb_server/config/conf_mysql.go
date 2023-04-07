package config

// 配置信息参数
type Mysql struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DB           string `yaml:"db"`
	UserName     string `yaml:"username"`
	Password     string `yaml:"password"`
	Config       string `yaml:"config" description:"字符集/时间/ssl"`
	MaxIdleConns int    `json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      string `yaml:"log_mode" description:"日志级别"`
}

// Connection
func (m Mysql) Dsn() string {

	return m.UserName + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DB + "?" + m.Config

}
