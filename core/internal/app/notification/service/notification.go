// Package service 通知模块 - 服务接口定义
// 定义通知管理的所有服务方法（列表、已读、偏好设置等）
package service

import (
	"context"

	notificationApi "your-finance/allfi/api/v1/notification"
)

// INotification 通知服务接口
type INotification interface {
	// ListNotifications 分页获取通知列表
	ListNotifications(ctx context.Context, userID int, page, pageSize int) ([]notificationApi.NotificationItem, *notificationApi.PaginationInfo, error)

	// GetUnreadCount 获取未读通知数量
	GetUnreadCount(ctx context.Context, userID int) (int64, error)

	// MarkRead 标记单条通知为已读
	MarkRead(ctx context.Context, notifID int) error

	// MarkAllRead 标记所有通知为已读
	MarkAllRead(ctx context.Context, userID int) error

	// GetPreferences 获取通知偏好设置
	GetPreferences(ctx context.Context, userID int) (*notificationApi.PreferenceItem, error)

	// UpdatePreferences 更新通知偏好设置
	UpdatePreferences(ctx context.Context, userID int, req *notificationApi.UpdatePreferencesReq) (*notificationApi.PreferenceItem, error)

	// Send 创建并发送通知（内部方法，不暴露路由）
	Send(ctx context.Context, userID int, notifType, title, message string) error
}

var localNotification INotification

// Notification 获取通知服务实例
func Notification() INotification {
	if localNotification == nil {
		panic("INotification 服务未注册，请检查 logic/notification 包的 init 函数")
	}
	return localNotification
}

// RegisterNotification 注册通知服务实现
// 由 logic 层在 init 函数中调用
func RegisterNotification(i INotification) {
	localNotification = i
}
