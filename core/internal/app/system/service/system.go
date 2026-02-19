// =================================================================================
// 系统管理服务接口定义
// 提供版本信息查询、在线更新、回滚等系统管理能力
// =================================================================================

package service

import (
	"context"

	systemApi "your-finance/allfi/api/v1/system"
)

// ISystem 系统管理服务接口
type ISystem interface {
	// GetVersion 获取当前版本信息
	// 包含版本号、构建时间、Git 提交哈希、运行模式、Go 版本
	GetVersion(ctx context.Context) (*systemApi.GetVersionRes, error)

	// CheckUpdate 检查 GitHub Releases 是否有新版本
	// 通过语义版本比较判断是否需要更新
	CheckUpdate(ctx context.Context) (*systemApi.CheckUpdateRes, error)

	// ApplyUpdate 执行版本更新
	// Docker 模式下转发给 updater 服务，宿主机模式下通过 git checkout 切换版本
	ApplyUpdate(ctx context.Context, targetVersion string) (*systemApi.ApplyUpdateRes, error)

	// Rollback 回滚到指定历史版本
	// 与 ApplyUpdate 类似，但会在历史记录中标记为回滚操作
	Rollback(ctx context.Context, targetVersion string) (*systemApi.RollbackRes, error)

	// GetUpdateStatus 获取当前更新/回滚操作的进度
	// 返回状态机状态：idle/updating/completed/failed
	GetUpdateStatus(ctx context.Context) (*systemApi.GetUpdateStatusRes, error)

	// GetUpdateHistory 获取历史更新记录
	// 从 data/update_history.json 文件中读取
	GetUpdateHistory(ctx context.Context) (*systemApi.GetUpdateHistoryRes, error)

	// GetAPIKeys 获取所有 API Key 配置（脱敏显示）
	GetAPIKeys(ctx context.Context) (*systemApi.GetAPIKeysRes, error)

	// UpdateAPIKey 更新指定服务商的 API Key（加密存储）
	UpdateAPIKey(ctx context.Context, provider string, apiKey string) error

	// DeleteAPIKey 删除指定服务商的 API Key
	DeleteAPIKey(ctx context.Context, provider string) error

	// GetAPIKeyPlain 获取指定服务商的 API Key 明文（内部使用，不对外暴露）
	GetAPIKeyPlain(ctx context.Context, provider string) string
}

// localSystem 系统管理服务实例（延迟注入）
var localSystem ISystem

// System 获取系统管理服务实例
// 如果服务未注册，会触发 panic
func System() ISystem {
	if localSystem == nil {
		panic("ISystem 服务未注册，请检查 logic/system 包的 init 函数")
	}
	return localSystem
}

// RegisterSystem 注册系统管理服务实现
// 由 logic 层在 init 函数中调用
func RegisterSystem(i ISystem) {
	localSystem = i
}
