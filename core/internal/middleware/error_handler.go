package middleware

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"your-finance/allfi/utility/response"
)

// ErrorHandler 错误处理中间件
// 捕获panic和业务错误，返回统一的错误响应
func ErrorHandler(r *ghttp.Request) {
	ctx := r.Context()

	// defer 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			// 记录 panic 日志（包含堆栈）
			g.Log().Errorf(ctx, "发生panic: %+v", err)

			// 返回500错误
			response.InternalError(r, "服务器内部错误，请稍后重试")
		}
	}()

	// 继续处理请求
	r.Middleware.Next()

	// 检查是否有错误
	if err := r.GetError(); err != nil {
		// 清除已设置的错误
		r.SetError(nil)

		// 解析错误类型并返回对应响应
		handleError(r, err)
	}
}

// handleError 处理不同类型的错误
func handleError(r *ghttp.Request, err error) {
	ctx := r.Context()

	// 记录错误日志
	g.Log().Error(ctx, "请求错误",
		"error", err,
		"path", r.URL.Path,
		"method", r.Method,
	)

	// 根据错误类型返回响应
	switch {
	case gerror.HasStack(err):
		// GoFrame 的错误（包含堆栈）
		code := gerror.Code(err)
		if code.Code() != -1 {
			// 有明确错误码
			response.Error(r, code.Code(), code.Message())
		} else {
			// 无错误码，使用错误信息
			response.Error(r, response.CodeInternalError, err.Error())
		}
	default:
		// 普通错误
		response.InternalError(r, err.Error())
	}
}
