// Package cron 报告定时任务
// 每日/每周/月报/年报自动生成
package cron

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	reportService "your-finance/allfi/internal/app/report/service"
)

// ReportJob 报告定时任务
// 定期检查是否需要生成报告（日报、周报、月报、年报）
type ReportJob struct {
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewReportJob 创建报告定时任务
// interval: 检查间隔，默认每小时检查一次
func NewReportJob(interval time.Duration) *ReportJob {
	if interval == 0 {
		interval = time.Hour
	}

	return &ReportJob{
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
func (j *ReportJob) Start() {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	j.wg.Add(1)
	go j.run()

	g.Log().Infof(context.Background(), "[Cron] 报告定时任务已启动，间隔: %v", j.interval)
}

// Stop 停止定时任务
func (j *ReportJob) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	j.running = false
	j.mu.Unlock()

	close(j.stopChan)
	j.wg.Wait()

	g.Log().Info(context.Background(), "[Cron] 报告定时任务已停止")
}

// run 运行定时任务循环
func (j *ReportJob) run() {
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

// execute 执行报告生成检查
// 根据当前时间判断是否需要生成日报、周报、月报或年报
func (j *ReportJob) execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	now := time.Now()

	// 每日报告：在 21:00-22:00 生成
	if now.Hour() >= 21 && now.Hour() < 22 {
		if _, err := reportService.Report().GenerateReport(ctx, "daily"); err != nil {
			g.Log().Errorf(ctx, "[Cron] 生成每日报告失败: %v", err)
		} else {
			g.Log().Info(ctx, "[Cron] 每日报告生成成功")
		}
	}

	// 每周报告：每周一 10:00-11:00 生成
	if now.Weekday() == time.Monday && now.Hour() >= 10 && now.Hour() < 11 {
		if _, err := reportService.Report().GenerateReport(ctx, "weekly"); err != nil {
			g.Log().Errorf(ctx, "[Cron] 生成每周报告失败: %v", err)
		} else {
			g.Log().Info(ctx, "[Cron] 每周报告生成成功")
		}
	}

	// 月报：每月 1 号 09:00-10:00 生成上月报告
	if now.Day() == 1 && now.Hour() >= 9 && now.Hour() < 10 {
		lastMonth := now.AddDate(0, -1, 0).Format("2006-01")
		if _, err := reportService.Report().GetMonthlyReport(ctx, lastMonth); err != nil {
			g.Log().Errorf(ctx, "[Cron] 生成月报失败 (%s): %v", lastMonth, err)
		} else {
			g.Log().Infof(ctx, "[Cron] 月报生成成功 (%s)", lastMonth)
		}
	}

	// 年报：每年 1 月 1 日 08:00-09:00 生成上年报告
	if now.Month() == time.January && now.Day() == 1 && now.Hour() >= 8 && now.Hour() < 9 {
		lastYear := now.AddDate(-1, 0, 0).Format("2006")
		if _, err := reportService.Report().GetAnnualReport(ctx, lastYear); err != nil {
			g.Log().Errorf(ctx, "[Cron] 生成年报失败 (%s): %v", lastYear, err)
		} else {
			g.Log().Infof(ctx, "[Cron] 年报生成成功 (%s)", lastYear)
		}
	}
}
