// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GasRecommendations is the golang structure of table gas_recommendations for DAO operations like Where/Data.
type GasRecommendations struct {
	g.Meta           `orm:"table:gas_recommendations, do:true"`
	Id               any //
	CreatedAt        any //
	UpdatedAt        any //
	Chain            any //
	RecommendedTime  any //
	EstimatedSavings any //
	Confidence       any //
	ValidUntil       any //
}
