// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// UserAchievements is the golang structure for table user_achievements.
type UserAchievements struct {
	Id            int       `json:"id"             orm:"id"             description:""` //
	CreatedAt     time.Time `json:"created_at"     orm:"created_at"     description:""` //
	UpdatedAt     time.Time `json:"updated_at"     orm:"updated_at"     description:""` //
	DeletedAt     time.Time `json:"deleted_at"     orm:"deleted_at"     description:""` //
	UserId        int       `json:"user_id"        orm:"user_id"        description:""` //
	AchievementId string    `json:"achievement_id" orm:"achievement_id" description:""` //
	UnlockedAt    time.Time `json:"unlocked_at"    orm:"unlocked_at"    description:""` //
}
