package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)



// 初始化zap实例
func InitzapLogger() *zap.Logger {

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "gin-app"))
	// 构造日志
	logger := zap.New(getEncoding(), caller, development, filed)
	//defer logger.Sync()

	return logger
}

func getEncoding() zapcore.Core {
	// 设置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // zapcore.CapitalLevelEncoder, // 大写编码器; zapcore.LowercaseLevelEncoder, // 小写编码器; zapcore.CapitalColorLevelEncoder //按级别显示不同颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // zapcore.FullCallerEncoder, // 全路径编码器; zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	loglevel := "debug"
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	atomicLevel.SetLevel(level)

	// 日志切割
	hook := lumberjack.Logger{
		Filename:   "./log/app.log", // 日志文件路径
		MaxSize:    128,             // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,              // 日志文件最多保存多少个备份
		MaxAge:     7,               // 文件最多保存多少天
		Compress:   true,            // 是否压缩
	}

	encoding := "json"
	if encoding == "json" {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			atomicLevel, // 日志级别
		)
		return core
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 将err日志单独输出到文件
	// encoder := getEncoder()
	// // test.log记录全量日志
	// logF, _ := os.Create("./test.log")
	// c1 := zapcore.NewCore(encoder, zapcore.AddSync(logF), zapcore.DebugLevel)
	// // test.err.log记录ERROR级别的日志
	// errF, _ := os.Create("./test.err.log")
	// c2 := zapcore.NewCore(encoder, zapcore.AddSync(errF), zap.ErrorLevel)
	// // 使用NewTee将c1和c2合并到core
	// core := zapcore.NewTee(c1, c2)

	return core
}

