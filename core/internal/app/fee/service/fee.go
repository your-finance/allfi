// =================================================================================
// 费用分析服务接口定义
// 提供交易费用（手续费 + Gas 费 + 提现费）分析能力
// =================================================================================

package service

import (
	"context"

	"your-finance/allfi/internal/app/fee/model"
)

// IFee 费用分析服务接口
type IFee interface {
	// GetFeeAnalytics 获取费用分析
	// userID: 用户 ID
	// period: 时间范围（7d/30d/90d/1y）
	// currency: 计价货币
	// 返回费用分析结果，包含总费用、分类明细、每日趋势
	GetFeeAnalytics(ctx context.Context, userID uint, period string, currency string) (*model.FeeAnalytics, error)
}

// localFee 费用分析服务实例（延迟注入）
var localFee IFee

// Fee 获取费用分析服务实例
// 如果服务未注册，会触发 panic
func Fee() IFee {
	if localFee == nil {
		panic("IFee 服务未注册，请检查 logic/fee 包的 init 函数")
	}
	return localFee
}

// RegisterFee 注册费用分析服务实现
// 由 logic 层在 init 函数中调用
func RegisterFee(i IFee) {
	localFee = i
}
