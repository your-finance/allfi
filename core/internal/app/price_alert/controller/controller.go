// Package controller 价格预警模块 - 路由和控制器
// 绑定价格预警 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	priceAlertApi "your-finance/allfi/api/v1/price_alert"
	"your-finance/allfi/internal/app/price_alert/service"
	"your-finance/allfi/internal/consts"
)

// Controller 价格预警控制器
type Controller struct{}

// Register 注册价格预警模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// Create 创建价格预警
func (c *Controller) Create(ctx context.Context, req *priceAlertApi.CreateReq) (res *priceAlertApi.CreateRes, err error) {
	userID := consts.GetUserID(ctx)

	alert, err := service.PriceAlert().CreateAlert(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &priceAlertApi.CreateRes{
		Alert: alert,
	}, nil
}

// List 获取价格预警列表
func (c *Controller) List(ctx context.Context, req *priceAlertApi.ListReq) (res *priceAlertApi.ListRes, err error) {
	userID := consts.GetUserID(ctx)

	alerts, err := service.PriceAlert().GetAlerts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &priceAlertApi.ListRes{
		Alerts: alerts,
	}, nil
}

// Update 更新价格预警
func (c *Controller) Update(ctx context.Context, req *priceAlertApi.UpdateReq) (res *priceAlertApi.UpdateRes, err error) {
	alert, err := service.PriceAlert().UpdateAlert(ctx, req)
	if err != nil {
		return nil, err
	}

	return &priceAlertApi.UpdateRes{
		Alert: alert,
	}, nil
}

// Delete 删除价格预警
func (c *Controller) Delete(ctx context.Context, req *priceAlertApi.DeleteReq) (res *priceAlertApi.DeleteRes, err error) {
	err = service.PriceAlert().DeleteAlert(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &priceAlertApi.DeleteRes{}, nil
}
