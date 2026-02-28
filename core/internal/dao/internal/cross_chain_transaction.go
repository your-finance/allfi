// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CrossChainTransactionDao 跨链交易数据访问对象
type CrossChainTransactionDao struct {
	table   string                       // 表名
	group   string                       // 数据库配置组名
	columns CrossChainTransactionColumns // 字段列表
}

// CrossChainTransactionColumns 字段列表
type CrossChainTransactionColumns struct {
	Id             string
	CreatedAt      string
	UpdatedAt      string
	UserId         string
	TxHash         string
	BridgeProtocol string
	SourceChain    string
	SourceToken    string
	SourceAmount   string
	DestChain      string
	DestToken      string
	DestAmount     string
	BridgeFee      string
	GasFee         string
	TotalFeeUsd    string
	Status         string
	InitiatedAt    string
	CompletedAt    string
}

// crossChainTransactionColumns 字段列表实例
var crossChainTransactionColumns = CrossChainTransactionColumns{
	Id:             "id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	UserId:         "user_id",
	TxHash:         "tx_hash",
	BridgeProtocol: "bridge_protocol",
	SourceChain:    "source_chain",
	SourceToken:    "source_token",
	SourceAmount:   "source_amount",
	DestChain:      "dest_chain",
	DestToken:      "dest_token",
	DestAmount:     "dest_amount",
	BridgeFee:      "bridge_fee",
	GasFee:         "gas_fee",
	TotalFeeUsd:    "total_fee_usd",
	Status:         "status",
	InitiatedAt:    "initiated_at",
	CompletedAt:    "completed_at",
}

// NewCrossChainTransactionDao 创建跨链交易数据访问对象
func NewCrossChainTransactionDao() *CrossChainTransactionDao {
	return &CrossChainTransactionDao{
		group:   "default",
		table:   "cross_chain_transactions",
		columns: crossChainTransactionColumns,
	}
}

// DB 获取数据库连接
func (dao *CrossChainTransactionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 获取表名
func (dao *CrossChainTransactionDao) Table() string {
	return dao.table
}

// Columns 获取字段列表
func (dao *CrossChainTransactionDao) Columns() CrossChainTransactionColumns {
	return dao.columns
}

// Group 获取数据库配置组名
func (dao *CrossChainTransactionDao) Group() string {
	return dao.group
}

// Ctx 创建带上下文的数据库模型
func (dao *CrossChainTransactionDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction 在事务中执行
func (dao *CrossChainTransactionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
