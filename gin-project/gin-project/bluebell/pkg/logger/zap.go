package logger

import (
	"bluebell/config"
	"bluebell/pkg/utils"
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var lg *zap.Logger

// Init 初始化lg
func Init(cfg *config.LogConfig, mode string) (err error) {
	writeSyncer := getLogWriter(cfg.Dir, cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder(cfg.Format)
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	switch {
	case mode == "release":
		core = zapcore.NewCore(encoder, writeSyncer, l)
	case mode == "test":
		core = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), l)
	default:
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)

	}

	lg := zap.New(
		core,
		zap.AddCaller(),                      // 开启开发模式，堆栈跟踪；日志打印输出文件名, 行号, 函数名
		zap.AddCallerSkip(1),                 // 向上跳 1 层
		zap.AddStacktrace(zapcore.WarnLevel), // warn以上级别才输出堆栈信息
		zap.Development(),                    // 可输出 dpanic, panic 级别的日志
		zap.Fields(
			zap.String("server", "bluebell"), // 设置初始化字段
		))

	zap.ReplaceGlobals(lg)
	zap.L().Info("init logger success")
	return
}

// 日志key配置
func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	switch {
	case format == "json":
		return zapcore.NewJSONEncoder(encoderConfig)

	default:
		return zapcore.NewConsoleEncoder(encoderConfig)

	}
}

// 日志文件切割
func getLogWriter(path, filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 日志文件
	if err := utils.PathExists(path); err != nil {
		panic(err)
	}

	file := fmt.Sprintf("./bluebell/%s/%s", path, filename)

	// 切割配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
