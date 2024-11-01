package bootstrap

import (
	"os"
	"time"

	"gin-rbac/global"
	"gin-rbac/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitializeLog() *zap.SugaredLogger {
	// 创建根目录
	createRootDir()

	// 设置日志等级
	var core zapcore.Core
	// 创建日志级别过滤器
	switch global.Config.System.Env {
	case "test":
		// 进入测试模式，日志输出到终端
		core = zapcore.NewTee(
			zapcore.NewCore(getZapCore(), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	case "prod":
		// 进入生产模式，日志输出到文件
		core = zapcore.NewCore(getZapCore(), getLogWriter(), zapcore.InfoLevel)
	default:
		// 进入开发模式，日志输出到终端和文件
		core = zapcore.NewTee(
			zapcore.NewCore(getZapCore(), getLogWriter(), zapcore.InfoLevel),
			zapcore.NewCore(getZapCore(), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	}

	var options []zap.Option // zap 配置项

	if global.Config.Log.ShowLine {
		options = append(options,
			zap.AddCaller(), // 开启开发模式，堆栈跟踪；日志打印输出文件名, 行号, 函数名
			// zap.AddCallerSkip(1),                 // 向上跳 1 层
			zap.AddStacktrace(zapcore.ErrorLevel), // warn以上级别才输出堆栈信息
			zap.Development(),                     // 可输出 dpanic, panic 级别的日志
		)
	}
	options = append(options,
		zap.Fields(
			zap.String("server", "app"), // 设置初始化字段
		))

	// 初始化 zap
	lg := zap.New(core, options...).Sugar()
	// zap.ReplaceGlobals(lg)
	return lg
}

// 创建根目录
func createRootDir() {
	if ok, _ := utils.PathExists(global.Config.Log.RootDir); !ok {
		_ = os.Mkdir(global.Config.Log.RootDir, os.ModePerm)
	}
}

// 扩展 Zap
func getZapCore() zapcore.Encoder {

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{ // 创建编码配置
		TimeKey:       "Time",                      // 时间键
		LevelKey:      "Level",                     // 日志级别键
		NameKey:       "Log",                       // 日志名称键
		CallerKey:     "Call",                      // 日志调用键
		MessageKey:    "Msg",                       // 日志消息键
		StacktraceKey: "Stacktrace",                // 堆栈跟踪键
		LineEnding:    zapcore.DefaultLineEnding,   // 行结束符,默认为 \n
		EncodeLevel:   zapcore.CapitalLevelEncoder, // 日志级别编码器,将日志级别转换为大写
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) { // 时间编码器,zapcore.ISO8601TimeEncoder 将时间格式化为 ISO8601 格式
			encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
		},
		EncodeDuration: zapcore.StringDurationEncoder, // 持续时间编码器,将持续时间编码为字符串
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用编码器,显示文件名和行号
	}

	// 设置编码器
	if global.Config.Log.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)

}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.Config.Log.RootDir + "/" + global.Config.Log.Filename,
		MaxSize:    global.Config.Log.MaxSize,
		MaxBackups: global.Config.Log.MaxBackups,
		MaxAge:     global.Config.Log.MaxAge,
		Compress:   global.Config.Log.Compress,
	}

	return zapcore.AddSync(file)
}
