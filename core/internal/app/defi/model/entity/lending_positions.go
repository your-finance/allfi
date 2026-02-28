// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:16:43
// =================================================================================

package entity

import (
	"time"
)

// LendingPositions is the golang structure for table lending_positions.
type LendingPositions struct {
	Id                   int       `json:"id"                    orm:"id"                    description:""` //
	UserId               int       `json:"user_id"               orm:"user_id"               description:""` //
	Protocol             string    `json:"protocol"              orm:"protocol"              description:""` //
	Chain                string    `json:"chain"                 orm:"chain"                 description:""` //
	WalletAddress        string    `json:"wallet_address"        orm:"wallet_address"        description:""` //
	SupplyToken          string    `json:"supply_token"          orm:"supply_token"          description:""` //
	SupplyAmount         float64   `json:"supply_amount"         orm:"supply_amount"         description:""` //
	SupplyValueUsd       float64   `json:"supply_value_usd"      orm:"supply_value_usd"      description:""` //
	SupplyApy            float64   `json:"supply_apy"            orm:"supply_apy"            description:""` //
	BorrowToken          string    `json:"borrow_token"          orm:"borrow_token"          description:""` //
	BorrowAmount         float64   `json:"borrow_amount"         orm:"borrow_amount"         description:""` //
	BorrowValueUsd       float64   `json:"borrow_value_usd"      orm:"borrow_value_usd"      description:""` //
	BorrowApy            float64   `json:"borrow_apy"            orm:"borrow_apy"            description:""` //
	HealthFactor         float64   `json:"health_factor"         orm:"health_factor"         description:""` //
	LiquidationThreshold float64   `json:"liquidation_threshold" orm:"liquidation_threshold" description:""` //
	Ltv                  float64   `json:"ltv"                   orm:"ltv"                   description:""` //
	NetApy               float64   `json:"net_apy"               orm:"net_apy"               description:""` //
	CreatedAt            time.Time `json:"created_at"            orm:"created_at"            description:""` //
	UpdatedAt            time.Time `json:"updated_at"            orm:"updated_at"            description:""` //
}
