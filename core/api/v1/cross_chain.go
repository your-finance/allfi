package v1

import (
	"your-finance/allfi/internal/integrations/bridge"
	"your-finance/allfi/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// CrossChainTransactionsReq 获取跨链交易列表请求
type CrossChainTransactionsReq struct {
	g.Meta   `path:"/cross-chain/transactions" method:"get" tags:"CrossChain" summary:"获取跨链交易列表"`
	Page     int `json:"page"      dc:"页码" v:"min:1"`
	PageSize int `json:"page_size" dc:"每页数量" v:"min:1|max:100"`
}

// CrossChainTransactionsRes 获取跨链交易列表响应
type CrossChainTransactionsRes struct {
	List  []*entity.CrossChainTransaction `json:"list"  dc:"交易列表"`
	Total int                             `json:"total" dc:"总数"`
	Page  int                             `json:"page"  dc:"当前页"`
}

// CrossChainFlowReq 获取资产流向请求
type CrossChainFlowReq struct {
	g.Meta    `path:"/cross-chain/flow" method:"get" tags:"CrossChain" summary:"获取资产流向数据"`
	StartTime string `json:"start_time" dc:"开始时间(YYYY-MM-DD)"`
	EndTime   string `json:"end_time"   dc:"结束时间(YYYY-MM-DD)"`
}

// CrossChainFlowRes 获取资产流向响应
type CrossChainFlowRes struct {
	Nodes []CrossChainFlowNode `json:"nodes" dc:"节点列表"`
	Links []CrossChainFlowLink `json:"links" dc:"连接列表"`
}

// CrossChainFlowNode 流向节点
type CrossChainFlowNode struct {
	Name  string `json:"name"  dc:"节点名称"`
	Value int    `json:"value" dc:"节点值"`
}

// CrossChainFlowLink 流向连接
type CrossChainFlowLink struct {
	Source string  `json:"source" dc:"源节点"`
	Target string  `json:"target" dc:"目标节点"`
	Value  float64 `json:"value"  dc:"流向值"`
}

// CrossChainFeeStatsReq 获取手续费统计请求
type CrossChainFeeStatsReq struct {
	g.Meta    `path:"/cross-chain/fees" method:"get" tags:"CrossChain" summary:"获取跨链手续费统计"`
	StartTime string `json:"start_time" dc:"开始时间(YYYY-MM-DD)"`
	EndTime   string `json:"end_time"   dc:"结束时间(YYYY-MM-DD)"`
}

// CrossChainFeeStatsRes 获取手续费统计响应
type CrossChainFeeStatsRes struct {
	TotalFee       float64            `json:"total_fee"        dc:"总手续费"`
	AvgFee         float64            `json:"avg_fee"          dc:"平均手续费"`
	FeeByBridge    map[string]float64 `json:"fee_by_bridge"    dc:"按跨链桥统计"`
	FeeByChain     map[string]float64 `json:"fee_by_chain"     dc:"按链统计"`
	TransactionCnt int                `json:"transaction_count" dc:"交易次数"`
}

// CrossChainBridgesReq 获取跨链桥列表请求
type CrossChainBridgesReq struct {
	g.Meta `path:"/cross-chain/bridges" method:"get" tags:"CrossChain" summary:"获取支持的跨链桥列表"`
}

// CrossChainBridgesRes 获取跨链桥列表响应
type CrossChainBridgesRes struct {
	Bridges []*bridge.BridgeInfo `json:"bridges" dc:"跨链桥列表"`
}

// CrossChainInitMockReq 初始化模拟数据请求
type CrossChainInitMockReq struct {
	g.Meta `path:"/cross-chain/init-mock" method:"post" tags:"CrossChain" summary:"初始化模拟数据"`
}

// CrossChainInitMockRes 初始化模拟数据响应
type CrossChainInitMockRes struct {
	Message string `json:"message" dc:"消息"`
}
