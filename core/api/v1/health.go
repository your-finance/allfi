// Package v1 健康检查 API 定义
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetHealthReq 获取健康状态请求
type GetHealthReq struct {
	g.Meta `path:"/health" method:"get" summary:"获取服务健康状态" tags:"系统"`
}

// GetHealthRes 获取健康状态响应
type GetHealthRes struct {
	Status   string                 `json:"status" dc:"服务状态"`
	Name     string                 `json:"name" dc:"服务名称"`
	Version  string                 `json:"version" dc:"服务版本"`
	Env      string                 `json:"env" dc:"运行环境"`
	Database map[string]interface{} `json:"database" dc:"数据库状态"`
}
