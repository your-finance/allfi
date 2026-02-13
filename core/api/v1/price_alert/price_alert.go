// Package price_alert 价格预警 API 定义
// 提供价格预警的增删改查接口
package price_alert

import "github.com/gogf/gf/v2/frame/g"

// CreateReq 创建价格预警请求
type CreateReq struct {
	g.Meta      `path:"/alerts" method:"post" summary:"创建价格预警" tags:"价格预警"`
	Symbol      string  `json:"symbol" v:"required" dc:"币种符号（如 BTC、ETH）"`
	Condition   string  `json:"condition" v:"required|in:above,below" dc:"预警条件（above=高于/below=低于）"`
	TargetPrice float64 `json:"target_price" v:"required|min:0" dc:"目标价格"`
	Note        string  `json:"note" dc:"备注"`
}

// AlertItem 价格预警条目
type AlertItem struct {
	ID           uint    `json:"id" dc:"预警 ID"`
	Symbol       string  `json:"symbol" dc:"币种符号"`
	Condition    string  `json:"condition" dc:"预警条件"`
	TargetPrice  float64 `json:"target_price" dc:"目标价格"`
	CurrentPrice float64 `json:"current_price" dc:"当前价格"`
	Note         string  `json:"note" dc:"备注"`
	IsActive     bool    `json:"is_active" dc:"是否启用"`
	TriggeredAt  string  `json:"triggered_at" dc:"触发时间"`
	CreatedAt    string  `json:"created_at" dc:"创建时间"`
}

// CreateRes 创建价格预警响应
type CreateRes struct {
	Alert *AlertItem `json:"alert" dc:"新建的预警信息"`
}

// ListReq 获取价格预警列表请求
type ListReq struct {
	g.Meta `path:"/alerts" method:"get" summary:"获取价格预警列表" tags:"价格预警"`
}

// ListRes 获取价格预警列表响应
type ListRes struct {
	Alerts []AlertItem `json:"alerts" dc:"预警列表"`
}

// UpdateReq 更新价格预警请求
type UpdateReq struct {
	g.Meta      `path:"/alerts/{id}" method:"put" summary:"更新价格预警" tags:"价格预警"`
	Id          uint    `json:"id" in:"path" v:"required|min:1" dc:"预警 ID"`
	Symbol      string  `json:"symbol" dc:"币种符号"`
	Condition   string  `json:"condition" dc:"预警条件（above/below）"`
	TargetPrice float64 `json:"target_price" dc:"目标价格"`
	Note        string  `json:"note" dc:"备注"`
	IsActive    *bool   `json:"is_active" dc:"是否启用"`
}

// UpdateRes 更新价格预警响应
type UpdateRes struct {
	Alert *AlertItem `json:"alert" dc:"更新后的预警信息"`
}

// DeleteReq 删除价格预警请求
type DeleteReq struct {
	g.Meta `path:"/alerts/{id}" method:"delete" summary:"删除价格预警" tags:"价格预警"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"预警 ID"`
}

// DeleteRes 删除价格预警响应
type DeleteRes struct{}
