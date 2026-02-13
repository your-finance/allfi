// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// ExchangeRates is the golang structure for table exchange_rates.
type ExchangeRates struct {
	Id           int       `json:"id"            orm:"id"            description:""` //
	CreatedAt    time.Time `json:"created_at"    orm:"created_at"    description:""` //
	UpdatedAt    time.Time `json:"updated_at"    orm:"updated_at"    description:""` //
	DeletedAt    time.Time `json:"deleted_at"    orm:"deleted_at"    description:""` //
	FromCurrency string    `json:"from_currency" orm:"from_currency" description:""` //
	ToCurrency   string    `json:"to_currency"   orm:"to_currency"   description:""` //
	Rate         float64   `json:"rate"          orm:"rate"          description:""` //
	Source       string    `json:"source"        orm:"source"        description:""` //
	FetchedAt    time.Time `json:"fetched_at"    orm:"fetched_at"    description:""` //
}
