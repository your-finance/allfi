package controller

import (
	"context"
	"time"

	v1 "your-finance/allfi/api/v1"
	"your-finance/allfi/internal/app/cross_chain/logic"

	"github.com/gogf/gf/v2/util/gconv"
)

// CrossChain 跨链交易控制器
var CrossChain = cCrossChain{}

type cCrossChain struct{}

// GetTransactions 获取跨链交易列表
func (c *cCrossChain) GetTransactions(ctx context.Context, req *v1.CrossChainTransactionsReq) (res *v1.CrossChainTransactionsRes, err error) {
	// 获取用户ID（从上下文中获取，实际应用中需要从认证中间件获取）
	userId := gconv.Int64(ctx.Value("user_id"))
	if userId == 0 {
		userId = 1 // 默认用户ID，用于测试
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 查询交易列表
	transactions, total, err := logic.CrossChain().GetTransactions(ctx, userId, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	res = &v1.CrossChainTransactionsRes{
		List:  transactions,
		Total: total,
		Page:  req.Page,
	}
	return res, nil
}

// GetAssetFlow 获取资产流向数据
func (c *cCrossChain) GetAssetFlow(ctx context.Context, req *v1.CrossChainFlowReq) (res *v1.CrossChainFlowRes, err error) {
	// 获取用户ID
	userId := gconv.Int64(ctx.Value("user_id"))
	if userId == 0 {
		userId = 1
	}

	// 解析时间范围
	var startTime, endTime *time.Time
	if req.StartTime != "" {
		t, err := time.Parse("2006-01-02", req.StartTime)
		if err == nil {
			startTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.Parse("2006-01-02", req.EndTime)
		if err == nil {
			endTime = &t
		}
	}

	// 查询资产流向
	flowData, err := logic.CrossChain().GetAssetFlow(ctx, userId, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 转换节点数据
	nodes := make([]v1.CrossChainFlowNode, len(flowData.Nodes))
	for i, node := range flowData.Nodes {
		nodes[i] = v1.CrossChainFlowNode{
			Name:  node.Name,
			Value: node.Value,
		}
	}

	// 转换连接数据
	links := make([]v1.CrossChainFlowLink, len(flowData.Links))
	for i, link := range flowData.Links {
		links[i] = v1.CrossChainFlowLink{
			Source: link.Source,
			Target: link.Target,
			Value:  link.Value,
		}
	}

	res = &v1.CrossChainFlowRes{
		Nodes: nodes,
		Links: links,
	}
	return res, nil
}

// GetFeeStats 获取跨链手续费统计
func (c *cCrossChain) GetFeeStats(ctx context.Context, req *v1.CrossChainFeeStatsReq) (res *v1.CrossChainFeeStatsRes, err error) {
	// 获取用户ID
	userId := gconv.Int64(ctx.Value("user_id"))
	if userId == 0 {
		userId = 1
	}

	// 解析时间范围
	var startTime, endTime *time.Time
	if req.StartTime != "" {
		t, err := time.Parse("2006-01-02", req.StartTime)
		if err == nil {
			startTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.Parse("2006-01-02", req.EndTime)
		if err == nil {
			endTime = &t
		}
	}

	// 查询手续费统计
	stats, err := logic.CrossChain().GetFeeStats(ctx, userId, startTime, endTime)
	if err != nil {
		return nil, err
	}

	res = &v1.CrossChainFeeStatsRes{
		TotalFee:       stats.TotalFee,
		AvgFee:         stats.AvgFee,
		FeeByBridge:    stats.FeeByBridge,
		FeeByChain:     stats.FeeByChain,
		TransactionCnt: stats.TransactionCnt,
	}
	return res, nil
}

// GetBridges 获取支持的跨链桥列表
func (c *cCrossChain) GetBridges(ctx context.Context, req *v1.CrossChainBridgesReq) (res *v1.CrossChainBridgesRes, err error) {
	bridges, err := logic.CrossChain().GetBridges(ctx)
	if err != nil {
		return nil, err
	}

	res = &v1.CrossChainBridgesRes{
		Bridges: bridges,
	}
	return res, nil
}

// InitMockData 初始化模拟数据（仅用于开发测试）
func (c *cCrossChain) InitMockData(ctx context.Context, req *v1.CrossChainInitMockReq) (res *v1.CrossChainInitMockRes, err error) {
	userId := gconv.Int64(ctx.Value("user_id"))
	if userId == 0 {
		userId = 1
	}

	err = logic.CrossChain().InitMockData(ctx, userId)
	if err != nil {
		return nil, err
	}

	res = &v1.CrossChainInitMockRes{
		Message: "Mock data initialized successfully",
	}
	return res, nil
}
