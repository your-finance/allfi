// Package cron 风险指标计算定时任务
// 每日自动计算风险指标（VaR、夏普比率、最大回撤等）
package cron

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	riskService "your-finance/allfi/internal/app/risk/service"
)

// RiskMetricsJob 风险指标计算定时任务
// 定期计算资产组合的风险指标并保存到数据库
type RiskMetricsJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewRiskMetricsJob 创建风险指标计算定时任务
// interval: 计算间隔，默认每日 00:30
func NewRiskMetricsJob(interval time.Duration) *RiskMetricsJob {
	if interval == 0 {
		interval = 24 * time.Hour // 默认每日执行
	}

	return &RiskMetricsJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *RiskMetricsJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 风险指标计算定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *RiskMetricsJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 风险指标计算定时任务已停止")
}

// run 运行定时任务循环
func (j *RiskMetricsJob) run() {
	defer j.wg.Done()

	// 计算首次执行时间（今日 00:30 或明日 00:30）
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 0, 30, 0, 0, now.Location())
	if now.After(nextRun) {
		// 今日 00:30 已过，设置为明日 00:30
		nextRun = nextRun.Add(24 * time.Hour)
	}

	// 等待到首次执行时间
	waitDuration := time.Until(nextRun)
	g.Log().Infof(context.Background(), "[Cron] 风险指标计算任务将在 %v 后首次执行（%s）",
		waitDuration.Round(time.Second), nextRun.Format("2006-01-02 15:04:05"))

	timer := time.NewTimer(waitDuration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			j.execute()
			// 重置定时器为下一个执行周期
			timer.Reset(j.interval)
		case <-j.stopChan:
			return
		}
	}
}

// execute 执行风险指标计算
// 计算过去 30 天的风险指标并保存到数据库
func (j *RiskMetricsJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	g.Log().Info(ctx, "[Cron] 开始计算风险指标...")

	// 计算风险指标（默认 30 天周期）
	metrics, err := riskService.Risk().CalculateMetrics(ctx, 30)
	if err != nil {
		g.Log().Errorf(ctx, "[Cron] 计算风险指标失败: %v", err)
		return
	}

	g.Log().Infof(ctx, "[Cron] 风险指标计算完成 - 夏普比率: %.2f, 最大回撤: %.2f%%, 波动率: %.2f%%",
		metrics.SharpeRatio, metrics.MaxDrawdown, metrics.Volatility)
}
