package bootstrap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongodb() *mongo.Client {
	// 连接到MongoDB，我这些配置是在其他包里面
	clientOptions := options.Client().ApplyURI("mongodb://mongo:qaz123@localhost:27017/?maxPoolSize=20&w=majority&timeoutMS=5000")
	// 配置日志
	clientOptions.SetLoggerOptions(logger("debug"))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("数据库连接失败！！！")
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("数据库连接成功！")

	return client
}

// 自定义日志
type CustomLogger struct {
	io.Writer
	mu sync.Mutex
}

func (logger *CustomLogger) Info(level int, msg string, _ ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	if options.LogLevel(level+1) == options.LogLevelDebug {
		fmt.Fprintf(logger, "level: %d DEBUG, message: %s\n", level, msg)
	} else {
		fmt.Fprintf(logger, "level: %d INFO, message: %s\n", level, msg)
	}
}
func (logger *CustomLogger) Error(err error, msg string, _ ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	fmt.Fprintf(logger, "error: %v, message: %s\n", err, msg)
}

// 记录日志
func logger(log string) *options.LoggerOptions {

	switch log {
	case "debug":
		return options.Logger().SetMaxDocumentLength(25).SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	case "info":
		buf := bytes.NewBuffer(nil)
		sink := &CustomLogger{Writer: buf}
		loggerOptions := options.
			Logger().
			SetSink(sink).
			SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug).
			SetComponentLevel(options.LogComponentConnection, options.LogLevelDebug)

		return loggerOptions
	default:
		return options.Logger().SetMaxDocumentLength(25).SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	}

}
