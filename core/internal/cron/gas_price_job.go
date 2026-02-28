package cron

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	gasService "your-finance/allfi/internal/app/gas/service"
)

// GasPriceJob Gas 价格定时记录任务
type GasPriceJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewGasPriceJob 创建 Gas 价格定时记录任务
func NewGasPriceJob(interval time.Duration) *GasPriceJob {
	if interval == 0 {
		interval = 5 * time.Minute // 默认 5 分钟记录一次
	}

	return &GasPriceJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (j *GasPriceJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] Gas 价格定时记录任务已启动，间隔: %v", j.interval)
}

func (j *GasPriceJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] Gas 价格定时记录任务已停止")
}

func (j *GasPriceJob) run() {
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

func (j *GasPriceJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	g.Log().Info(ctx, "[Cron] 开始记录 Gas 价格...")

	if err := gasService.Gas().RecordGasPrice(ctx); err != nil {
		g.Log().Errorf(ctx, "[Cron] 记录 Gas 价格失败: %v", err)
	} else {
		g.Log().Info(ctx, "[Cron] Gas 价格记录完成")
	}
}
