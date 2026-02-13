// Package model 资产健康评分模块业务 DTO
// 定义健康评分相关的数据传输对象
package model

// HealthScoreResult 健康评分结果
type HealthScoreResult struct {
	OverallScore float64                `json:"overall_score"` // 综合评分（0-100）
	Level        string                 `json:"level"`         // 健康等级（excellent/good/fair/poor）
	Details      []HealthScoreDimension `json:"details"`       // 各维度评分细项
	Currency     string                 `json:"currency"`      // 计价货币
	Weakest      string                 `json:"weakest"`       // 最弱维度标识
	Advice       []string               `json:"advice"`        // 改善建议
}

// HealthScoreDimension 健康评分维度
type HealthScoreDimension struct {
	Category    string  `json:"category"`    // 评分类别标识
	Score       float64 `json:"score"`       // 分数（0-100）
	Weight      float64 `json:"weight"`      // 权重（0-1）
	MaxScore    float64 `json:"max_score"`   // 满分值
	Description string  `json:"description"` // 维度说明
	Suggestion  string  `json:"suggestion"`  // 改善建议
	Value       float64 `json:"value"`       // 实际计算值（占比/数量等）
}
