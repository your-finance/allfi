// Package middleware 统一响应中间件
// 拦截 Controller 返回值，将其包装为统一的 JSON 响应格式：
// {"code": 0, "message": "success", "data": {...}, "timestamp": 1234567890}
package middleware

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// defaultResponse 统一响应结构
// 与 internal/utils/response.go 和 utility/response/response.go 中定义的格式保持一致
type defaultResponse struct {
	Code      int         `json:"code"`          // 状态码（0=成功，非0=错误）
	Message   string      `json:"message"`       // 响应消息
	Data      interface{} `json:"data,omitempty"` // 响应数据
	Timestamp int64       `json:"timestamp"`     // Unix 秒级时间戳
}

// MiddlewareHandlerResponse 统一响应处理中间件
// 拦截所有 Controller 的返回值，包装为项目统一的 JSON 格式
//
// 处理逻辑：
//  1. 先调用 r.Middleware.Next() 执行后续处理器
//  2. 如果响应已经被写入（如手动调用 WriteJson），则不再处理
//  3. 如果存在错误（r.GetError()），提取错误码映射到 HTTP 状态码，返回错误 JSON
//  4. 如果是成功响应，将 r.GetHandlerResponse() 包装在统一格式中
func MiddlewareHandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果响应已被写入（Handler 中手动写入了响应），则跳过统一包装
	if r.Response.BufferLength() > 0 {
		return
	}

	// 设置响应 Content-Type
	r.Response.Header().Set("Content-Type", "application/json")

	// 检查是否存在错误
	if err := r.GetError(); err != nil {
		// 清除错误，避免框架重复处理
		r.SetError(nil)

		// 根据错误构建错误响应
		handleErrorResponse(r, err)
		return
	}

	// 获取 Handler 返回的数据
	res := r.GetHandlerResponse()

	// 构建成功响应
	r.Response.WriteJson(defaultResponse{
		Code:      0,
		Message:   "success",
		Data:      res,
		Timestamp: time.Now().Unix(),
	})
}

// handleErrorResponse 处理错误响应
// 从 gerror 中提取错误码，映射为对应的 HTTP 状态码和业务错误码
func handleErrorResponse(r *ghttp.Request, err error) {
	// 提取 GoFrame 错误码
	code := gerror.Code(err)
	bizCode := code.Code()
	message := err.Error()

	// 如果没有明确的错误码，默认使用内部错误（2001）
	if bizCode == gcode.CodeNil.Code() {
		bizCode = 2001 // CodeInternalError
	}

	// 映射 HTTP 状态码
	httpStatus := mapBizCodeToHTTPStatus(bizCode)
	r.Response.WriteStatus(httpStatus)

	// 如果错误码自带消息且不为空，优先使用错误码消息
	if code.Message() != "" {
		message = code.Message()
	}

	// 写入错误响应 JSON
	r.Response.WriteJson(defaultResponse{
		Code:      bizCode,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// mapBizCodeToHTTPStatus 将业务错误码映射为 HTTP 状态码
//
// 本项目有两套错误码定义（正在迁移统一中）：
//   - internal/utils/response.go: 1001-参数错误, 1002-验证失败, 1003-不存在, 1004-重复, 1005-未授权
//   - utility/response/response.go: 1001-参数错误, 1002-未授权, 1003-禁止, 1004-不存在, 1005-冲突, 1006-限流
//
// 映射规则（兼容两套错误码）：
//   - 1001: 400 Bad Request（参数错误）
//   - 1002: 401 Unauthorized（未授权 - utility 定义）
//   - 1003: 404 Not Found（资源不存在 - utils 定义）/ 403 Forbidden（utility 定义）
//   - 1004: 409 Conflict（重复条目 - utils 定义）/ 404 Not Found（utility 定义）
//   - 1005: 401 Unauthorized（未授权 - utils 定义）/ 409 Conflict（utility 定义）
//   - 1006: 429 Too Many Requests（请求过多）
//   - 2001-2999: 500 Internal Server Error（服务端错误）
//   - 3001-3999: 500 Internal Server Error（业务错误）
//   - 其他: 500 Internal Server Error
func mapBizCodeToHTTPStatus(bizCode int) int {
	// 精确匹配高优先级错误码
	switch bizCode {
	case 1001: // 参数错误（两套定义一致）
		return http.StatusBadRequest
	case 1002: // 未授权（utility/response）/ 验证失败（utils/response）
		return http.StatusUnauthorized
	case 1003: // 资源不存在（utils/response）/ 禁止访问（utility/response）
		return http.StatusNotFound
	case 1004: // 重复条目（utils/response）/ 资源不存在（utility/response）
		return http.StatusConflict
	case 1005: // 未授权（utils/response）/ 资源冲突（utility/response）
		return http.StatusUnauthorized
	case 1006: // 请求过多
		return http.StatusTooManyRequests
	}

	// 范围匹配
	switch {
	case bizCode >= 1001 && bizCode <= 1999:
		// 其余客户端错误统一返回 400
		return http.StatusBadRequest
	case bizCode >= 2001 && bizCode <= 2999:
		// 服务端错误
		return http.StatusInternalServerError
	case bizCode >= 3001 && bizCode <= 3999:
		// 业务错误
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
