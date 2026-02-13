// Package middleware 全局中间件注册
// 提供统一的中间件注册入口，将所有全局中间件绑定到 GoFrame 服务器实例
package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 注册全局中间件
// 按照执行顺序绑定所有默认中间件到服务器实例：
//  1. CORS - 跨域请求处理（最先执行，确保预检请求被正确响应）
//  2. Context - 上下文注入（生成 TraceID 等请求级元数据）
//  3. Logger - 请求日志（记录请求开始/结束、耗时等信息）
//  4. ErrorHandler - 错误处理（捕获 panic 和业务错误）
//  5. MiddlewareHandlerResponse - 统一响应包装（将 Controller 返回值包装为标准 JSON 格式）
//
// 注意：Auth 认证中间件不在此处注册，而是按路由组单独绑定，
// 因为部分路由（如 /auth/*、/health）不需要认证
func Register(s *ghttp.Server) {
	s.BindMiddlewareDefault(
		CORS,
		Context,
		Logger,
		ErrorHandler,
		MiddlewareHandlerResponse,
	)
}
