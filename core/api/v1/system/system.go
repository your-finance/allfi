// Package system 系统管理 API 定义
// 提供版本信息查询、在线更新、回滚等系统管理功能
package system

import "github.com/gogf/gf/v2/frame/g"

// GetVersionReq 获取版本信息请求
type GetVersionReq struct {
	g.Meta `path:"/system/version" method:"get" summary:"获取版本信息" tags:"系统管理" security:""`
}

// GetVersionRes 获取版本信息响应
type GetVersionRes struct {
	Version   string `json:"version" dc:"当前版本号"`
	BuildTime string `json:"build_time" dc:"构建时间"`
	GitCommit string `json:"git_commit" dc:"Git 提交哈希"`
	RunMode   string `json:"run_mode" dc:"运行模式（docker/host）"`
	GoVersion string `json:"go_version" dc:"Go 版本"`
}

// CheckUpdateReq 检查更新请求
type CheckUpdateReq struct {
	g.Meta `path:"/system/update/check" method:"get" summary:"检查新版本" tags:"系统管理"`
}

// CheckUpdateRes 检查更新响应
type CheckUpdateRes struct {
	HasUpdate      bool   `json:"has_update" dc:"是否有新版本"`
	CurrentVersion string `json:"current_version" dc:"当前版本"`
	LatestVersion  string `json:"latest_version" dc:"最新版本"`
	ReleaseNotes   string `json:"release_notes" dc:"更新说明"`
	ReleaseURL     string `json:"release_url" dc:"发布页面 URL"`
	PublishedAt    string `json:"published_at" dc:"发布时间"`
}

// ApplyUpdateReq 执行更新请求
type ApplyUpdateReq struct {
	g.Meta        `path:"/system/update/apply" method:"post" summary:"执行更新" tags:"系统管理"`
	TargetVersion string `json:"target_version" v:"required" dc:"目标版本号"`
}

// ApplyUpdateRes 执行更新响应
type ApplyUpdateRes struct {
	Status  string `json:"status" dc:"更新状态（started/failed）"`
	Message string `json:"message" dc:"状态信息"`
}

// RollbackReq 版本回滚请求
type RollbackReq struct {
	g.Meta        `path:"/system/update/rollback" method:"post" summary:"版本回滚" tags:"系统管理"`
	TargetVersion string `json:"target_version" v:"required" dc:"目标回滚版本"`
}

// RollbackRes 版本回滚响应
type RollbackRes struct {
	Status  string `json:"status" dc:"回滚状态"`
	Message string `json:"message" dc:"状态信息"`
}

// GetUpdateStatusReq 获取更新状态请求
type GetUpdateStatusReq struct {
	g.Meta `path:"/system/update/status" method:"get" summary:"获取更新进度" tags:"系统管理"`
}

// GetUpdateStatusRes 获取更新状态响应
type GetUpdateStatusRes struct {
	State    string `json:"state" dc:"状态（idle/updating/completed/failed）"`
	Step     int    `json:"step" dc:"当前步骤"`
	Total    int    `json:"total" dc:"总步骤数"`
	StepName string `json:"step_name" dc:"当前步骤名称"`
	Message  string `json:"message" dc:"详细信息"`
}

// GetUpdateHistoryReq 获取更新历史请求
type GetUpdateHistoryReq struct {
	g.Meta `path:"/system/update/history" method:"get" summary:"更新历史" tags:"系统管理"`
}

// UpdateRecord 更新记录
type UpdateRecord struct {
	Version   string `json:"version" dc:"版本号"`
	GitCommit string `json:"git_commit" dc:"Git 提交"`
	Timestamp string `json:"timestamp" dc:"更新时间"`
	Status    string `json:"status" dc:"状态（success/failed/rolled_back）"`
}

// GetUpdateHistoryRes 获取更新历史响应
type GetUpdateHistoryRes struct {
	Records []UpdateRecord `json:"records" dc:"更新记录列表"`
}
