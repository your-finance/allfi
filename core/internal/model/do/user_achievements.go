// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserAchievements is the golang structure of table user_achievements for DAO operations like Where/Data.
type UserAchievements struct {
	g.Meta        `orm:"table:user_achievements, do:true"`
	Id            any //
	CreatedAt     any //
	UpdatedAt     any //
	DeletedAt     any //
	UserId        any //
	AchievementId any //
	UnlockedAt    any //
}
