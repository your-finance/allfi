// Package controller 健康检查模块路由注册
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	healthApi "your-finance/allfi/api/v1/health"
	"your-finance/allfi/internal/app/health/service"
)

// HealthController 健康检查控制器
type HealthController struct{}

// GetHealth 获取服务健康状态
//
// 对应路由: GET /health
// 调用 service.Health().GetHealthStatus 获取状态信息
func (c *HealthController) GetHealth(ctx context.Context, req *healthApi.GetHealthReq) (res *healthApi.GetHealthRes, err error) {
	// 调用 Service 层
	status, err := service.Health().GetHealthStatus(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "获取健康状态失败")
	}

	// 提取版本信息
	version := "unknown"
	if v, ok := status["version"].(string); ok {
		version = v
	}

	// 构造响应（匹配子目录 API 包的 GetHealthRes 定义）
	res = &healthApi.GetHealthRes{
		Status:    status["status"].(string),
		Version:   version,
		Timestamp: time.Now().Unix(),
	}

	return res, nil
}

// Register 注册健康检查路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&HealthController{})
}
