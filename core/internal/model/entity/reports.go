// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// Reports is the golang structure for table reports.
type Reports struct {
	Id            int       `json:"id"             orm:"id"             description:""` //
	CreatedAt     time.Time `json:"created_at"     orm:"created_at"     description:""` //
	UpdatedAt     time.Time `json:"updated_at"     orm:"updated_at"     description:""` //
	DeletedAt     time.Time `json:"deleted_at"     orm:"deleted_at"     description:""` //
	UserId        int       `json:"user_id"        orm:"user_id"        description:""` //
	Type          string    `json:"type"           orm:"type"           description:""` //
	Period        string    `json:"period"         orm:"period"         description:""` //
	TotalValue    float32   `json:"total_value"    orm:"total_value"    description:""` //
	Change        float32   `json:"change"         orm:"change"         description:""` //
	ChangePercent float32   `json:"change_percent" orm:"change_percent" description:""` //
	TopGainers    string    `json:"top_gainers"    orm:"top_gainers"    description:""` //
	TopLosers     string    `json:"top_losers"     orm:"top_losers"     description:""` //
	BtcBenchmark  float32   `json:"btc_benchmark"  orm:"btc_benchmark"  description:""` //
	EthBenchmark  float32   `json:"eth_benchmark"  orm:"eth_benchmark"  description:""` //
	Content       string    `json:"content"        orm:"content"        description:""` //
	GeneratedAt   time.Time `json:"generated_at"   orm:"generated_at"   description:""` //
}
