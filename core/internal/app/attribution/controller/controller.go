// Package controller 资产归因分析模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	attrApi "your-finance/allfi/api/v1/attribution"
	"your-finance/allfi/internal/app/attribution/service"
)

// AttributionController 资产归因分析控制器
type AttributionController struct{}

// Get 获取资产归因分析
func (c *AttributionController) Get(ctx context.Context, req *attrApi.GetReq) (res *attrApi.GetRes, err error) {
	// 调用服务层获取归因分析
	result, err := service.Attribution().GetAttribution(ctx, req.Days, req.Currency)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &attrApi.GetRes{
		TotalReturn:  result.TotalChange,
		Days:         req.Days,
		Currency:     req.Currency,
		Attributions: make([]attrApi.AttributionItem, 0, len(result.Assets)),
	}

	// 计算总收益率
	if len(result.Assets) > 0 {
		var totalStartValue float64
		for _, a := range result.Assets {
			totalStartValue += a.StartValue
		}
		if totalStartValue > 0 {
			res.TotalPercent = result.TotalChange / totalStartValue * 100
		}
	}

	// 转换各资产归因
	for _, a := range result.Assets {
		var weight float64
		if result.TotalChange != 0 {
			weight = a.TotalChange / result.TotalChange * 100
		}
		var returnPct float64
		if a.StartValue > 0 {
			returnPct = a.TotalChange / a.StartValue * 100
		}

		res.Attributions = append(res.Attributions, attrApi.AttributionItem{
			Symbol:       a.Symbol,
			Source:       "", // 由前端根据需要补充
			Contribution: a.TotalChange,
			Weight:       weight,
			Return:       returnPct,
		})
	}

	return res, nil
}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&AttributionController{})
}
