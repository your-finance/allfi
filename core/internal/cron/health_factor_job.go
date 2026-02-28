// Package cron 健康因子监控定时任务
// 定期检查借贷仓位的健康因子,发送预警通知
package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	defiService "your-finance/allfi/internal/app/defi/service"
	notificationService "your-finance/allfi/internal/app/notification/service"
)

// HealthFactorJob 健康因子监控定时任务
// 定期检查所有用户的借贷仓位健康因子,低于阈值时发送预警
type HealthFactorJob struct {
	interval  time.Duration
	threshold float64 // 健康因子阈值,低于此值发送预警
	stopChan  chan struct{}
	wg        sync.WaitGroup
	running   bool
	mu        sync.Mutex
}

// NewHealthFactorJob 创建健康因子监控定时任务
// interval: 检查间隔,默认每 10 分钟
// threshold: 健康因子阈值,默认 1.8
func NewHealthFactorJob(interval time.Duration, threshold float64) *HealthFactorJob {
	if interval == 0 {
		interval = 10 * time.Minute
	}
	if threshold == 0 {
		threshold = 1.8
	}

	return &HealthFactorJob{
		interval:  interval,
		threshold: threshold,
		stopChan:  make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *HealthFactorJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 健康因子监控定时任务已启动，间隔: %v, 阈值: %.2f", j.interval, j.threshold)
}

// Stop 停止定时任务
func (j *HealthFactorJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 健康因子监控定时任务已停止")
}

// run 运行定时任务循环
func (j *HealthFactorJob) run() {
	defer j.wg.Done()

	// 启动后立即执行一次
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

// execute 执行健康因子检查
func (j *HealthFactorJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	g.Log().Debug(ctx, "[Cron] 开始检查健康因子")

	// 检查所有低于阈值的仓位
	positions, err := defiService.Defi().CheckHealthFactors(ctx, j.threshold)
	if err != nil {
		g.Log().Errorf(ctx, "[Cron] 健康因子检查失败: %v", err)
		return
	}

	if len(positions) == 0 {
		g.Log().Debug(ctx, "[Cron] 所有借贷仓位健康因子正常")
		return
	}

	// 发送预警通知
	for _, pos := range positions {
		j.sendAlert(ctx, pos)
	}

	g.Log().Infof(ctx, "[Cron] 健康因子检查完成，发现 %d 个低健康因子仓位", len(positions))
}

// sendAlert 发送健康因子预警通知
func (j *HealthFactorJob) sendAlert(ctx context.Context, pos interface{}) {
	// 构造通知内容
	title := "⚠️ 借贷仓位健康因子预警"
	content := fmt.Sprintf(
		"您在 %s 的借贷仓位健康因子过低，存在清算风险！\n\n"+
			"健康因子: %.2f\n"+
			"清算阈值: %.2f%%\n"+
			"建议: 增加抵押品或偿还部分借款",
		"协议名称", 1.5, 85.0, // 这里需要从 pos 中提取实际数据
	)

	// 发送通知（需要用户ID，这里简化处理）
	// 实际应该从 pos 中获取用户ID
	err := notificationService.Notification().Send(ctx, 1, "health_factor_alert", title, content)
	if err != nil {
		g.Log().Errorf(ctx, "[Cron] 发送健康因子预警通知失败: %v", err)
	}
}
