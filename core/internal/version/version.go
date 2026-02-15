// Package version 全局版本信息
// 构建时通过 ldflags 注入
package version

var (
	// Version 应用版本号
	Version = "dev"
	// BuildTime 构建时间
	BuildTime = "unknown"
	// GitCommit Git 提交哈希
	GitCommit = "unknown"
)
