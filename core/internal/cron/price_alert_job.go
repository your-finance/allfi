// Package cron 价格预警定时任务
// 定期检查价格预警条件
package cron

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	priceAlertService "your-finance/allfi/internal/app/price_alert/service"
)

// PriceAlertJob 价格预警定时任务
// 定期调用 CheckAlerts 检查所有活跃的价格预警
type PriceAlertJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewPriceAlertJob 创建价格预警定时任务
// interval: 检查间隔，默认每 5 分钟
func NewPriceAlertJob(interval time.Duration) *PriceAlertJob {
	if interval == 0 {
		interval = 5 * time.Minute
	}

	return &PriceAlertJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *PriceAlertJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 价格预警定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *PriceAlertJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 价格预警定时任务已停止")
}

// run 运行定时任务循环
func (j *PriceAlertJob) run() {
	defer j.wg.Done()

	// 启动后等一个周期再开始（避免启动时并发压力）
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

// execute 执行价格预警检查
func (j *PriceAlertJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := priceAlertService.PriceAlert().CheckAlerts(ctx); err != nil {
		g.Log().Errorf(ctx, "[Cron] 价格预警检查失败: %v", err)
	}
}
