// Package health_score 资产健康评分 API 定义
// 提供投资组合健康度评分接口
package health_score

import "github.com/gogf/gf/v2/frame/g"

// GetReq 获取资产健康评分请求
type GetReq struct {
	g.Meta   `path:"/portfolio/health" method:"get" summary:"获取资产健康评分" tags:"资产健康"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// ScoreDetail 评分细项
type ScoreDetail struct {
	Category    string  `json:"category" dc:"评分类别"`
	Score       float64 `json:"score" dc:"分数（0-100）"`
	Weight      float64 `json:"weight" dc:"权重"`
	Description string  `json:"description" dc:"说明"`
	Suggestion  string  `json:"suggestion" dc:"改善建议"`
}

// GetRes 获取资产健康评分响应
type GetRes struct {
	OverallScore float64       `json:"overall_score" dc:"综合评分（0-100）"`
	Level        string        `json:"level" dc:"健康等级（excellent/good/fair/poor）"`
	Details      []ScoreDetail `json:"details" dc:"评分细项"`
	Currency     string        `json:"currency" dc:"计价货币"`
	UpdatedAt    string        `json:"updated_at" dc:"评分更新时间"`
}
