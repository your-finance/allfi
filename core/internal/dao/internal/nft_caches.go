// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NftCachesDao is the data access object for the table nft_caches.
type NftCachesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  NftCachesColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// NftCachesColumns defines and stores column names for the table nft_caches.
type NftCachesColumns struct {
	Id              string //
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
	UserId          string //
	WalletAddress   string //
	ContractAddress string //
	TokenId         string //
	Name            string //
	Description     string //
	ImageUrl        string //
	Collection      string //
	CollectionSlug  string //
	Chain           string //
	FloorPrice      string //
	FloorCurrency   string //
	FloorPriceUsd   string //
	CachedAt        string //
}

// nftCachesColumns holds the columns for the table nft_caches.
var nftCachesColumns = NftCachesColumns{
	Id:              "id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
	UserId:          "user_id",
	WalletAddress:   "wallet_address",
	ContractAddress: "contract_address",
	TokenId:         "token_id",
	Name:            "name",
	Description:     "description",
	ImageUrl:        "image_url",
	Collection:      "collection",
	CollectionSlug:  "collection_slug",
	Chain:           "chain",
	FloorPrice:      "floor_price",
	FloorCurrency:   "floor_currency",
	FloorPriceUsd:   "floor_price_usd",
	CachedAt:        "cached_at",
}

// NewNftCachesDao creates and returns a new DAO object for table data access.
func NewNftCachesDao(handlers ...gdb.ModelHandler) *NftCachesDao {
	return &NftCachesDao{
		group:    "default",
		table:    "nft_caches",
		columns:  nftCachesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NftCachesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NftCachesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NftCachesDao) Columns() NftCachesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NftCachesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NftCachesDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *NftCachesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
