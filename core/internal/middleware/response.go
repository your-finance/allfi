// Package middleware 统一响应中间件
// 拦截 Controller 返回值，将其包装为统一的 JSON 响应格式：
// {"code": 0, "message": "success", "data": {...}, "timestamp": 1234567890}
//
// 本中间件参照 GoFrame v2 官方 ghttp.MiddlewareHandlerResponse 实现，
// 在其基础上增加了 timestamp 字段
package middleware

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// defaultResponse 统一响应结构
type defaultResponse struct {
	Code      int         `json:"code"`      // 状态码（0=成功，非0=错误）
	Message   string      `json:"message"`   // 响应消息
	Data      interface{} `json:"data"`      // 响应数据（成功时为具体数据，错误时为 nil）
	Timestamp int64       `json:"timestamp"` // Unix 秒级时间戳
}

// MiddlewareHandlerResponse 统一响应处理中间件
// 参照 GoFrame v2 官方实现（ghttp.MiddlewareHandlerResponse），
// 拦截 Controller 返回的 (res, err)，包装为项目统一的 JSON 格式。
//
// 处理逻辑：
//  1. 调用 r.Middleware.Next() 执行后续处理器
//  2. 如果响应已被写入（Handler 中手动调用了 WriteJson/WriteJsonExit），则不再处理
//  3. 如果存在错误（r.GetError()），提取 gerror 错误码，返回错误 JSON
//  4. 如果是成功响应，将 r.GetHandlerResponse() 包装在统一格式中
//  5. 处理非 200 HTTP 状态码（404、403 等）
func MiddlewareHandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果响应已被写入（Handler 中手动写入了响应，如 auth.go 的 WriteJsonExit），则跳过统一包装
	// BufferLength：检查缓冲区中待发送的字节
	// BytesWritten：检查已直接发送到客户端的字节（绕过缓冲区的写入）
	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)

	if err != nil {
		// 有错误：提取错误码和错误消息
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		msg = err.Error()
	} else {
		// 无错误：检查 HTTP 状态码（处理 404、403 等框架级错误）
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			msg = http.StatusText(r.Response.Status)
		} else {
			// 成功响应
			code = gcode.CodeOK
			msg = "success"
		}
	}

	r.Response.WriteJson(defaultResponse{
		Code:      code.Code(),
		Message:   msg,
		Data:      res,
		Timestamp: time.Now().Unix(),
	})
}
