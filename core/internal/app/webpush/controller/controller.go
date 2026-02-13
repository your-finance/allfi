// Package controller WebPush 推送模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	wpApi "your-finance/allfi/api/v1/webpush"
	"your-finance/allfi/internal/app/webpush/service"
)

// WebpushController WebPush 推送控制器
type WebpushController struct{}

// GetVAPID 获取 VAPID 公钥
func (c *WebpushController) GetVAPID(ctx context.Context, req *wpApi.GetVAPIDReq) (res *wpApi.GetVAPIDRes, err error) {
	// 调用服务层获取 VAPID 公钥
	key, err := service.Webpush().GetVAPIDPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	res = &wpApi.GetVAPIDRes{
		VAPIDPublicKey: key,
	}
	return res, nil
}

// Subscribe 订阅 WebPush
func (c *WebpushController) Subscribe(ctx context.Context, req *wpApi.SubscribeReq) (res *wpApi.SubscribeRes, err error) {
	// 调用服务层订阅
	err = service.Webpush().Subscribe(ctx, req.Endpoint, req.Keys.P256dh, req.Keys.Auth)
	if err != nil {
		return nil, err
	}

	return &wpApi.SubscribeRes{}, nil
}

// Unsubscribe 取消订阅 WebPush
func (c *WebpushController) Unsubscribe(ctx context.Context, req *wpApi.UnsubscribeReq) (res *wpApi.UnsubscribeRes, err error) {
	// 调用服务层取消订阅
	err = service.Webpush().Unsubscribe(ctx, req.Endpoint)
	if err != nil {
		return nil, err
	}

	return &wpApi.UnsubscribeRes{}, nil
}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&WebpushController{})
}
