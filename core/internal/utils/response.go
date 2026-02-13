// Package utils 统一响应格式
// 提供标准化的 API 响应结构和工具函数
package utils

import (
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	jsoniter "github.com/json-iterator/go"
)

// 错误码定义
const (
	// 成功
	CodeSuccess = 0

	// 客户端错误 (1001-1999)
	CodeInvalidParams    = 1001 // 参数错误
	CodeValidationFailed = 1002 // 验证失败
	CodeResourceNotFound = 1003 // 资源不存在
	CodeDuplicateEntry   = 1004 // 重复条目
	CodeUnauthorized     = 1005 // 未授权

	// 服务器错误 (2001-2999)
	CodeInternalError    = 2001 // 内部错误
	CodeDatabaseError    = 2002 // 数据库错误
	CodeExternalAPIError = 2003 // 外部 API 错误
	CodeEncryptionError  = 2004 // 加密错误
	CodeCacheError       = 2005 // 缓存错误
	CodeConfigError      = 2006 // 配置错误

	// 业务错误 (3001-3999)
	CodeExchangeAPIError   = 3001 // 交易所 API 错误
	CodeBlockchainAPIError = 3002 // 区块链 API 错误
	CodeInvalidAPIKey      = 3003 // API Key 无效
	CodeInvalidAddress     = 3004 // 地址无效
	CodeBalanceFetchFailed = 3005 // 余额获取失败
	CodeRateFetchFailed    = 3006 // 汇率获取失败
	CodeSnapshotFailed     = 3007 // 快照创建失败
)

// 错误消息映射
var ErrorMessages = map[int]string{
	CodeSuccess:            "成功",
	CodeInvalidParams:      "请求参数错误",
	CodeValidationFailed:   "数据验证失败",
	CodeResourceNotFound:   "资源不存在",
	CodeDuplicateEntry:     "数据已存在",
	CodeUnauthorized:       "未授权访问",
	CodeInternalError:      "服务器内部错误",
	CodeDatabaseError:      "数据库操作失败",
	CodeExternalAPIError:   "外部服务调用失败",
	CodeEncryptionError:    "数据加密/解密失败",
	CodeCacheError:         "缓存操作失败",
	CodeConfigError:        "配置加载失败",
	CodeExchangeAPIError:   "交易所 API 调用失败",
	CodeBlockchainAPIError: "区块链 API 调用失败",
	CodeInvalidAPIKey:      "API Key 无效",
	CodeInvalidAddress:     "钱包地址格式无效",
	CodeBalanceFetchFailed: "获取余额失败",
	CodeRateFetchFailed:    "获取汇率失败",
	CodeSnapshotFailed:     "创建资产快照失败",
}

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// PaginationData 分页数据结构
type PaginationData struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// NewResponse 创建新的响应
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// Success 成功响应
func Success(r *ghttp.Request, data interface{}) {
	response := NewResponse(CodeSuccess, "成功", data)
	r.Response.WriteJson(response)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(r *ghttp.Request, message string, data interface{}) {
	response := NewResponse(CodeSuccess, message, data)
	r.Response.WriteJson(response)
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(r *ghttp.Request, list interface{}, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize
	pagination := PaginationData{
		List: list,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}
	response := NewResponse(CodeSuccess, "成功", pagination)
	r.Response.WriteJson(response)
}

// Error 错误响应
func Error(r *ghttp.Request, code int, message string) {
	if message == "" {
		if msg, ok := ErrorMessages[code]; ok {
			message = msg
		} else {
			message = "未知错误"
		}
	}
	response := NewResponse(code, message, nil)
	r.Response.WriteJson(response)
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(r *ghttp.Request, code int, message string, data interface{}) {
	if message == "" {
		if msg, ok := ErrorMessages[code]; ok {
			message = msg
		} else {
			message = "未知错误"
		}
	}
	response := NewResponse(code, message, data)
	r.Response.WriteJson(response)
}

// ErrorInternal 内部错误响应
func ErrorInternal(r *ghttp.Request, err error) {
	message := "服务器内部错误"
	if err != nil {
		// 生产环境不应该暴露具体错误信息
		// 这里仅用于开发调试，生产环境应记录日志但不返回
		message = err.Error()
	}
	Error(r, CodeInternalError, message)
}

// ErrorNotFound 资源不存在响应
func ErrorNotFound(r *ghttp.Request, resource string) {
	message := resource + " 不存在"
	Error(r, CodeResourceNotFound, message)
}

// ErrorValidation 验证失败响应
func ErrorValidation(r *ghttp.Request, err error) {
	message := "数据验证失败"
	if err != nil {
		message = err.Error()
	}
	Error(r, CodeValidationFailed, message)
}

// ErrorDatabase 数据库错误响应
func ErrorDatabase(r *ghttp.Request, err error) {
	message := "数据库操作失败"
	// 生产环境应记录日志但不返回具体错误
	Error(r, CodeDatabaseError, message)
}

// GetErrorMessage 获取错误码对应的消息
func GetErrorMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// 用于 net/http 的简化错误码常量
const (
	ErrInvalidParams   = CodeInvalidParams
	ErrNotFound        = CodeResourceNotFound
	ErrInternalServer  = CodeInternalError
	ErrExternalAPI     = CodeExternalAPIError
	ErrUnauthorized    = CodeUnauthorized
	ErrDatabaseError   = CodeDatabaseError
	ErrEncryptionError = CodeEncryptionError
)

// SuccessResponse 创建成功响应（用于 net/http）
func SuccessResponse(data interface{}, message string) *Response {
	if message == "" {
		message = "成功"
	}
	return NewResponse(CodeSuccess, message, data)
}

// ErrorResponse 创建错误响应（用于 net/http）
func ErrorResponse(code int, message string) *Response {
	if message == "" {
		message = GetErrorMessage(code)
	}
	return NewResponse(code, message, nil)
}

// ToJSON 将响应转换为 JSON 字节
func (r *Response) ToJSON() ([]byte, error) {
	return jsoniter.Marshal(r)
}
