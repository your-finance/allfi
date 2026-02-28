// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:16:43
// =================================================================================

package entity

import (
	"time"
)

// LendingRateHistory is the golang structure for table lending_rate_history.
type LendingRateHistory struct {
	Id                int       `json:"id"                  orm:"id"                  description:""` //
	Protocol          string    `json:"protocol"            orm:"protocol"            description:""` //
	Chain             string    `json:"chain"               orm:"chain"               description:""` //
	Token             string    `json:"token"               orm:"token"               description:""` //
	SupplyApy         float64   `json:"supply_apy"          orm:"supply_apy"          description:""` //
	BorrowApyStable   float64   `json:"borrow_apy_stable"   orm:"borrow_apy_stable"   description:""` //
	BorrowApyVariable float64   `json:"borrow_apy_variable" orm:"borrow_apy_variable" description:""` //
	UtilizationRate   float64   `json:"utilization_rate"    orm:"utilization_rate"    description:""` //
	RecordedAt        time.Time `json:"recorded_at"         orm:"recorded_at"         description:""` //
}
