// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// =================================================================================

package entity

import (
	"time"
)

// GasRecommendations is the golang structure for table gas_recommendations.
type GasRecommendations struct {
	Id               int       `json:"id"                orm:"id"                description:""` //
	CreatedAt        time.Time `json:"created_at"        orm:"created_at"        description:""` //
	UpdatedAt        time.Time `json:"updated_at"        orm:"updated_at"        description:""` //
	Chain            string    `json:"chain"             orm:"chain"             description:""` //
	RecommendedTime  string    `json:"recommended_time"  orm:"recommended_time"  description:""` //
	EstimatedSavings float32   `json:"estimated_savings" orm:"estimated_savings" description:""` //
	Confidence       float32   `json:"confidence"        orm:"confidence"        description:""` //
	ValidUntil       time.Time `json:"valid_until"       orm:"valid_until"       description:""` //
}
