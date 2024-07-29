package config

// 配置信息参数
type Elastic struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	UserName        string `yaml:"username"`
	Password        string `yaml:"password"`
	SetGzip         bool   `yaml:"set_gizp"`
	SetSniff        bool   `yaml:"set_sniff"`
	HealthcheckTime int    `yaml:"health_check_tiame" description:"心跳检测时间"`
	SetErrorLog     string `json:"set_errorlog" yaml:"set_errorlog"`
	SetInfoLog      string `json:"set_infolog" yaml:"set_infolog"`
}
