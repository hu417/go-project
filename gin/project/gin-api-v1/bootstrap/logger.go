package bootstrap

import (
	"os"
	"time"

	"gin-api-demo/global"
	"gin-api-demo/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitializeLog() *zap.Logger {
	// 创建根目录
	createRootDir()

	// 设置日志等级
	var core zapcore.Core
	// 创建日志级别过滤器
	switch global.Conf.App.Env {
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

	if global.Conf.Log.ShowLine {
		options = append(options,
			zap.AddCaller(), // 开启开发模式，堆栈跟踪；日志打印输出文件名, 行号, 函数名
			//zap.AddCallerSkip(1),                  // 向上跳 1 层
			zap.AddStacktrace(zapcore.ErrorLevel), // error以上级别才输出堆栈信息
			zap.Development(),                     // 可输出 dpanic, panic 级别的日志
		)
	}
	options = append(options,
		zap.Fields(
			zap.String("server", "app"), // 设置初始化字段
		))

	// 初始化 zap
	lg := zap.New(core, options...)
	zap.ReplaceGlobals(lg)
	return lg
}

// 创建根目录
func createRootDir() {
	if ok, _ := utils.PathExists(global.Conf.Log.RootDir); !ok {
		_ = os.Mkdir(global.Conf.Log.RootDir, os.ModePerm)
	}
}

// // 设置日志等级
// func setLogLevel() {
// 	switch global.Conf.Log.Level {
// 	case "debug":
// 		level = zap.DebugLevel
// 		options = append(options, zap.AddStacktrace(level))
// 	case "info":
// 		level = zap.InfoLevel
// 	case "warn":
// 		level = zap.WarnLevel
// 	case "error":
// 		level = zap.ErrorLevel
// 		options = append(options, zap.AddStacktrace(level))
// 	case "dpanic":
// 		level = zap.DPanicLevel
// 	case "panic":
// 		level = zap.PanicLevel
// 	case "fatal":
// 		level = zap.FatalLevel
// 	default:
// 		level = zap.InfoLevel
// 	}
// }

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
		FunctionKey:   zapcore.OmitKey,             //
		LineEnding:    zapcore.DefaultLineEnding,   // 行结束符,默认为 \n
		EncodeLevel:   zapcore.CapitalLevelEncoder, // 日志级别编码器,将日志级别转换为大写
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) { // 时间编码器,zapcore.ISO8601TimeEncoder 将时间格式化为 ISO8601 格式
			encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
		},
		EncodeDuration: zapcore.StringDurationEncoder, // 持续时间编码器,将持续时间编码为字符串
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用编码器,显示文件名和行号
	}

	// 设置编码器
	if global.Conf.Log.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)

}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   global.Conf.Log.RootDir + "/" + global.Conf.Log.Filename,
		MaxSize:    global.Conf.Log.MaxSize,
		MaxBackups: global.Conf.Log.MaxBackups,
		MaxAge:     global.Conf.Log.MaxAge,
		Compress:   global.Conf.Log.Compress,
	}
	// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	return zapcore.AddSync(lumberJackLogger)
}
