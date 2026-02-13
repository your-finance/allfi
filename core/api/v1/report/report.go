// Package report 报告 API 定义
// 提供报告列表、详情、月报、年报、手动生成接口
package report

import "github.com/gogf/gf/v2/frame/g"

// ListReq 获取报告列表请求
type ListReq struct {
	g.Meta `path:"/reports" method:"get" summary:"获取报告列表" tags:"报告"`
	Type   string `json:"type" in:"query" dc:"报告类型筛选（daily/weekly/monthly/annual）"`
	Limit  int    `json:"limit" in:"query" d:"20" dc:"返回数量限制"`
}

// ReportSummary 报告摘要
type ReportSummary struct {
	ID         uint    `json:"id" dc:"报告 ID"`
	Type       string  `json:"type" dc:"报告类型"`
	Title      string  `json:"title" dc:"报告标题"`
	TotalValue float64 `json:"total_value" dc:"总资产价值"`
	PnL        float64 `json:"pnl" dc:"盈亏"`
	PnLPercent float64 `json:"pnl_percent" dc:"盈亏比例"`
	CreatedAt  string  `json:"created_at" dc:"生成时间"`
}

// ListRes 获取报告列表响应
type ListRes struct {
	Reports []ReportSummary `json:"reports" dc:"报告列表"`
}

// GetReq 获取单个报告详情请求
type GetReq struct {
	g.Meta `path:"/reports/{id}" method:"get" summary:"获取报告详情" tags:"报告"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"报告 ID"`
}

// GetRes 获取单个报告详情响应
type GetRes struct {
	Report interface{} `json:"report" dc:"报告详情（包含完整数据）"`
}

// GetMonthlyReq 获取月度报告请求
type GetMonthlyReq struct {
	g.Meta `path:"/reports/monthly/{month}" method:"get" summary:"获取月度报告" tags:"报告"`
	Month  string `json:"month" in:"path" v:"required" dc:"月份（格式: 2024-01）"`
}

// GetMonthlyRes 获取月度报告响应
type GetMonthlyRes struct {
	Report interface{} `json:"report" dc:"月度报告详情"`
}

// GetAnnualReq 获取年度报告请求
type GetAnnualReq struct {
	g.Meta `path:"/reports/annual/{year}" method:"get" summary:"获取年度报告" tags:"报告"`
	Year   string `json:"year" in:"path" v:"required" dc:"年份（格式: 2024）"`
}

// GetAnnualRes 获取年度报告响应
type GetAnnualRes struct {
	Report interface{} `json:"report" dc:"年度报告详情"`
}

// GenerateReq 手动生成报告请求
type GenerateReq struct {
	g.Meta `path:"/reports/generate" method:"post" summary:"手动生成报告" tags:"报告"`
	Type   string `json:"type" v:"required|in:daily,weekly,monthly,annual" dc:"报告类型"`
}

// GenerateRes 手动生成报告响应
type GenerateRes struct {
	Report *ReportSummary `json:"report" dc:"新生成的报告摘要"`
}

// CompareReq 报告对比请求
type CompareReq struct {
	g.Meta    `path:"/reports/compare" method:"get" summary:"对比两份报告" tags:"报告"`
	ReportID1 int `json:"report_id_1" v:"required" in:"query" dc:"报告 1 ID"`
	ReportID2 int `json:"report_id_2" v:"required" in:"query" dc:"报告 2 ID"`
}

// CompareRes 报告对比响应
type CompareRes struct {
	Report1    *CompareReportSummary `json:"report_1" dc:"报告 1 摘要"`
	Report2    *CompareReportSummary `json:"report_2" dc:"报告 2 摘要"`
	ValueDiff  float64               `json:"value_diff" dc:"总价值差异"`
	ChangeDiff float64               `json:"change_diff" dc:"变化差异"`
}

// CompareReportSummary 报告摘要（用于对比）
type CompareReportSummary struct {
	ID          uint    `json:"id" dc:"报告 ID"`
	Type        string  `json:"type" dc:"报告类型"`
	Period      string  `json:"period" dc:"报告周期"`
	TotalValue  float64 `json:"total_value" dc:"总资产价值"`
	Change      float64 `json:"change" dc:"变化金额"`
	ChangePct   float64 `json:"change_percent" dc:"变化百分比"`
	GeneratedAt string  `json:"generated_at" dc:"生成时间"`
}
