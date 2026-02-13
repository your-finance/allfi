// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// WalletAddresses is the golang structure of table wallet_addresses for DAO operations like Where/Data.
type WalletAddresses struct {
	g.Meta     `orm:"table:wallet_addresses, do:true"`
	Id         any //
	CreatedAt  any //
	UpdatedAt  any //
	DeletedAt  any //
	UserId     any //
	Blockchain any //
	Address    any //
	Label      any //
	IsActive   any //
}
