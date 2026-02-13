// Package cron 通知定时任务
// 每日摘要生成和推送
package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	notificationDao "your-finance/allfi/internal/app/notification/dao"
	notificationService "your-finance/allfi/internal/app/notification/service"
	"your-finance/allfi/internal/model/entity"
)

// NotificationJob 通知定时任务
// 定期检查并为启用每日摘要的用户生成通知
type NotificationJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewNotificationJob 创建通知定时任务
// interval: 检查间隔，默认每小时检查一次
func NewNotificationJob(interval time.Duration) *NotificationJob {
	if interval == 0 {
		interval = time.Hour
	}

	return &NotificationJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *NotificationJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 通知定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *NotificationJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 通知定时任务已停止")
}

// run 运行定时任务循环
func (j *NotificationJob) run() {
	defer j.wg.Done()

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.execute()
		case <-j.stopChan:
			return
		}
	}
}

// execute 执行每日摘要检查
// 查询启用每日摘要的用户偏好，判断当前时间是否匹配用户设定的摘要时间
func (j *NotificationJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 查询所有启用每日摘要的用户偏好
	var prefs []entity.NotificationPreferences
	err := notificationDao.NotificationPreferences.Ctx(ctx).
		Where("enable_daily_digest", 1).
		Scan(&prefs)
	if err != nil {
		g.Log().Errorf(ctx, "[Cron] 获取通知偏好列表失败: %v", err)
		return
	}

	// 当前小时（HH:00 格式）
	now := time.Now()
	currentHour := now.Format("15:00")

	for _, pref := range prefs {
		// 检查当前时间是否匹配用户设定的摘要时间（精确到小时）
		if len(pref.DigestTime) >= 5 && pref.DigestTime[:2]+":00" == currentHour {
			// 生成每日摘要内容
			msg := fmt.Sprintf("您的每日资产摘要 - %s", now.Format("2006-01-02"))
			err := notificationService.Notification().Send(ctx, pref.UserId, "digest", "每日摘要", msg)
			if err != nil {
				g.Log().Errorf(ctx, "[Cron] 生成每日摘要失败 (用户 %d): %v", pref.UserId, err)
				continue
			}
			g.Log().Infof(ctx, "[Cron] 每日摘要已生成 (用户 %d)", pref.UserId)
		}
	}
}
