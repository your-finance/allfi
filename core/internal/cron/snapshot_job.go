// Package cron 定时任务
// 资产快照定时任务 + CronManager 管理器
package cron

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	assetService "your-finance/allfi/internal/app/asset/service"
	exchangeRateService "your-finance/allfi/internal/app/exchange_rate/service"
)

// SnapshotJob 资产快照定时任务
// 定期刷新汇率并创建资产快照
type SnapshotJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewSnapshotJob 创建快照定时任务
// interval: 快照间隔，默认 1 小时
func NewSnapshotJob(interval time.Duration) *SnapshotJob {
	if interval == 0 {
		interval = time.Hour // 默认 1 小时
	}

	return &SnapshotJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *SnapshotJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 快照定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *SnapshotJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 快照定时任务已停止")
}

// run 运行定时任务循环
func (j *SnapshotJob) run() {
	defer j.wg.Done()

	// 启动时立即执行一次
	j.execute()

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

// execute 执行快照任务
// 先刷新汇率缓存，再刷新全部资产（内部已包含快照创建逻辑）
func (j *SnapshotJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	g.Log().Info(ctx, "[Cron] 开始创建资产快照...")

	// 刷新汇率缓存
	if _, err := exchangeRateService.ExchangeRate().RefreshRates(ctx); err != nil {
		g.Log().Errorf(ctx, "[Cron] 刷新汇率缓存失败: %v", err)
	}

	// 刷新全部资产（内部包含快照创建逻辑）
	if _, err := assetService.Asset().RefreshAll(ctx); err != nil {
		g.Log().Errorf(ctx, "[Cron] 刷新资产/创建快照失败: %v", err)
		return
	}

	g.Log().Info(ctx, "[Cron] 资产快照创建成功")
}

// CronManager 定时任务管理器
// 统一管理所有 cron job 的启动和停止
type CronManager struct {
	snapshotJob     *SnapshotJob
	notificationJob *NotificationJob
	priceAlertJob   *PriceAlertJob
	reportJob       *ReportJob
	strategyJob     *StrategyJob
	riskAlertJob    *RiskAlertJob
}

// NewCronManager 创建定时任务管理器
// 无需传入 service 参数，所有 job 内部通过 service 单例获取接口
func NewCronManager() *CronManager {
	return &CronManager{
		snapshotJob:     NewSnapshotJob(time.Hour),
		notificationJob: NewNotificationJob(time.Hour),
		priceAlertJob:   NewPriceAlertJob(5 * time.Minute),
		reportJob:       NewReportJob(time.Hour),
		strategyJob:     NewStrategyJob(30 * time.Minute),
		riskAlertJob:    NewRiskAlertJob(time.Hour),
	}
}

// Start 启动所有定时任务
func (m *CronManager) Start() {
	if m.snapshotJob != nil {
		m.snapshotJob.Start()
	}
	if m.notificationJob != nil {
		m.notificationJob.Start()
	}
	if m.priceAlertJob != nil {
		m.priceAlertJob.Start()
	}
	if m.reportJob != nil {
		m.reportJob.Start()
	}
	if m.strategyJob != nil {
		m.strategyJob.Start()
	}
	if m.riskAlertJob != nil {
		m.riskAlertJob.Start()
	}

	g.Log().Info(context.Background(), "[Cron] 所有定时任务已启动")
}

// Stop 停止所有定时任务
func (m *CronManager) Stop() {
	if m.snapshotJob != nil {
		m.snapshotJob.Stop()
	}
	if m.notificationJob != nil {
		m.notificationJob.Stop()
	}
	if m.priceAlertJob != nil {
		m.priceAlertJob.Stop()
	}
	if m.reportJob != nil {
		m.reportJob.Stop()
	}
	if m.strategyJob != nil {
		m.strategyJob.Stop()
	}
	if m.riskAlertJob != nil {
		m.riskAlertJob.Stop()
	}

	g.Log().Info(context.Background(), "[Cron] 所有定时任务已停止")
}
