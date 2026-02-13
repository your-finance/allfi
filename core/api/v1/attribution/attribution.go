// Package attribution 资产归因分析 API 定义
// 提供资产收益归因分析接口
package attribution

import "github.com/gogf/gf/v2/frame/g"

// GetReq 获取资产归因分析请求
type GetReq struct {
	g.Meta   `path:"/analytics/attribution" method:"get" summary:"获取资产归因分析" tags:"归因分析"`
	Days     int    `json:"days" in:"query" d:"7" dc:"分析天数"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// AttributionItem 归因分析条目
type AttributionItem struct {
	Symbol       string  `json:"symbol" dc:"币种符号"`
	Source       string  `json:"source" dc:"来源"`
	Contribution float64 `json:"contribution" dc:"贡献金额"`
	Weight       float64 `json:"weight" dc:"权重（百分比）"`
	Return       float64 `json:"return" dc:"收益率"`
}

// GetRes 获取资产归因分析响应
type GetRes struct {
	TotalReturn   float64            `json:"total_return" dc:"组合总收益"`
	TotalPercent  float64            `json:"total_percent" dc:"组合总收益率"`
	Attributions  []AttributionItem  `json:"attributions" dc:"归因明细"`
	Days          int                `json:"days" dc:"分析天数"`
	Currency      string             `json:"currency" dc:"计价货币"`
}
