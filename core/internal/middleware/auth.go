// Package middleware 认证中间件
package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"

	"your-finance/allfi/internal/dao"
	"your-finance/allfi/internal/model/entity"
)

// getJWTSecret 从 system_config 表读取 JWT 密钥
// auth logic 在首次设置 PIN 时自动生成并存入数据库
func getJWTSecret(r *ghttp.Request) []byte {
	var config entity.SystemConfig
	err := dao.SystemConfig.Ctx(r.Context()).
		Where(dao.SystemConfig.Columns().ConfigKey, "auth.jwt_secret").
		Scan(&config)
	if err != nil || config.ConfigValue == "" {
		return nil
	}
	// auth logic 使用 base64 编码存储密钥
	decoded, err := base64.StdEncoding.DecodeString(config.ConfigValue)
	if err != nil {
		return []byte(config.ConfigValue)
	}
	return decoded
}

// Auth 认证中间件
// 验证 JWT Token，提取用户信息到请求上下文
func Auth(r *ghttp.Request) {
	// 从 Header 中获取 Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "未提供认证 Token",
		})
		return
	}

	// 验证 Token 格式: Bearer <token>
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "Token 格式错误，应为: Bearer <token>",
		})
		return
	}

	tokenString := parts[1]

	// 从数据库获取 JWT 密钥（与 auth logic 保持一致）
	jwtSecret := getJWTSecret(r)
	if jwtSecret == nil {
		g.Log().Error(r.Context(), "JWT 密钥未初始化，请先设置 PIN")
		r.Response.WriteJsonExit(g.Map{
			"code":    500,
			"message": "服务器未初始化，请先设置 PIN",
		})
		return
	}

	// 解析 Token（兼容 auth logic 的 MapClaims 格式）
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		g.Log().Warning(r.Context(), "Token 验证失败", "error", err)
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "Token 无效或已过期",
		})
		return
	}

	// 单用户模式：设置默认用户参数
	r.SetParam("user_id", 1)
	r.SetParam("username", "allfi-user")

	g.Log().Debug(r.Context(), "认证成功")

	// 继续处理请求
	r.Middleware.Next()
}

// OptionalAuth 可选认证中间件
// 尝试验证 Token，但不强制要求
// 如果 Token 有效，将用户信息存入上下文；如果无效或不存在，继续处理请求
func OptionalAuth(r *ghttp.Request) {
	// 从 Header 中获取 Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Middleware.Next()
		return
	}

	// 验证 Token 格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		r.Middleware.Next()
		return
	}

	tokenString := parts[1]

	// 从数据库获取 JWT 密钥
	jwtSecret := getJWTSecret(r)
	if jwtSecret == nil {
		r.Middleware.Next()
		return
	}

	// 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		r.Middleware.Next()
		return
	}

	// 单用户模式：设置默认用户参数
	r.SetParam("user_id", 1)
	r.SetParam("username", "allfi-user")

	// 继续处理请求
	r.Middleware.Next()
}
