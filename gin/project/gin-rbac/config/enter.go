package config

type Config struct {
	System *System // System configuration 系统配置
	MySQL  *MySQL  // MySQL configuration mysql配置
	Redis  *Redis  // Redis configuration redis配置
	Log    *Log    // Log configuration 日志配置
	JWT    *JWT    // JWT configuration JWT配置
}
