package log

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const ctxLoggerKey = "zapLogger"

type Logger struct {
	*zap.Logger
}

func NewLog(conf *viper.Viper) *Logger {
	// log address "out.log" User-defined
	lp := conf.GetString("log.log_file_name") // 日志文件
	lv := conf.GetString("log.log_level")     // 日志级别

	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch lv {
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

	// 日志文件切割
	hook := lumberjack.Logger{
		Filename:   lp,                             // Log file path
		MaxSize:    conf.GetInt("log.max_size"),    // Maximum size unit for each log file: M
		MaxBackups: conf.GetInt("log.max_backups"), // The maximum number of backups that can be saved for log files
		MaxAge:     conf.GetInt("log.max_age"),     // Maximum number of days the file can be saved
		Compress:   conf.GetBool("log.compress"),   // Compression or not
	}

	var encoder zapcore.Encoder
	// 定义key
	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	// 输出样式
	switch conf.GetString("log.encoding") {
	case "json":
		encoder = zapcore.NewJSONEncoder(cfg)
	case "console":
		encoder = zapcore.NewConsoleEncoder(cfg)
	default:
		encoder = zapcore.NewJSONEncoder(cfg)
	}

	// core
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // Print to console and file
		level,
	)
	option := []zap.Option{}
	if conf.GetString("env") != "prod" {
		option = append(option, zap.Development())
	}
	option = append(
		option,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel))
	return &Logger{
		zap.New(
			core,
			option...,
		)}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// WithValue Adds a field to the specified context
func (l *Logger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}
