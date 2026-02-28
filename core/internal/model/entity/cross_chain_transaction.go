package entity

import (
	"time"
)

// CrossChainTransaction 跨链交易实体
type CrossChainTransaction struct {
	Id        int64     `json:"id"         orm:"id,primary" description:"主键ID"`
	CreatedAt time.Time `json:"created_at" orm:"created_at" description:"创建时间"`
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at" description:"更新时间"`
	UserId    int64     `json:"user_id"    orm:"user_id"    description:"用户ID"`

	// 交易信息
	TxHash         string `json:"tx_hash"         orm:"tx_hash"         description:"交易哈希"`
	BridgeProtocol string `json:"bridge_protocol" orm:"bridge_protocol" description:"跨链桥协议"`

	// 源链信息
	SourceChain  string  `json:"source_chain"  orm:"source_chain"  description:"源链"`
	SourceToken  string  `json:"source_token"  orm:"source_token"  description:"源代币"`
	SourceAmount float64 `json:"source_amount" orm:"source_amount" description:"源金额"`

	// 目标链信息
	DestChain  string  `json:"dest_chain"  orm:"dest_chain"  description:"目标链"`
	DestToken  string  `json:"dest_token"  orm:"dest_token"  description:"目标代币"`
	DestAmount float64 `json:"dest_amount" orm:"dest_amount" description:"目标金额"`

	// 费用
	BridgeFee   float64 `json:"bridge_fee"    orm:"bridge_fee"    description:"跨链桥费用"`
	GasFee      float64 `json:"gas_fee"       orm:"gas_fee"       description:"Gas费用"`
	TotalFeeUsd float64 `json:"total_fee_usd" orm:"total_fee_usd" description:"总费用(USD)"`

	// 状态
	Status      string     `json:"status"       orm:"status"       description:"状态"`
	InitiatedAt time.Time  `json:"initiated_at" orm:"initiated_at" description:"发起时间"`
	CompletedAt *time.Time `json:"completed_at" orm:"completed_at" description:"完成时间"`
}
