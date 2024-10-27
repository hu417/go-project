package middleware

import (
	"bytes"
	"io"
	"time"

	"gin-rbac/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// 定义请求体最大大小
const maxRequestBodySize = 1024 * 1024 // 1 MB

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 记录请求日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipSpecialRoutes(c) {
			c.Next()
			return
		}

		// 开始计时 time.Now().UTC() 来获取 UTC 时间，这有助于统一时间记录
		start := time.Now().UTC()

		// 生成请求 ID
		requestId := c.GetString("request_id")
		if requestId == "" {
			requestId = uuid.New().String()
			c.Set("request_id", requestId)
		}

		// 获取 response 内容
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, maxRequestBodySize))
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		// 开始记录日志的逻辑
		cost := time.Since(start)
		responStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.String("request_id", requestId),
			zap.Int("status", responStatus),
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", cost.String()),
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			// 请求的内容
			logFields = append(logFields, zap.String("Request Body", string(requestBody)))

			// 响应的内容
			logFields = append(logFields, zap.String("Response Body", w.body.String()))
		}

		if responStatus > 400 && responStatus <= 499 {
			// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意

			global.Log.Warn("HTTP Warning "+cast.ToString(responStatus), logFields)
		} else if responStatus >= 500 && responStatus <= 599 {
			// 除了内部错误，记录 error
			global.Log.Error("HTTP Error "+cast.ToString(responStatus), logFields)
		} else {
			global.Log.Debug("HTTP Access Log", logFields)
		}
	}
}
