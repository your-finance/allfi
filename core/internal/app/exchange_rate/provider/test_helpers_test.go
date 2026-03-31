package provider

import (
	"os"
	"testing"
)

// requireOnlineProviderTests 用于显式保护依赖公网的集成测试。
// 默认 `go test ./...` 不跑这类测试，避免 CI 或离线环境超时。
func requireOnlineProviderTests(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("跳过外网集成测试")
	}
	if os.Getenv("ALLFI_RUN_ONLINE_TESTS") != "1" {
		t.Skip("跳过外网集成测试；如需执行请设置 ALLFI_RUN_ONLINE_TESTS=1")
	}
}
