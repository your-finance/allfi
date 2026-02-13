// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ExchangeAccounts is the golang structure of table exchange_accounts for DAO operations like Where/Data.
type ExchangeAccounts struct {
	g.Meta                 `orm:"table:exchange_accounts, do:true"`
	Id                     any //
	CreatedAt              any //
	UpdatedAt              any //
	DeletedAt              any //
	UserId                 any //
	ExchangeName           any //
	ApiKeyEncrypted        any //
	ApiSecretEncrypted     any //
	ApiPassphraseEncrypted any //
	Label                  any //
	Note                   any //
	IsActive               any //
}
