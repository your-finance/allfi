// Package service 价格预警模块 - 服务接口定义
// 定义价格预警的增删改查和触发检查
package service

import (
	"context"

	priceAlertApi "your-finance/allfi/api/v1/price_alert"
)

// IPriceAlert 价格预警服务接口
type IPriceAlert interface {
	// CreateAlert 创建价格预警
	CreateAlert(ctx context.Context, userID int, req *priceAlertApi.CreateReq) (*priceAlertApi.AlertItem, error)

	// GetAlerts 获取预警列表
	GetAlerts(ctx context.Context, userID int) ([]priceAlertApi.AlertItem, error)

	// UpdateAlert 更新预警（暂停/恢复）
	UpdateAlert(ctx context.Context, req *priceAlertApi.UpdateReq) (*priceAlertApi.AlertItem, error)

	// DeleteAlert 删除预警
	DeleteAlert(ctx context.Context, alertID int) error

	// CheckAlerts 检查触发条件（内部方法，cron 调用）
	CheckAlerts(ctx context.Context) error
}

var localPriceAlert IPriceAlert

// PriceAlert 获取价格预警服务实例
func PriceAlert() IPriceAlert {
	if localPriceAlert == nil {
		panic("IPriceAlert 服务未注册，请检查 logic/price_alert 包的 init 函数")
	}
	return localPriceAlert
}

// RegisterPriceAlert 注册价格预警服务实现
// 由 logic 层在 init 函数中调用
func RegisterPriceAlert(i IPriceAlert) {
	localPriceAlert = i
}
