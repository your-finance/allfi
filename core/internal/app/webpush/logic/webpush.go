// Package logic WebPush 推送模块 - 业务逻辑实现
package logic

import (
	"context"
	"encoding/json"
	"fmt"

	webpushLib "github.com/SherClockHolmes/webpush-go"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/webpush/dao"
	"your-finance/allfi/internal/app/webpush/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/model/entity"
)

// sWebpush WebPush 推送服务实现
type sWebpush struct{}

// New 创建 WebPush 推送服务实例
func New() service.IWebpush {
	return &sWebpush{}
}

// vapidConfig VAPID 配置（在 main.go 中初始化）
var vapidConfig struct {
	PublicKey  string
	PrivateKey string
	Email      string
}

// SetVAPIDConfig 设置 VAPID 密钥配置
// 由 main.go 在启动时调用
func SetVAPIDConfig(publicKey, privateKey, email string) {
	vapidConfig.PublicKey = publicKey
	vapidConfig.PrivateKey = privateKey
	vapidConfig.Email = email
}

// GetVAPIDPublicKey 获取 VAPID 公钥
func (s *sWebpush) GetVAPIDPublicKey(ctx context.Context) (string, error) {
	if vapidConfig.PublicKey == "" {
		// 尝试从配置文件读取
		publicKey := g.Cfg().MustGet(ctx, "webpush.vapidPublicKey").String()
		if publicKey != "" {
			vapidConfig.PublicKey = publicKey
			vapidConfig.PrivateKey = g.Cfg().MustGet(ctx, "webpush.vapidPrivateKey").String()
			vapidConfig.Email = g.Cfg().MustGet(ctx, "webpush.email", "admin@allfi.local").String()
		}
	}

	if vapidConfig.PublicKey == "" {
		// 自动生成 VAPID 密钥对
		priv, pub, err := webpushLib.GenerateVAPIDKeys()
		if err != nil {
			return "", gerror.Wrap(err, "生成 VAPID 密钥对失败")
		}
		vapidConfig.PublicKey = pub
		vapidConfig.PrivateKey = priv
		if vapidConfig.Email == "" {
			vapidConfig.Email = "admin@allfi.local"
		}
		g.Log().Info(ctx, "已自动生成 VAPID 密钥对", "publicKey", pub[:20]+"...")
	}

	return vapidConfig.PublicKey, nil
}

// Subscribe 订阅 WebPush 推送
//
// 业务逻辑:
// 1. 获取当前用户ID
// 2. 存储 endpoint + p256dh + auth 到 web_push_subscriptions 表
// 3. 如果已存在相同 endpoint 则更新
func (s *sWebpush) Subscribe(ctx context.Context, endpoint string, p256dh string, auth string) error {
	if endpoint == "" || p256dh == "" || auth == "" {
		return gerror.New("endpoint、p256dh、auth 不能为空")
	}

	userID := consts.GetUserID(ctx)

	// 检查是否已存在该 endpoint 的订阅
	count, err := dao.WebPushSubscriptions.Ctx(ctx).
		Where(dao.WebPushSubscriptions.Columns().Endpoint, endpoint).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询订阅记录失败")
	}

	if count > 0 {
		// 更新已有订阅
		_, err = dao.WebPushSubscriptions.Ctx(ctx).
			Where(dao.WebPushSubscriptions.Columns().Endpoint, endpoint).
			Update(g.Map{
				dao.WebPushSubscriptions.Columns().P256Dh: p256dh,
				dao.WebPushSubscriptions.Columns().Auth:   auth,
				dao.WebPushSubscriptions.Columns().UserId: userID,
			})
		if err != nil {
			return gerror.Wrap(err, "更新订阅记录失败")
		}
		g.Log().Info(ctx, "更新 WebPush 订阅成功", "endpoint", endpoint[:50]+"...")
	} else {
		// 创建新订阅
		sub := &entity.WebPushSubscriptions{
			UserId:   userID,
			Endpoint: endpoint,
			P256Dh:   p256dh,
			Auth:     auth,
		}
		_, err = dao.WebPushSubscriptions.Ctx(ctx).Insert(sub)
		if err != nil {
			return gerror.Wrap(err, "创建订阅记录失败")
		}
		g.Log().Info(ctx, "创建 WebPush 订阅成功", "userId", userID)
	}

	return nil
}

// Unsubscribe 取消订阅 WebPush 推送
//
// 业务逻辑:
// 根据 endpoint 删除订阅记录
func (s *sWebpush) Unsubscribe(ctx context.Context, endpoint string) error {
	if endpoint == "" {
		return gerror.New("endpoint 不能为空")
	}

	_, err := dao.WebPushSubscriptions.Ctx(ctx).
		Where(dao.WebPushSubscriptions.Columns().Endpoint, endpoint).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "删除订阅记录失败")
	}

	g.Log().Info(ctx, "取消 WebPush 订阅成功")
	return nil
}

// SendPush 向指定用户发送 WebPush 通知
//
// 业务逻辑:
// 1. 查询用户的订阅信息
// 2. 构建推送载荷（JSON 格式）
// 3. 使用 web-push 库发送通知
func (s *sWebpush) SendPush(ctx context.Context, userID int, title string, body string) error {
	// 查询用户的订阅信息
	var sub entity.WebPushSubscriptions
	err := dao.WebPushSubscriptions.Ctx(ctx).
		Where(dao.WebPushSubscriptions.Columns().UserId, userID).
		Scan(&sub)
	if err != nil {
		return gerror.Wrapf(err, "未找到用户 %d 的 WebPush 订阅", userID)
	}

	if sub.Endpoint == "" {
		return gerror.Newf("用户 %d 无有效的 WebPush 订阅", userID)
	}

	// 确保 VAPID 密钥可用
	_, err = s.GetVAPIDPublicKey(ctx)
	if err != nil {
		return err
	}

	// 构建推送载荷
	payload := map[string]string{
		"title": title,
		"body":  body,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return gerror.Wrap(err, "序列化推送数据失败")
	}

	// 构建 webpush 订阅对象
	subscription := &webpushLib.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpushLib.Keys{
			P256dh: sub.P256Dh,
			Auth:   sub.Auth,
		},
	}

	// 发送推送
	resp, err := webpushLib.SendNotification(payloadJSON, subscription, &webpushLib.Options{
		Subscriber:      vapidConfig.Email,
		VAPIDPublicKey:  vapidConfig.PublicKey,
		VAPIDPrivateKey: vapidConfig.PrivateKey,
		TTL:             60, // 60 秒 TTL
	})
	if err != nil {
		return gerror.Wrap(err, "发送 WebPush 失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("WebPush 推送服务返回错误: %d", resp.StatusCode)
	}

	g.Log().Info(ctx, "发送 WebPush 成功",
		"userId", userID,
		"title", title,
		"statusCode", resp.StatusCode,
	)

	return nil
}
