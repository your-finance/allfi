package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// CORS 跨域中间件
// 从配置文件读取CORS设置并设置响应头
func CORS(r *ghttp.Request) {
	ctx := r.Context()

	// 从配置读取CORS设置
	corsEnabled := g.Cfg().MustGet(ctx, "cors.enabled", true).Bool()
	if !corsEnabled {
		r.Middleware.Next()
		return
	}

	allowOrigin := g.Cfg().MustGet(ctx, "cors.allowOrigin", "*").String()
	allowMethods := g.Cfg().MustGet(ctx, "cors.allowMethods", "GET,POST,PUT,DELETE,OPTIONS").String()
	allowHeaders := g.Cfg().MustGet(ctx, "cors.allowHeaders", "Origin,Content-Type,Accept,Authorization").String()
	exposeHeaders := g.Cfg().MustGet(ctx, "cors.exposeHeaders", "Content-Length,Content-Type").String()
	allowCredentials := g.Cfg().MustGet(ctx, "cors.allowCredentials", false).Bool()
	maxAge := g.Cfg().MustGet(ctx, "cors.maxAge", 3600).Int()

	// 设置CORS响应头
	r.Response.CORSDefault()
	r.Response.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	r.Response.Header().Set("Access-Control-Allow-Methods", allowMethods)
	r.Response.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	r.Response.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
	r.Response.Header().Set("Access-Control-Max-Age", g.NewVar(maxAge).String())

	if allowCredentials {
		r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// 如果是 OPTIONS 预检请求，直接返回
	if r.Method == "OPTIONS" {
		r.Response.WriteStatus(204)
		return
	}

	// 继续处理请求
	r.Middleware.Next()
}
