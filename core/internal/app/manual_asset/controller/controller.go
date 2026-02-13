// Package controller 手动资产模块 - 控制器
// 绑定手动资产 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	manualAssetApi "your-finance/allfi/api/v1/manual_asset"
	"your-finance/allfi/internal/app/manual_asset/service"
	"your-finance/allfi/internal/consts"
)

// Controller 手动资产控制器
type Controller struct{}

// Register 注册手动资产模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// List 获取手动资产列表
func (c *Controller) List(ctx context.Context, req *manualAssetApi.ListReq) (res *manualAssetApi.ListRes, err error) {
	userID := consts.GetUserID(ctx)

	assets, err := service.ManualAsset().ListManualAssets(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &manualAssetApi.ListRes{
		Assets: assets,
	}, nil
}

// Create 添加手动资产
func (c *Controller) Create(ctx context.Context, req *manualAssetApi.CreateReq) (res *manualAssetApi.CreateRes, err error) {
	userID := consts.GetUserID(ctx)

	asset, err := service.ManualAsset().CreateManualAsset(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &manualAssetApi.CreateRes{
		Asset: asset,
	}, nil
}

// Update 更新手动资产
func (c *Controller) Update(ctx context.Context, req *manualAssetApi.UpdateReq) (res *manualAssetApi.UpdateRes, err error) {
	asset, err := service.ManualAsset().UpdateManualAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return &manualAssetApi.UpdateRes{
		Asset: asset,
	}, nil
}

// Delete 删除手动资产
func (c *Controller) Delete(ctx context.Context, req *manualAssetApi.DeleteReq) (res *manualAssetApi.DeleteRes, err error) {
	err = service.ManualAsset().DeleteManualAsset(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &manualAssetApi.DeleteRes{}, nil
}
