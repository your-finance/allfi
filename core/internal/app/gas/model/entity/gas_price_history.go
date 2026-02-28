// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// =================================================================================

package entity

import (
	"time"
)

// GasPriceHistory is the golang structure for table gas_price_history.
type GasPriceHistory struct {
	Id         int       `json:"id"          orm:"id"          description:""` //
	CreatedAt  time.Time `json:"created_at"  orm:"created_at"  description:""` //
	Chain      string    `json:"chain"       orm:"chain"       description:""` //
	Low        float32   `json:"low"         orm:"low"         description:""` //
	Standard   float32   `json:"standard"    orm:"standard"    description:""` //
	Fast       float32   `json:"fast"        orm:"fast"        description:""` //
	Instant    float32   `json:"instant"     orm:"instant"     description:""` //
	BaseFee    float32   `json:"base_fee"    orm:"base_fee"    description:""` //
	RecordedAt time.Time `json:"recorded_at" orm:"recorded_at" description:""` //
}
