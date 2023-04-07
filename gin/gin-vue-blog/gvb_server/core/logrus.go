package core

// import (
// 	"bytes"
// 	"fmt"
// 	//"gvb_server/global"
// 	"os"
// 	"path"

// 	"github.com/sirupsen/logrus"
// )

// // 颜色
// const (
// 	red    = 31
// 	yellow = 33
// 	blue   = 36
// 	groy   = 37
// )

// type LogFormatter struct{}

// // format 实现Formatter(entry *logrus.Entry) ([]byte,error)接口
// func (f *LogFormatter) format(entry *logrus.Entry) ([]byte, error) {
// 	// 根据不同的level展示颜色
// 	var levelColor int
// 	switch entry.Level {
// 	case logrus.DebugLevel, logrus.TraceLevel:
// 		levelColor = groy
// 	case logrus.WarnLevel:
// 		levelColor = yellow
// 	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
// 		levelColor = red
// 	default:
// 		levelColor = blue
// 	}

// 	var b *bytes.Buffer
// 	if entry.Buffer != nil {
// 		b = entry.Buffer
// 	} else {
// 		b = &bytes.Buffer{}
// 	}

// 	// 自定义日期格式
// 	timestamp := entry.Time.Format("2006-01-02 15:04:05")
// 	if entry.HasCaller() {
// 		// 自定义文件路径
// 		funcVal := entry.Caller.Function
// 		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
// 		// 自定义输出格式
// 		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
// 	} else {
// 		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s\n", timestamp, levelColor, entry.Level, entry.Message)
// 	}
// 	return b.Bytes(), nil
// }

// func InitLogrus() *logrus.Logger {
// 	mlog := logrus.New()       // 新建一个实例
// 	mlog.SetOutput(os.Stdout)  // 设置输出类型
// 	mlog.SetReportCaller(true) // 开启返回函数的名和行号
// 	// mlog.SetFormatter(&LogFormatter{}) // 设置自定义的formartter
// 	// level, err := logrus.PanicLevel(global.Config.Logger.Level)
// 	// if err != nil {
// 	// 	level = logrus.InfoLevel
// 	// }

// 	mlog.SetLevel(logrus.DebugLevel) // 设置日志等级
// 	return mlog
// }

// func InitDefaultLogger() {
// 	// 全局log
// 	logrus.SetOutput(os.Stdout)
// 	logrus.SetReportCaller(true)
// 	//logrus.SetFormatter(&LogFormatter{})

// 	logrus.SetLevel(logrus.DebugLevel)
// }
