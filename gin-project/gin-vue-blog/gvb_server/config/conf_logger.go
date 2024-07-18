package config

type Logger struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                           // 级别
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                        // 输出格式
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                        // 日志前缀
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                 // 日志文件夹
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine"`                 // 显示行
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`  // 输出控制台
	MaxSize       int    `mapstructure:"maxsize" json:"maxsize" yaml:"maxsize"`                     // 日志文件最大的尺寸(M), 超限后开始自动分割
	MaxBackups    int    `mapstructure:"maxbackups" json:"maxbackups" yaml:"maxbackups"`            // 保留旧文件的最大个数
	MaxAge        int    `mapstructure:"maxage" json:"maxage" yaml:"maxage"`                        // 保留旧文件的最大天数
	Compress      bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}
