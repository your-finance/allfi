// Package cron 风险提醒定时任务
// 每小时检查一次资产组合风险（集中度/波动/平台/缓冲）
package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	healthScoreService "your-finance/allfi/internal/app/health_score/service"
	notificationService "your-finance/allfi/internal/app/notification/service"
)

// 风险评分低于此阈值时发送通知（60 分）
const riskScoreThreshold = 60.0

// RiskAlertJob 风险提醒定时任务
// 定期评估资产组合的健康评分，当评分低于阈值时发送风险提醒通知
type RiskAlertJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewRiskAlertJob 创建风险提醒定时任务
// interval: 检查间隔，默认每小时
func NewRiskAlertJob(interval time.Duration) *RiskAlertJob {
	if interval == 0 {
		interval = 1 * time.Hour
	}

	return &RiskAlertJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *RiskAlertJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 风险提醒定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *RiskAlertJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 风险提醒定时任务已停止")
}

// run 运行定时任务循环
func (j *RiskAlertJob) run() {
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

// execute 执行风险检查
// 获取健康评分 → 检查各维度得分 → 低于阈值时发送通知
func (j *RiskAlertJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 获取资产健康评分（默认 USD 计价）
	result, err := healthScoreService.HealthScore().GetHealthScore(ctx, "USD")
	if err != nil {
		g.Log().Errorf(ctx, "[Cron] 风险检查失败: %v", err)
		return
	}

	if result == nil {
		return
	}

	// 检查综合评分是否低于阈值
	if result.OverallScore < riskScoreThreshold {
		// 构建风险提醒消息
		msg := fmt.Sprintf(
			"资产组合健康评分 %.0f 分（%s），最弱维度: %s。建议: %s",
			result.OverallScore,
			result.Level,
			result.Weakest,
			getFirstAdvice(result.Advice),
		)

		// TODO: 当前单用户模式，发送给用户 1
		err := notificationService.Notification().Send(
			ctx, 1, "risk_alert", "风险提醒", msg,
		)
		if err != nil {
			g.Log().Errorf(ctx, "[Cron] 发送风险提醒通知失败: %v", err)
		} else {
			g.Log().Infof(ctx, "[Cron] 风险提醒已发送，健康评分: %.0f", result.OverallScore)
		}
	}

	// 检查各维度得分，对得分低于阈值的维度单独提醒
	for _, dim := range result.Details {
		if dim.Score < riskScoreThreshold && dim.Suggestion != "" {
			msg := fmt.Sprintf("风险维度「%s」得分 %.0f，%s", dim.Description, dim.Score, dim.Suggestion)
			err := notificationService.Notification().Send(
				ctx, 1, "risk_alert", "风险维度提醒", msg,
			)
			if err != nil {
				g.Log().Errorf(ctx, "[Cron] 发送风险维度通知失败 (%s): %v", dim.Category, err)
			}
		}
	}
}

// getFirstAdvice 获取第一条建议，若无建议则返回默认提示
func getFirstAdvice(advice []string) string {
	if len(advice) > 0 {
		return advice[0]
	}
	return "请关注资产组合分散化配置"
}
