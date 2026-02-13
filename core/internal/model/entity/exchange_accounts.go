// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// ExchangeAccounts is the golang structure for table exchange_accounts.
type ExchangeAccounts struct {
	Id                     int       `json:"id"                       orm:"id"                       description:""` //
	CreatedAt              time.Time `json:"created_at"               orm:"created_at"               description:""` //
	UpdatedAt              time.Time `json:"updated_at"               orm:"updated_at"               description:""` //
	DeletedAt              time.Time `json:"deleted_at"               orm:"deleted_at"               description:""` //
	UserId                 int       `json:"user_id"                  orm:"user_id"                  description:""` //
	ExchangeName           string    `json:"exchange_name"            orm:"exchange_name"            description:""` //
	ApiKeyEncrypted        string    `json:"api_key_encrypted"        orm:"api_key_encrypted"        description:""` //
	ApiSecretEncrypted     string    `json:"api_secret_encrypted"     orm:"api_secret_encrypted"     description:""` //
	ApiPassphraseEncrypted string    `json:"api_passphrase_encrypted" orm:"api_passphrase_encrypted" description:""` //
	Label                  string    `json:"label"                    orm:"label"                    description:""` //
	Note                   string    `json:"note"                     orm:"note"                     description:""` //
	IsActive               float64   `json:"is_active"                orm:"is_active"                description:""` //
}
