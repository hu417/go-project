package conf

// AppConfig 项目配置
type Config struct {
	System   `mapstructure:"system"`
	Log      `mapstructure:"log"`
	Database `mapstructure:"database"`
	Redis    `mapstructure:"redis"`
	Jwt      `mapstructure:"jwt"`
}

// SystemConfig 系统配置
type System struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      string `mapstructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

// LogConfig 日志配置
type Log struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	RootDir    string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`
	Format     string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine   bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`    // day
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type Jwt struct {
	Secret                  string `mapstructure:"secret" json:"secret" yaml:"secret"`
	JwtTtl                  int64  `mapstructure:"jwt_ttl" json:"jwt_ttl" yaml:"jwt_ttl"`                                                          // token 有效期（秒）
	JwtBlacklistGracePeriod int64  `mapstructure:"jwt_blacklist_grace_period" json:"jwt_blacklist_grace_period" yaml:"jwt_blacklist_grace_period"` // 黑名单宽限时间（秒)
	RefreshGracePeriod      int64  `mapstructure:"refresh_grace_period" json:"refresh_grace_period" yaml:"refresh_grace_period"`                   // token 自动刷新宽限时间（秒）
}
