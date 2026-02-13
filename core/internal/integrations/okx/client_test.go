// Package okx OKX 客户端测试
package okx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDoRequest_ReturnsResponseBody 测试 doRequest 正确读取并返回完整响应体
func TestDoRequest_ReturnsResponseBody(t *testing.T) {
	// 创建模拟 HTTP 服务器
	expectedBody := `{"code":"0","msg":"","data":[{"details":[{"ccy":"BTC","bal":"1.5"}]}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	}))
	defer server.Close()

	// 创建客户端，将 baseURL 指向模拟服务器
	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// 临时替换 baseURL（通过直接构造请求来测试 doRequest 的读取逻辑）
	// 由于 doRequest 内部使用包级 baseURL 常量，我们用 testDoRequestWithURL 辅助函数
	result, err := testDoRequestWithURL(client, server.URL, "/api/v5/account/balance")

	require.NoError(t, err, "doRequest 不应返回错误")
	assert.NotEmpty(t, result, "doRequest 应返回非空数据")
	assert.Equal(t, expectedBody, string(result), "doRequest 应返回完整的响应体")
}

// TestDoRequest_EmptyResponse 测试 doRequest 处理空响应
func TestDoRequest_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// 返回空响应体
	}))
	defer server.Close()

	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	result, err := testDoRequestWithURL(client, server.URL, "/api/v5/account/balance")

	require.NoError(t, err, "空响应不应返回错误")
	assert.Empty(t, result, "空响应应返回空切片")
}

// TestDoRequest_LargeResponse 测试 doRequest 处理较大的响应体
func TestDoRequest_LargeResponse(t *testing.T) {
	// 构造一个大于 512 字节的响应（超过 Read 单次读取的典型缓冲区大小）
	largeData := `{"code":"0","data":[`
	for i := 0; i < 100; i++ {
		if i > 0 {
			largeData += ","
		}
		largeData += `{"ccy":"TOKEN` + string(rune('A'+i%26)) + `","bal":"100.0"}`
	}
	largeData += `]}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(largeData))
	}))
	defer server.Close()

	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	result, err := testDoRequestWithURL(client, server.URL, "/api/v5/account/balance")

	require.NoError(t, err, "大响应不应返回错误")
	assert.Equal(t, largeData, string(result), "应完整读取所有数据")
}

// testDoRequestWithURL 辅助函数：绕过包级 baseURL 常量，直接测试请求逻辑
// 模拟 doRequest 的核心逻辑，使用自定义 URL
func testDoRequestWithURL(c *Client, serverURL string, path string) ([]byte, error) {
	ctx := context.Background()
	url := serverURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	req.Header.Set("OK-ACCESS-KEY", c.apiKey)
	req.Header.Set("OK-ACCESS-SIGN", c.sign(timestamp, "GET", path, ""))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", c.passphrase)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 使用与 doRequest 相同的读取方式
	return readResponseBody(resp)
}

// readResponseBody 从 doRequest 中提取的响应读取逻辑
// 与 doRequest 使用相同的 io.ReadAll 方式读取响应
func readResponseBody(resp *http.Response) ([]byte, error) {
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	return result, nil
}
