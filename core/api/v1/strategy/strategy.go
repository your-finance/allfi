// Package strategy 策略引擎 API 定义
// 提供策略的增删改查和分析接口
package strategy

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ListReq 获取策略列表请求
type ListReq struct {
	g.Meta `path:"/strategies" method:"get" summary:"获取策略列表" tags:"策略引擎"`
}

// StrategyItem 策略条目
type StrategyItem struct {
	ID        uint        `json:"id" dc:"策略 ID"`
	Name      string      `json:"name" dc:"策略名称"`
	Type      string      `json:"type" dc:"策略类型（rebalance/dca/stop_limit）"`
	Config    interface{} `json:"config" dc:"策略配置（JSON 对象）"`
	IsActive  bool        `json:"is_active" dc:"是否启用"`
	CreatedAt string      `json:"created_at" dc:"创建时间"`
	UpdatedAt string      `json:"updated_at" dc:"更新时间"`
}

// ListRes 获取策略列表响应
type ListRes struct {
	Strategies []StrategyItem `json:"strategies" dc:"策略列表"`
}

// CreateReq 创建策略请求
type CreateReq struct {
	g.Meta `path:"/strategies" method:"post" summary:"创建策略" tags:"策略引擎"`
	Name   string      `json:"name" v:"required|max-length:100" dc:"策略名称"`
	Type   string      `json:"type" v:"required|in:rebalance,dca,stop_limit" dc:"策略类型"`
	Config interface{} `json:"config" v:"required" dc:"策略配置（JSON 对象）"`
}

// CreateRes 创建策略响应
type CreateRes struct {
	Strategy *StrategyItem `json:"strategy" dc:"新建的策略信息"`
}

// UpdateReq 更新策略请求
type UpdateReq struct {
	g.Meta   `path:"/strategies/{id}" method:"put" summary:"更新策略" tags:"策略引擎"`
	Id       uint        `json:"id" in:"path" v:"required|min:1" dc:"策略 ID"`
	Name     string      `json:"name" dc:"策略名称"`
	Type     string      `json:"type" dc:"策略类型"`
	Config   interface{} `json:"config" dc:"策略配置"`
	IsActive *bool       `json:"is_active" dc:"是否启用"`
}

// UpdateRes 更新策略响应
type UpdateRes struct {
	Strategy *StrategyItem `json:"strategy" dc:"更新后的策略信息"`
}

// DeleteReq 删除策略请求
type DeleteReq struct {
	g.Meta `path:"/strategies/{id}" method:"delete" summary:"删除策略" tags:"策略引擎"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"策略 ID"`
}

// DeleteRes 删除策略响应
type DeleteRes struct{}

// GetAnalysisReq 获取策略分析请求
type GetAnalysisReq struct {
	g.Meta `path:"/strategies/{id}/analysis" method:"get" summary:"获取策略分析" tags:"策略引擎"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"策略 ID"`
}

// AnalysisItem 策略分析结果
type AnalysisItem struct {
	StrategyID     uint                   `json:"strategy_id" dc:"策略 ID"`
	CurrentAlloc   map[string]float64     `json:"current_alloc" dc:"当前配比"`
	TargetAlloc    map[string]float64     `json:"target_alloc" dc:"目标配比"`
	Deviation      map[string]float64     `json:"deviation" dc:"偏离度"`
	Recommendations []RecommendationItem  `json:"recommendations" dc:"调仓建议"`
}

// RecommendationItem 调仓建议条目
type RecommendationItem struct {
	Symbol    string  `json:"symbol" dc:"币种"`
	Action    string  `json:"action" dc:"操作（buy/sell）"`
	Amount    float64 `json:"amount" dc:"数量"`
	ValueUSD  float64 `json:"value_usd" dc:"价值（USD）"`
	Reason    string  `json:"reason" dc:"建议原因"`
}

// GetAnalysisRes 获取策略分析响应
type GetAnalysisRes struct {
	Analysis *AnalysisItem `json:"analysis" dc:"策略分析结果"`
}

// GetRebalanceReq 获取再平衡建议请求
type GetRebalanceReq struct {
	g.Meta `path:"/strategies/{id}/rebalance" method:"get" summary:"获取再平衡建议" tags:"策略引擎"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"策略 ID"`
}

// GetRebalanceRes 获取再平衡建议响应（复用 AnalysisResult 结构）
type GetRebalanceRes = GetAnalysisRes
