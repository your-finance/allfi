// Package controller 报告模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	reportApi "your-finance/allfi/api/v1/report"
	"your-finance/allfi/internal/app/report/service"
)

// Controller 报告控制器
type Controller struct{}

// List 获取报告列表
func (c *Controller) List(ctx context.Context, req *reportApi.ListReq) (res *reportApi.ListRes, err error) {
	return service.Report().GetReports(ctx, req.Type, req.Limit)
}

// Get 获取报告详情
func (c *Controller) Get(ctx context.Context, req *reportApi.GetReq) (res *reportApi.GetRes, err error) {
	return service.Report().GetReport(ctx, req.Id)
}

// GetMonthly 获取月度报告
func (c *Controller) GetMonthly(ctx context.Context, req *reportApi.GetMonthlyReq) (res *reportApi.GetMonthlyRes, err error) {
	return service.Report().GetMonthlyReport(ctx, req.Month)
}

// GetAnnual 获取年度报告
func (c *Controller) GetAnnual(ctx context.Context, req *reportApi.GetAnnualReq) (res *reportApi.GetAnnualRes, err error) {
	return service.Report().GetAnnualReport(ctx, req.Year)
}

// Generate 手动生成报告
func (c *Controller) Generate(ctx context.Context, req *reportApi.GenerateReq) (res *reportApi.GenerateRes, err error) {
	return service.Report().GenerateReport(ctx, req.Type)
}

// Compare 对比两份报告
func (c *Controller) Compare(ctx context.Context, req *reportApi.CompareReq) (res *reportApi.CompareRes, err error) {
	return service.Report().Compare(ctx, req.ReportID1, req.ReportID2)
}

// Register 注册报告模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
