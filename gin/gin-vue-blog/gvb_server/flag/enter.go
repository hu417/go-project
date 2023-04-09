package flag

import (
	sys_flag "flag"
)

type Options struct {
	DB      bool
	Version bool
}

// Parse解析命令行参数
func Parse() Options {
	version := sys_flag.Bool("v", false, "项目版本")
	db := sys_flag.Bool("db", false, "初始化数据库")
	// 解析命令行参数写入注册的flag里
	sys_flag.Parse()
	return Options{
		Version: *version,
		DB:      *db,
	}

}

// IsWebStop 是否停止web项目
func IsWebStop(options Options) bool {
	if options.DB {
		return true
	}
	if options.Version {
		return true
	}

	return false
}

// SwitchOption 根据命令执行不同的函数
func SwitchOption(options Options) {
	if options.DB {
		Makmigrations()
	}
	if options.Version {
		Version()
	}
}
