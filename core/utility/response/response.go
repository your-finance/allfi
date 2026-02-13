// Package response 提供统一的API响应格式
package response

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Response 标准响应结构
type Response struct {
	Code      int         `json:"code" dc:"状态码(0=成功,非0=错误)"`
	Message   string      `json:"message" dc:"响应消息"`
	Data      interface{} `json:"data,omitempty" dc:"响应数据"`
	Timestamp int64       `json:"timestamp" dc:"时间戳"`
}

// Success 成功响应（200）
//
// 参数:
//   - r: HTTP请求对象
//   - data: 响应数据（可选）
func Success(r *ghttp.Request, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	r.Response.WriteJson(&Response{
		Code:      0,
		Message:   "success",
		Data:      responseData,
		Timestamp: time.Now().Unix(),
	})
}

// SuccessWithMessage 成功响应（自定义消息）
//
// 参数:
//   - r: HTTP请求对象
//   - message: 自定义成功消息
//   - data: 响应数据（可选）
func SuccessWithMessage(r *ghttp.Request, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	r.Response.WriteJson(&Response{
		Code:      0,
		Message:   message,
		Data:      responseData,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
//
// 参数:
//   - r: HTTP请求对象
//   - code: 错误码
//   - message: 错误消息
func Error(r *ghttp.Request, code int, message string) {
	ctx := r.Context()

	// 记录错误日志
	g.Log().Error(ctx, "API错误响应",
		"code", code,
		"message", message,
		"path", r.URL.Path,
		"method", r.Method,
	)

	r.Response.WriteJson(&Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// ErrorWithData 错误响应（带数据）
//
// 参数:
//   - r: HTTP请求对象
//   - code: 错误码
//   - message: 错误消息
//   - data: 附加数据
func ErrorWithData(r *ghttp.Request, code int, message string, data interface{}) {
	ctx := r.Context()

	g.Log().Error(ctx, "API错误响应",
		"code", code,
		"message", message,
		"data", data,
		"path", r.URL.Path,
	)

	r.Response.WriteJson(&Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// 常用错误码定义
const (
	// 客户端错误 (1xxx)
	CodeInvalidParams   = 1001 // 参数错误
	CodeUnauthorized    = 1002 // 未授权
	CodeForbidden       = 1003 // 禁止访问
	CodeNotFound        = 1004 // 资源不存在
	CodeConflict        = 1005 // 资源冲突
	CodeTooManyRequests = 1006 // 请求过多

	// 服务器错误 (2xxx)
	CodeInternalError    = 2001 // 内部错误
	CodeDatabaseError    = 2002 // 数据库错误
	CodeExternalAPIError = 2003 // 外部API错误
	CodeEncryptionError  = 2004 // 加密错误
	CodeCacheError       = 2005 // 缓存错误

	// 业务错误 (3xxx)
	CodeExchangeAPIError  = 3001 // 交易所API错误
	CodeBlockchainError   = 3002 // 区块链查询错误
	CodePriceServiceError = 3003 // 价格服务错误
	CodeAssetNotFound     = 3004 // 资产不存在
	CodeAccountNotFound   = 3005 // 账户不存在
)

// 便捷错误响应方法

// InvalidParams 参数错误响应
func InvalidParams(r *ghttp.Request, message ...string) {
	msg := "参数错误"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(r, CodeInvalidParams, msg)
}

// Unauthorized 未授权响应
func Unauthorized(r *ghttp.Request, message ...string) {
	msg := "未授权访问"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(r, CodeUnauthorized, msg)
}

// NotFound 资源不存在响应
func NotFound(r *ghttp.Request, message ...string) {
	msg := "资源不存在"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(r, CodeNotFound, msg)
}

// InternalError 内部错误响应
func InternalError(r *ghttp.Request, message ...string) {
	msg := "服务器内部错误"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(r, CodeInternalError, msg)
}

// DatabaseError 数据库错误响应
func DatabaseError(r *ghttp.Request, message ...string) {
	msg := "数据库操作失败"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(r, CodeDatabaseError, msg)
}

// ExternalAPIError 外部API错误响应
func ExternalAPIError(r *ghttp.Request, message ...string) {
	msg := "外部服务调用失败"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(r, CodeExternalAPIError, msg)
}
