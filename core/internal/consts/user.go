// Package consts 全局常量定义
package consts

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	// DefaultUserID 默认用户 ID（单用户模式）
	DefaultUserID = 1
)

// GetUserID 从请求上下文中获取当前用户 ID
// 优先从 GoFrame 请求参数中获取（由 Auth 中间件写入）
// 如果获取不到（如定时任务、非 HTTP 请求等场景），返回默认值 1（单用户模式）
func GetUserID(ctx context.Context) int {
	// 尝试从 GoFrame HTTP 请求上下文中获取
	r := ghttp.RequestFromCtx(ctx)
	if r != nil {
		userID := r.GetParam("user_id").Int()
		if userID > 0 {
			return userID
		}
	}
	// 单用户模式默认返回 1
	return DefaultUserID
}
