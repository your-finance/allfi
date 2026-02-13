// Package webpush WebPush 推送 API 定义
// 提供 VAPID 公钥获取、订阅、取消订阅接口
package webpush

import "github.com/gogf/gf/v2/frame/g"

// GetVAPIDReq 获取 VAPID 公钥请求
type GetVAPIDReq struct {
	g.Meta `path:"/notifications/webpush/vapid" method:"get" summary:"获取 VAPID 公钥" tags:"WebPush"`
}

// GetVAPIDRes 获取 VAPID 公钥响应
type GetVAPIDRes struct {
	VAPIDPublicKey string `json:"vapid_public_key" dc:"VAPID 公钥"`
}

// SubscribeReq 订阅 WebPush 请求
type SubscribeReq struct {
	g.Meta   `path:"/notifications/webpush/subscribe" method:"post" summary:"订阅 WebPush 推送" tags:"WebPush"`
	Endpoint string         `json:"endpoint" v:"required" dc:"推送服务端点 URL"`
	Keys     *SubscribeKeys `json:"keys" v:"required" dc:"加密密钥"`
}

// SubscribeKeys WebPush 订阅密钥
type SubscribeKeys struct {
	P256dh string `json:"p256dh" v:"required" dc:"P-256 Diffie-Hellman 公钥"`
	Auth   string `json:"auth" v:"required" dc:"认证密钥"`
}

// SubscribeRes 订阅 WebPush 响应
type SubscribeRes struct{}

// UnsubscribeReq 取消订阅 WebPush 请求
type UnsubscribeReq struct {
	g.Meta   `path:"/notifications/webpush/unsubscribe" method:"post" summary:"取消订阅 WebPush 推送" tags:"WebPush"`
	Endpoint string `json:"endpoint" v:"required" dc:"推送服务端点 URL"`
}

// UnsubscribeRes 取消订阅 WebPush 响应
type UnsubscribeRes struct{}
