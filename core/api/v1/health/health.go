// Package health 健康检查 API 定义
package health

import "github.com/gogf/gf/v2/frame/g"

// GetHealthReq 健康检查请求
type GetHealthReq struct {
	g.Meta `path:"/health" method:"get" summary:"健康检查" tags:"系统"`
}

// GetHealthRes 健康检查响应
type GetHealthRes struct {
	Status    string `json:"status" dc:"服务状态"`
	Version   string `json:"version" dc:"服务版本"`
	Timestamp int64  `json:"timestamp" dc:"时间戳（Unix 秒）"`
}
