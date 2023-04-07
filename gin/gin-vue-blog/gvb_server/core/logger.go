package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gvb_server/global"
	"gvb_server/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化日志
func InitLogger() *zap.SugaredLogger {

	// 判断文件是否存在，若不存在，则创建
	if ok, _ := utils.PathExist(global.Config.Logger.Director); !ok {
		fmt.Printf("create %v directory \n", global.Config.Logger.Director)
		_ = os.Mkdir(global.Config.Logger.Director, os.ModePerm)
	}

	//====== 日志輸出級別 ======
	// 调试级别
	//debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
	// 	return lev == zap.DebugLevel
	// })
	// // 日志级别
	// infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
	// 	return lev == zap.InfoLevel
	// })
	// // 警告级别
	// warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
	// 	return lev == zap.WarnLevel
	// })
	// // 错误级别
	// errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
	// 	return lev >= zap.ErrorLevel
	// })

	stSeparator := string(filepath.Separator) // 分隔符
	stRootDir, _ := os.Getwd()                // 工作路径
	// stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly) + ".log" // DateOnly表示日期
	// 项目更目录 + 分隔符   + 文件名  + 分隔符      + 日志  + 后缀
	// stLogFilePath := stRootDir + stSeparator + global.Config.Logger.Director + stSeparator + level + time.Now().Format("2006-01-02") + ".log"
	stLogFilePath := stRootDir + stSeparator + global.Config.Logger.Director

	// 实例化Logger
	cores := [...]zapcore.Core{
		// zapcore.NewCore(getEncoder(), zapcore.AddSync(os.Stdout), debugPriority),         //打印到控制台
		getEncoderCore(fmt.Sprintf("%s/server_%s.log", stLogFilePath, global.Config.Logger.Level), getLevelPriority()), // 写入到level级别文件中
	}

	logger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	// 是否显示行号
	if global.Config.Logger.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger.Sugar()
}

// getEncoderCore 获取Encoder的zapcore.Core,
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	//写入到文件
	writer := getWriteSyncer(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 获取日志输出级别
func getLevelPriority() zapcore.LevelEnabler {
	switch global.Config.Logger.Level {
	case "debug", "Debug":
		return zap.DebugLevel
	case "info", "Info":
		return zap.InfoLevel
	case "warn", "Warn":
		return zap.WarnLevel
	case "error", "Error":
		return zap.ErrorLevel
	case "dpanic", "DPanic":
		return zap.DPanicLevel
	case "panic", "Panic":
		return zap.PanicLevel
	case "fatal", "Fatal":
		return zap.FatalLevel
	}
	return zap.InfoLevel
}

// getEncoder 获取zapcore.Encoder 日志格式
func getEncoder() zapcore.Encoder {
	if global.Config.Logger.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// 日志输出路径: 文件、控制台、双向输出
func getWriteSyncer(file string) zapcore.WriteSyncer {

	luberjackSyncer := &lumberjack.Logger{
		Filename:   file,                            // 文件路径
		MaxSize:    global.Config.Logger.MaxSize,    // 日志文件最大的尺寸(M), 超限后开始自动分割
		MaxBackups: global.Config.Logger.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     global.Config.Logger.MaxAge,     // 保留旧文件的最大天数
		Compress:   global.Config.Logger.Compress,   // 是否压缩/归档旧文件
	}

	// 双向输出
	if global.Config.Logger.LogInConsole {
		fileWriter := zapcore.AddSync(luberjackSyncer) // 写入文件
		consoleWriter := zapcore.AddSync(os.Stdout)    // 输出到控制台
		return zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)
	}

	return zapcore.AddSync(luberjackSyncer)
}

// ##### 日志输出格式 #####
// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(global.Config.Logger.Prefix + "[" + "2006/01/02 15:04:05.000" + "]"))
}

// cEncodeLevel 自定义日志级别显示
// func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
// 	enc.AppendString("[" + level.CapitalString() + "]")
// }

// cEncodeCaller 自定义行号显示
// func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
// 	enc.AppendString("[" + caller.TrimmedPath() + "]")
// }

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {

	config = zapcore.EncoderConfig{
		TimeKey:       "time",
		EncodeTime:    CustomTimeEncoder,
		LevelKey:      "level",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		LineEnding:    zapcore.DefaultLineEnding,
		NameKey:       "logger",
		StacktraceKey: global.Config.Logger.StacktraceKey,
		MessageKey:    "message",
		CallerKey:     "caller",
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, // zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeName:   zapcore.FullNameEncoder,
	}

	switch {
	case global.Config.Logger.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case global.Config.Logger.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case global.Config.Logger.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case global.Config.Logger.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}
