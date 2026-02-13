// Package logic 通知业务逻辑
// 实现通知的分页查询、已读标记、偏好设置和异步推送
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	notificationApi "your-finance/allfi/api/v1/notification"
	"your-finance/allfi/internal/app/notification/dao"
	"your-finance/allfi/internal/app/notification/service"
	webpushService "your-finance/allfi/internal/app/webpush/service"
	"your-finance/allfi/internal/model/entity"
)

// sNotification 通知服务实现
type sNotification struct{}

// New 创建通知服务实例
func New() service.INotification {
	return &sNotification{}
}

// ListNotifications 分页获取通知列表
func (s *sNotification) ListNotifications(ctx context.Context, userID int, page, pageSize int) ([]notificationApi.NotificationItem, *notificationApi.PaginationInfo, error) {
	// 参数校验
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var list []entity.Notifications
	var total int
	err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().UserId, userID).
		WhereNull(dao.Notifications.Columns().DeletedAt).
		OrderDesc(dao.Notifications.Columns().CreatedAt).
		Page(page, pageSize).
		ScanAndCount(&list, &total, true)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询通知列表失败")
	}

	// 转换为 API 响应格式
	items := make([]notificationApi.NotificationItem, 0, len(list))
	for _, n := range list {
		items = append(items, s.toNotificationItem(&n))
	}

	// 构建分页信息
	totalInt64 := int64(total)
	totalPages := int(totalInt64+int64(pageSize)-1) / pageSize
	pagination := &notificationApi.PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		Total:      totalInt64,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return items, pagination, nil
}

// GetUnreadCount 获取未读通知数量
func (s *sNotification) GetUnreadCount(ctx context.Context, userID int) (int64, error) {
	count, err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().UserId, userID).
		Where(dao.Notifications.Columns().IsRead, 0).
		WhereNull(dao.Notifications.Columns().DeletedAt).
		Count()
	if err != nil {
		return 0, gerror.Wrap(err, "查询未读数量失败")
	}
	return int64(count), nil
}

// MarkRead 标记单条通知为已读
func (s *sNotification) MarkRead(ctx context.Context, notifID int) error {
	_, err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().Id, notifID).
		Data(g.Map{
			dao.Notifications.Columns().IsRead:    1,
			dao.Notifications.Columns().UpdatedAt: gtime.Now(),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "标记已读失败")
	}
	return nil
}

// MarkAllRead 标记所有通知为已读
func (s *sNotification) MarkAllRead(ctx context.Context, userID int) error {
	_, err := dao.Notifications.Ctx(ctx).
		Where(dao.Notifications.Columns().UserId, userID).
		Where(dao.Notifications.Columns().IsRead, 0).
		WhereNull(dao.Notifications.Columns().DeletedAt).
		Data(g.Map{
			dao.Notifications.Columns().IsRead:    1,
			dao.Notifications.Columns().UpdatedAt: gtime.Now(),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "标记全部已读失败")
	}
	return nil
}

// GetPreferences 获取通知偏好设置
// 如果不存在则返回默认值
func (s *sNotification) GetPreferences(ctx context.Context, userID int) (*notificationApi.PreferenceItem, error) {
	var pref entity.NotificationPreferences
	err := dao.NotificationPreferences.Ctx(ctx).
		Where(dao.NotificationPreferences.Columns().UserId, userID).
		WhereNull(dao.NotificationPreferences.Columns().DeletedAt).
		Scan(&pref)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, gerror.Wrap(err, "查询通知偏好失败")
	}

	// 如果不存在，返回默认偏好
	if err != nil || pref.Id == 0 {
		return &notificationApi.PreferenceItem{
			EmailEnabled:   false,
			PushEnabled:    true,
			PriceAlert:     true,
			PortfolioAlert: true,
			SystemNotice:   true,
		}, nil
	}

	return &notificationApi.PreferenceItem{
		EmailEnabled:   pref.EnableDailyDigest == 1,
		PushEnabled:    pref.EnablePriceAlert == 1,
		PriceAlert:     pref.EnablePriceAlert == 1,
		PortfolioAlert: pref.EnableAssetAlert == 1,
		SystemNotice:   true,
	}, nil
}

// UpdatePreferences 更新通知偏好设置
func (s *sNotification) UpdatePreferences(ctx context.Context, userID int, req *notificationApi.UpdatePreferencesReq) (*notificationApi.PreferenceItem, error) {
	// 检查是否已存在偏好设置
	var existing entity.NotificationPreferences
	err := dao.NotificationPreferences.Ctx(ctx).
		Where(dao.NotificationPreferences.Columns().UserId, userID).
		WhereNull(dao.NotificationPreferences.Columns().DeletedAt).
		Scan(&existing)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, gerror.Wrap(err, "查询通知偏好失败")
	}

	now := gtime.Now()

	if err != nil || existing.Id == 0 {
		// 不存在则创建
		data := g.Map{
			dao.NotificationPreferences.Columns().UserId:    userID,
			dao.NotificationPreferences.Columns().CreatedAt: now,
			dao.NotificationPreferences.Columns().UpdatedAt: now,
		}
		if req.EmailEnabled != nil {
			data[dao.NotificationPreferences.Columns().EnableDailyDigest] = boolToFloat(*req.EmailEnabled)
		}
		if req.PushEnabled != nil {
			data[dao.NotificationPreferences.Columns().EnablePriceAlert] = boolToFloat(*req.PushEnabled)
		}
		if req.PriceAlert != nil {
			data[dao.NotificationPreferences.Columns().EnablePriceAlert] = boolToFloat(*req.PriceAlert)
		}
		if req.PortfolioAlert != nil {
			data[dao.NotificationPreferences.Columns().EnableAssetAlert] = boolToFloat(*req.PortfolioAlert)
		}

		_, err = dao.NotificationPreferences.Ctx(ctx).Data(data).Insert()
		if err != nil {
			return nil, gerror.Wrap(err, "创建通知偏好失败")
		}
	} else {
		// 存在则更新
		updateData := g.Map{
			dao.NotificationPreferences.Columns().UpdatedAt: now,
		}
		if req.EmailEnabled != nil {
			updateData[dao.NotificationPreferences.Columns().EnableDailyDigest] = boolToFloat(*req.EmailEnabled)
		}
		if req.PushEnabled != nil {
			updateData[dao.NotificationPreferences.Columns().EnablePriceAlert] = boolToFloat(*req.PushEnabled)
		}
		if req.PriceAlert != nil {
			updateData[dao.NotificationPreferences.Columns().EnablePriceAlert] = boolToFloat(*req.PriceAlert)
		}
		if req.PortfolioAlert != nil {
			updateData[dao.NotificationPreferences.Columns().EnableAssetAlert] = boolToFloat(*req.PortfolioAlert)
		}

		_, err = dao.NotificationPreferences.Ctx(ctx).
			Where(dao.NotificationPreferences.Columns().UserId, userID).
			Data(updateData).
			Update()
		if err != nil {
			return nil, gerror.Wrap(err, "更新通知偏好失败")
		}
	}

	g.Log().Info(ctx, "更新通知偏好成功", "userId", userID)

	return s.GetPreferences(ctx, userID)
}

// Send 创建并发送通知（内部方法，不暴露路由）
func (s *sNotification) Send(ctx context.Context, userID int, notifType, title, message string) error {
	now := gtime.Now()
	_, err := dao.Notifications.Ctx(ctx).Data(g.Map{
		dao.Notifications.Columns().UserId:    userID,
		dao.Notifications.Columns().Type:      notifType,
		dao.Notifications.Columns().Title:     title,
		dao.Notifications.Columns().Content:   message,
		dao.Notifications.Columns().IsRead:    0,
		dao.Notifications.Columns().CreatedAt: now,
		dao.Notifications.Columns().UpdatedAt: now,
	}).Insert()
	if err != nil {
		return gerror.Wrap(err, "保存通知失败")
	}

	g.Log().Info(ctx, "发送通知成功",
		"userId", userID,
		"type", notifType,
		"title", title,
	)

	// 异步推送 WebPush（不阻塞主流程）
	go func() {
		pushCtx := context.Background()

		// 检查用户是否启用了推送通知
		pref, err := s.GetPreferences(pushCtx, userID)
		if err != nil {
			g.Log().Warning(pushCtx, "查询通知偏好失败，跳过 WebPush", "error", err)
			return
		}
		if !pref.PushEnabled {
			return
		}

		// 调用 WebPush 服务发送推送
		if err := webpushService.Webpush().SendPush(pushCtx, userID, title, message); err != nil {
			g.Log().Warning(pushCtx, "WebPush 推送失败",
				"userId", userID,
				"title", title,
				"error", err,
			)
		}
	}()

	return nil
}

// toNotificationItem 将数据库实体转换为 API 响应格式
func (s *sNotification) toNotificationItem(n *entity.Notifications) notificationApi.NotificationItem {
	return notificationApi.NotificationItem{
		ID:        uint(n.Id),
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Content,
		IsRead:    n.IsRead == 1,
		CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// boolToFloat 将布尔值转换为浮点数（SQLite 兼容）
func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
