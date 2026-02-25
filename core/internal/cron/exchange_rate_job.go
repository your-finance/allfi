package cron

import (
	"context"
	"sync"
	"time"

	exchangeRateService "your-finance/allfi/internal/app/exchange_rate/service"

	"github.com/gogf/gf/v2/frame/g"
)

// ExchangeRateJob 汇率定时更新任务
type ExchangeRateJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewExchangeRateJob 创建汇率定时更新任务
func NewExchangeRateJob(interval time.Duration) *ExchangeRateJob {
	if interval == 0 {
		interval = 30 * time.Minute // 默认30分钟更新一次汇率
	}

	return &ExchangeRateJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (j *ExchangeRateJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 汇率定时更新任务已启动，间隔: %v", j.interval)
}

func (j *ExchangeRateJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 汇率定时更新任务已停止")
}

func (j *ExchangeRateJob) run() {
	defer j.wg.Done()

	// 启动时不立即执行（因为 snapshot_job 启动时已经会执行一次），避免并发冲突
	// j.execute()

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

func (j *ExchangeRateJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	g.Log().Info(ctx, "[Cron] 开始定时刷新汇率...")

	if _, err := exchangeRateService.ExchangeRate().RefreshRates(ctx); err != nil {
		g.Log().Errorf(ctx, "[Cron] 刷新汇率失败: %v", err)
	} else {
		g.Log().Info(ctx, "[Cron] 汇率定时刷新完成")
	}
}
