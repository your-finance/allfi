// Package version 全局版本信息
// 构建时通过 ldflags 注入
package version

import "strings"

var (
	// Version 应用版本号（构建时注入，如 v0.1.7-7-ge7bb477）
	Version = "dev"
	// BuildTime 构建时间
	BuildTime = "unknown"
	// GitCommit Git 提交哈希
	GitCommit = "unknown"
)

// ShortVersion 返回简短版本号
// 对于正式 tag（如 v0.1.7），直接返回去掉 v 前缀的版本
// 对于开发版本（如 v0.1.7-7-ge7bb477），只返回基础 tag 部分（如 0.1.7-dev）
func ShortVersion() string {
	v := strings.TrimPrefix(Version, "v")
	// 如果包含 - 说明是 tag 之后的提交，截取基础版本并加 -dev 后缀
	if base, _, found := strings.Cut(v, "-"); found {
		return base + "-dev"
	}
	return v
}
