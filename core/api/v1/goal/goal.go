// Package goal 目标追踪 API 定义
// 提供投资目标的增删改查接口
package goal

import "github.com/gogf/gf/v2/frame/g"

// ListReq 获取目标列表请求
type ListReq struct {
	g.Meta `path:"/goals" method:"get" summary:"获取目标列表" tags:"目标追踪"`
}

// GoalItem 目标条目
type GoalItem struct {
	ID           uint    `json:"id" dc:"目标 ID"`
	Title        string  `json:"title" dc:"目标标题"`
	Type         string  `json:"type" dc:"目标类型（asset_value/holding_amount/return_rate）"`
	TargetValue  float64 `json:"target_value" dc:"目标值"`
	CurrentValue float64 `json:"current_value" dc:"当前值"`
	Currency     string  `json:"currency" dc:"货币"`
	Progress     float64 `json:"progress" dc:"进度（百分比）"`
	Deadline     string  `json:"deadline" dc:"截止日期"`
	IsCompleted  bool    `json:"is_completed" dc:"是否已完成"`
	CreatedAt    string  `json:"created_at" dc:"创建时间"`
	UpdatedAt    string  `json:"updated_at" dc:"更新时间"`
}

// ListRes 获取目标列表响应
type ListRes struct {
	Goals []GoalItem `json:"goals" dc:"目标列表"`
}

// CreateReq 创建目标请求
type CreateReq struct {
	g.Meta      `path:"/goals" method:"post" summary:"创建投资目标" tags:"目标追踪"`
	Title       string  `json:"title" v:"required|max-length:100" dc:"目标标题"`
	Type        string  `json:"type" v:"required|in:asset_value,holding_amount,return_rate" dc:"目标类型"`
	TargetValue float64 `json:"target_value" v:"required|min:0" dc:"目标值"`
	Currency    string  `json:"currency" d:"USD" dc:"货币"`
	Deadline    string  `json:"deadline" dc:"截止日期（ISO 8601 格式，可选）"`
}

// CreateRes 创建目标响应
type CreateRes struct {
	Goal *GoalItem `json:"goal" dc:"新建的目标信息"`
}

// UpdateReq 更新目标请求
type UpdateReq struct {
	g.Meta      `path:"/goals/{id}" method:"put" summary:"更新投资目标" tags:"目标追踪"`
	Id          uint    `json:"id" in:"path" v:"required|min:1" dc:"目标 ID"`
	Title       string  `json:"title" dc:"目标标题"`
	Type        string  `json:"type" dc:"目标类型"`
	TargetValue float64 `json:"target_value" dc:"目标值"`
	Currency    string  `json:"currency" dc:"货币"`
	Deadline    string  `json:"deadline" dc:"截止日期"`
}

// UpdateRes 更新目标响应
type UpdateRes struct {
	Goal *GoalItem `json:"goal" dc:"更新后的目标信息"`
}

// DeleteReq 删除目标请求
type DeleteReq struct {
	g.Meta `path:"/goals/{id}" method:"delete" summary:"删除投资目标" tags:"目标追踪"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"目标 ID"`
}

// DeleteRes 删除目标响应
type DeleteRes struct{}
