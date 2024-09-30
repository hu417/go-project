package conf

import "fmt"

// MySQLConfig 数据库配置
type Database struct {
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	DBName              string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	UserName            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Config              string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             bool   `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	Log_live            string `mapstructure:"log_live" json:"log_live" yaml:"log_live"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

func (d *Database) Dns() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		d.UserName,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.Config,
	)
}
