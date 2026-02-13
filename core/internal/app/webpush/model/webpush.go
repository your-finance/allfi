// Package model WebPush 推送模块 - 数据传输对象
package model

// SubscribeInput 订阅请求输入
type SubscribeInput struct {
	Endpoint string `json:"endpoint" v:"required" dc:"推送服务端点 URL"`
	P256dh   string `json:"p256dh" v:"required" dc:"P-256 Diffie-Hellman 公钥"`
	Auth     string `json:"auth" v:"required" dc:"认证密钥"`
}

// UnsubscribeInput 取消订阅请求输入
type UnsubscribeInput struct {
	Endpoint string `json:"endpoint" v:"required" dc:"推送服务端点 URL"`
}

// PushMessage 推送消息
type PushMessage struct {
	Title string `json:"title" dc:"通知标题"`
	Body  string `json:"body" dc:"通知正文"`
}
