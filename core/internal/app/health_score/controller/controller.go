// Package controller 资产健康评分模块控制器
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	healthScoreApi "your-finance/allfi/api/v1/health_score"
	"your-finance/allfi/internal/app/health_score/service"
)

// HealthScoreController 资产健康评分控制器
type HealthScoreController struct{}

// Get 获取资产健康评分
//
// 对应路由: GET /portfolio/health
// 查询参数: currency — 计价货币，默认 USD
func (c *HealthScoreController) Get(ctx context.Context, req *healthScoreApi.GetReq) (res *healthScoreApi.GetRes, err error) {
	// 调用 Service 层
	result, err := service.HealthScore().GetHealthScore(ctx, req.Currency)
	if err != nil {
		return nil, gerror.Wrap(err, "获取健康评分失败")
	}

	// 将业务 DTO 转换为 API 响应
	var details []healthScoreApi.ScoreDetail
	for _, d := range result.Details {
		details = append(details, healthScoreApi.ScoreDetail{
			Category:    d.Category,
			Score:       d.Score,
			Weight:      d.Weight,
			Description: d.Description,
			Suggestion:  d.Suggestion,
		})
	}

	res = &healthScoreApi.GetRes{
		OverallScore: result.OverallScore,
		Level:        result.Level,
		Details:      details,
		Currency:     result.Currency,
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	return res, nil
}

// Register 注册资产健康评分路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&HealthScoreController{})
}
