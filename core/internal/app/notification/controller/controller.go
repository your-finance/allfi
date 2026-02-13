// Package controller 通知模块 - 控制器
// 绑定通知 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	notificationApi "your-finance/allfi/api/v1/notification"
	"your-finance/allfi/internal/app/notification/service"
	"your-finance/allfi/internal/consts"
)

// Controller 通知控制器
type Controller struct{}

// Register 注册通知模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// List 获取通知列表
func (c *Controller) List(ctx context.Context, req *notificationApi.ListReq) (res *notificationApi.ListRes, err error) {
	userID := consts.GetUserID(ctx)

	list, pagination, err := service.Notification().ListNotifications(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &notificationApi.ListRes{
		List:       list,
		Pagination: pagination,
	}, nil
}

// GetUnreadCount 获取未读通知数量
func (c *Controller) GetUnreadCount(ctx context.Context, req *notificationApi.GetUnreadCountReq) (res *notificationApi.GetUnreadCountRes, err error) {
	userID := consts.GetUserID(ctx)

	count, err := service.Notification().GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &notificationApi.GetUnreadCountRes{
		Count: count,
	}, nil
}

// MarkRead 标记通知为已读
func (c *Controller) MarkRead(ctx context.Context, req *notificationApi.MarkReadReq) (res *notificationApi.MarkReadRes, err error) {
	err = service.Notification().MarkRead(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &notificationApi.MarkReadRes{}, nil
}

// MarkAllRead 标记所有通知为已读
func (c *Controller) MarkAllRead(ctx context.Context, req *notificationApi.MarkAllReadReq) (res *notificationApi.MarkAllReadRes, err error) {
	userID := consts.GetUserID(ctx)

	err = service.Notification().MarkAllRead(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &notificationApi.MarkAllReadRes{}, nil
}

// GetPreferences 获取通知偏好设置
func (c *Controller) GetPreferences(ctx context.Context, req *notificationApi.GetPreferencesReq) (res *notificationApi.GetPreferencesRes, err error) {
	userID := consts.GetUserID(ctx)

	preferences, err := service.Notification().GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &notificationApi.GetPreferencesRes{
		Preferences: preferences,
	}, nil
}

// UpdatePreferences 更新通知偏好设置
func (c *Controller) UpdatePreferences(ctx context.Context, req *notificationApi.UpdatePreferencesReq) (res *notificationApi.UpdatePreferencesRes, err error) {
	userID := consts.GetUserID(ctx)

	preferences, err := service.Notification().UpdatePreferences(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &notificationApi.UpdatePreferencesRes{
		Preferences: preferences,
	}, nil
}
