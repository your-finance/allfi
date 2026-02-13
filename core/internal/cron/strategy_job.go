// Package cron 策略检查定时任务
// 定期检查所有活跃策略的偏离度，触发提醒通知
package cron

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	notificationService "your-finance/allfi/internal/app/notification/service"
	strategyDao "your-finance/allfi/internal/app/strategy/dao"
	strategyService "your-finance/allfi/internal/app/strategy/service"
	"your-finance/allfi/internal/model/entity"
)

// 策略偏离度超过此阈值时发送通知（10%）
const deviationThreshold = 0.10

// StrategyJob 策略检查定时任务
// 遍历所有活跃策略，分析偏离度，超过阈值时发送通知
type StrategyJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewStrategyJob 创建策略检查定时任务
// interval: 检查间隔，默认每 30 分钟
func NewStrategyJob(interval time.Duration) *StrategyJob {
	if interval == 0 {
		interval = 30 * time.Minute
	}

	return &StrategyJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *StrategyJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 策略检查定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *StrategyJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 策略检查定时任务已停止")
}

// run 运行定时任务循环
func (j *StrategyJob) run() {
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

// execute 执行策略检查
// 查询所有活跃策略 → 逐个分析偏离度 → 超过阈值发送通知
func (j *StrategyJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 查询所有活跃策略
	var strategies []entity.Strategies
	err := strategyDao.Strategies.Ctx(ctx).
		Where("is_active", 1).
		Scan(&strategies)
	if err != nil {
		g.Log().Errorf(ctx, "[Cron] 查询活跃策略失败: %v", err)
		return
	}

	if len(strategies) == 0 {
		return
	}

	// 逐个策略分析偏离度
	for _, s := range strategies {
		analysis, err := strategyService.Strategy().GetAnalysis(ctx, uint(s.Id))
		if err != nil {
			g.Log().Errorf(ctx, "[Cron] 策略分析失败 (ID=%d, %s): %v", s.Id, s.Name, err)
			continue
		}

		if analysis == nil || analysis.Analysis == nil {
			continue
		}

		// 检查是否有任何币种偏离度超过阈值
		for symbol, deviation := range analysis.Analysis.Deviation {
			if math.Abs(deviation) > deviationThreshold {
				// 构造通知消息
				msg := fmt.Sprintf(
					"策略「%s」中 %s 偏离目标配比 %.1f%%，建议调仓",
					s.Name, symbol, deviation*100,
				)
				err := notificationService.Notification().Send(
					ctx, s.UserId, "strategy_alert", "策略偏离提醒", msg,
				)
				if err != nil {
					g.Log().Errorf(ctx, "[Cron] 发送策略偏离通知失败 (策略 %d): %v", s.Id, err)
				}
				// 每个策略只发送一次通知（取最大偏离的币种）
				break
			}
		}
	}
}
