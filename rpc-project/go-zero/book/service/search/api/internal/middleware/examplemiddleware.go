package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExampleMiddleware struct {
}

func NewExampleMiddleware() *ExampleMiddleware {
	return &ExampleMiddleware{}
}

func (m *ExampleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		logx.Info("---- 这是路由中间件 ----")

		// Passthrough to next handler if need
		next(w, r)
	}
}
