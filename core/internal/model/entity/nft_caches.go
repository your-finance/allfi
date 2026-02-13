// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// NftCaches is the golang structure for table nft_caches.
type NftCaches struct {
	Id              int       `json:"id"               orm:"id"               description:""` //
	CreatedAt       time.Time `json:"created_at"       orm:"created_at"       description:""` //
	UpdatedAt       time.Time `json:"updated_at"       orm:"updated_at"       description:""` //
	DeletedAt       time.Time `json:"deleted_at"       orm:"deleted_at"       description:""` //
	UserId          int       `json:"user_id"          orm:"user_id"          description:""` //
	WalletAddress   string    `json:"wallet_address"   orm:"wallet_address"   description:""` //
	ContractAddress string    `json:"contract_address" orm:"contract_address" description:""` //
	TokenId         string    `json:"token_id"         orm:"token_id"         description:""` //
	Name            string    `json:"name"             orm:"name"             description:""` //
	Description     string    `json:"description"      orm:"description"      description:""` //
	ImageUrl        string    `json:"image_url"        orm:"image_url"        description:""` //
	Collection      string    `json:"collection"       orm:"collection"       description:""` //
	CollectionSlug  string    `json:"collection_slug"  orm:"collection_slug"  description:""` //
	Chain           string    `json:"chain"            orm:"chain"            description:""` //
	FloorPrice      float32   `json:"floor_price"      orm:"floor_price"      description:""` //
	FloorCurrency   string    `json:"floor_currency"   orm:"floor_currency"   description:""` //
	FloorPriceUsd   float32   `json:"floor_price_usd"  orm:"floor_price_usd"  description:""` //
	CachedAt        time.Time `json:"cached_at"        orm:"cached_at"        description:""` //
}
