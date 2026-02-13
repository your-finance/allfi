// Package model 资产归因分析模块 - 数据传输对象
package model

import "time"

// AttributionResult 归因分析结果
type AttributionResult struct {
	Range          string             `json:"range"`           // 分析范围标识（1d/7d/30d）
	StartTime      time.Time          `json:"start_time"`      // 起始时间
	EndTime        time.Time          `json:"end_time"`        // 结束时间
	TotalChange    float64            `json:"total_change"`    // 总价值变化
	PriceEffect    float64            `json:"price_effect"`    // 总价格效应
	QuantityEffect float64            `json:"quantity_effect"` // 总数量效应
	Assets         []AssetAttribution `json:"assets"`          // 各资产归因明细
	Currency       string             `json:"currency"`        // 计价货币
}

// AssetAttribution 单个资产的归因分析
type AssetAttribution struct {
	Symbol            string  `json:"symbol"`             // 资产符号
	Name              string  `json:"name"`               // 资产名称
	StartBalance      float64 `json:"start_balance"`      // 起始余额
	EndBalance        float64 `json:"end_balance"`        // 结束余额
	StartPrice        float64 `json:"start_price"`        // 起始价格
	EndPrice          float64 `json:"end_price"`          // 结束价格
	StartValue        float64 `json:"start_value"`        // 起始价值
	EndValue          float64 `json:"end_value"`          // 结束价值
	TotalChange       float64 `json:"total_change"`       // 总变化
	PriceEffect       float64 `json:"price_effect"`       // 价格效应 = 起始数量 * (结束价格 - 起始价格)
	QuantityEffect    float64 `json:"quantity_effect"`    // 数量效应 = 起始价格 * (结束数量 - 起始数量)
	InteractionEffect float64 `json:"interaction_effect"` // 交叉效应 = (结束数量 - 起始数量) * (结束价格 - 起始价格)
}
