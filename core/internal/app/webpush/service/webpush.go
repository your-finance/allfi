// Package service WebPush 推送模块 - Service 层接口定义
package service

import (
	"context"
)

// IWebpush WebPush 推送服务接口
type IWebpush interface {
	// GetVAPIDPublicKey 获取 VAPID 公钥（供前端使用）
	GetVAPIDPublicKey(ctx context.Context) (string, error)

	// Subscribe 订阅 WebPush 推送
	// endpoint: 推送服务端点 URL
	// p256dh: P-256 Diffie-Hellman 公钥
	// auth: 认证密钥
	Subscribe(ctx context.Context, endpoint string, p256dh string, auth string) error

	// Unsubscribe 取消订阅 WebPush 推送
	// endpoint: 推送服务端点 URL
	Unsubscribe(ctx context.Context, endpoint string) error

	// SendPush 向指定用户发送 WebPush 通知（内部方法）
	// userID: 目标用户ID
	// title: 通知标题
	// body: 通知正文
	SendPush(ctx context.Context, userID int, title string, body string) error
}

var localWebpush IWebpush

// Webpush 获取 WebPush 推送服务实例
func Webpush() IWebpush {
	if localWebpush == nil {
		panic("IWebpush 服务未注册，请检查 logic/webpush 包的 init 函数")
	}
	return localWebpush
}

// RegisterWebpush 注册 WebPush 推送服务实现
// 由 logic 层在 init 函数中调用
func RegisterWebpush(i IWebpush) {
	localWebpush = i
}
