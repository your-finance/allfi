// =================================================================================
// 资产健康评分服务接口定义
// 提供投资组合多维度健康评估能力，量化投资风险
// =================================================================================

package service

import (
	"context"

	"your-finance/allfi/internal/app/health_score/model"
)

// IHealthScore 资产健康评分服务接口
type IHealthScore interface {
	// GetHealthScore 获取资产健康评分
	// currency: 计价货币（USD/CNY 等）
	// 返回综合评分、各维度得分和改善建议
	GetHealthScore(ctx context.Context, currency string) (*model.HealthScoreResult, error)
}

// localHealthScore 健康评分服务实例（延迟注入）
var localHealthScore IHealthScore

// HealthScore 获取健康评分服务实例
// 如果服务未注册，会触发 panic
func HealthScore() IHealthScore {
	if localHealthScore == nil {
		panic("IHealthScore 服务未注册，请检查 logic/health_score 包的 init 函数")
	}
	return localHealthScore
}

// RegisterHealthScore 注册健康评分服务实现
// 由 logic 层在 init 函数中调用
func RegisterHealthScore(i IHealthScore) {
	localHealthScore = i
}
