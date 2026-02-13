// Package service 报告模块 - 服务接口定义
// 提供资产报告的列表、详情、生成功能
package service

import (
	"context"

	reportApi "your-finance/allfi/api/v1/report"
)

// IReport 报告服务接口
type IReport interface {
	// GetReports 获取报告列表
	GetReports(ctx context.Context, reportType string, limit int) (*reportApi.ListRes, error)

	// GetReport 获取单个报告详情
	GetReport(ctx context.Context, id uint) (*reportApi.GetRes, error)

	// GenerateReport 生成报告（daily/weekly/monthly/annual）
	GenerateReport(ctx context.Context, reportType string) (*reportApi.GenerateRes, error)

	// GetMonthlyReport 获取/生成月度报告
	GetMonthlyReport(ctx context.Context, month string) (*reportApi.GetMonthlyRes, error)

	// GetAnnualReport 获取/生成年度报告
	GetAnnualReport(ctx context.Context, year string) (*reportApi.GetAnnualRes, error)

	// Compare 对比两份报告
	Compare(ctx context.Context, reportID1, reportID2 int) (*reportApi.CompareRes, error)
}

var localReport IReport

// Report 获取报告服务实例
func Report() IReport {
	if localReport == nil {
		panic("IReport 服务未注册，请检查 logic/report 包的 init 函数")
	}
	return localReport
}

// RegisterReport 注册报告服务实现
// 由 logic 层在 init 函数中调用
func RegisterReport(i IReport) {
	localReport = i
}
