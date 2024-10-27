package utils

import (
	"io"
	"log"
	"os"
	"time"

	"gin-api-demo/global"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

// 自定义 gorm Writer
func getGormLogWriter() logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if global.Conf.Mysql.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   global.Conf.Log.RootDir + "/" + global.Conf.Mysql.LogFilename,
			MaxSize:    global.Conf.Log.MaxSize,
			MaxBackups: global.Conf.Log.MaxBackups,
			MaxAge:     global.Conf.Log.MaxAge,
			Compress:   global.Conf.Log.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func GetGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch global.Conf.Mysql.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond,                 // 慢 SQL 阈值
		LogLevel:                  logMode,                                // 日志级别
		IgnoreRecordNotFoundError: false,                                  // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  !global.Conf.Mysql.EnableFileLogWriter, // 禁用彩色打印
	})
}
