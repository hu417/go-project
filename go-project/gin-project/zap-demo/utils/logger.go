package utils

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化Logger
func InitLogger(mode, format string) {
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println(err) //print the error if obtained
	}
	file := fmt.Sprintf("%s/app-%s.log", directory, FormatTime())
	writeSyncer := getLogWriter(
		file,
		10,
		2,
		30,
	)
	encoder := getEncoder(format)

	var core zapcore.Core
	// 创建日志级别过滤器
	switch mode {
	case "test":
		// 进入测试模式，日志输出到终端
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	case "prod":
		// 进入生产模式，日志输出到文件
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	default:
		// 进入开发模式，日志输出到终端和文件
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	}

	lg := zap.New(core,
		zap.AddCaller(),                      // 开启开发模式，堆栈跟踪；日志打印输出文件名, 行号, 函数名
		zap.AddCallerSkip(0),                 // 向上跳 1 层
		zap.AddStacktrace(zapcore.WarnLevel), // warn以上级别才输出堆栈信息
		zap.Development(),                    // 可输出 dpanic, panic 级别的日志
		zap.Fields(
			zap.String("server", "app"), // 设置初始化字段
		))
	defer lg.Sync() // 确保缓冲区中的日志条目被刷新

	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可

}

// zap配置编码器
func getEncoder(format string) zapcore.Encoder {
	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{ // 创建编码配置
		TimeKey:        "Time",                        // 时间键
		LevelKey:       "Level",                       // 日志级别键
		NameKey:        "Log",                         // 日志名称键
		CallerKey:      "Call",                        // 日志调用键
		MessageKey:     "Msg",                         // 日志消息键
		StacktraceKey:  "Stacktrace",                  // 堆栈跟踪键
		LineEnding:     zapcore.DefaultLineEnding,     // 行结束符,默认为 \n
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 日志级别编码器,将日志级别转换为大写
		EncodeTime:     CustomTimeEncoder,             // 时间编码器,zapcore.ISO8601TimeEncoder 将时间格式化为 ISO8601 格式
		EncodeDuration: zapcore.StringDurationEncoder, // 持续时间编码器,将持续时间编码为字符串
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用编码器,显示文件名和行号
	}

	switch format {
	case "json", "JSON":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 写入日志文件
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 日志文件分割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,   // 每个日志文件最大 10 MB
		MaxBackups: maxBackup, // 保留最近的 5 个日志文件
		MaxAge:     maxAge,    // 保留最近 30 天的日志
		Compress:   true,      // 旧日志文件压缩
	}

	// 利用io.MultiWriter支持文件和终端两个输出目标
	// dest := io.MultiWriter(lumberJackLogger, os.Stdout)

	return zapcore.AddSync(lumberJackLogger)
}

// 格式化时间
func FormatTime() string {
	return time.Now().Format("2006-01-02_15")
}
