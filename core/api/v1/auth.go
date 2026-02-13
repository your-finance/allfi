package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/model/entity"
)

// RegisterReq 用户注册请求
type RegisterReq struct {
	g.Meta   `path:"/auth/register" method:"post" summary:"用户注册" tags:"认证"`
	Username string `json:"username" v:"required|length:3,50" dc:"用户名"`
	Email    string `json:"email" v:"required|email" dc:"邮箱"`
	Password string `json:"password" v:"required|length:6,50" dc:"密码"`
	Nickname string `json:"nickname" v:"max-length:50" dc:"昵称（可选）"`
}

// RegisterRes 用户注册响应
type RegisterRes struct {
	User  *entity.User `json:"user" dc:"用户信息"`
	Token string       `json:"token" dc:"JWT Token"`
}

// LoginReq 用户登录请求
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" summary:"用户登录" tags:"认证"`
	Username string `json:"username" v:"required" dc:"用户名或邮箱"`
	Password string `json:"password" v:"required" dc:"密码"`
}

// LoginRes 用户登录响应
type LoginRes struct {
	User  *entity.User `json:"user" dc:"用户信息"`
	Token string       `json:"token" dc:"JWT Token"`
}

// GetProfileReq 获取用户信息请求
type GetProfileReq struct {
	g.Meta `path:"/auth/profile" method:"get" summary:"获取当前用户信息" tags:"认证"`
}

// GetProfileRes 获取用户信息响应
type GetProfileRes struct {
	User *entity.User `json:"user" dc:"用户信息"`
}

// UpdateProfileReq 更新用户信息请求
type UpdateProfileReq struct {
	g.Meta   `path:"/auth/profile" method:"put" summary:"更新当前用户信息" tags:"认证"`
	Nickname string `json:"nickname" v:"max-length:50" dc:"昵称"`
	Avatar   string `json:"avatar" v:"url|max-length:255" dc:"头像URL"`
}

// UpdateProfileRes 更新用户信息响应
type UpdateProfileRes struct {
	User *entity.User `json:"user" dc:"用户信息"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	g.Meta      `path:"/auth/password" method:"put" summary:"修改密码" tags:"认证"`
	OldPassword string `json:"old_password" v:"required" dc:"旧密码"`
	NewPassword string `json:"new_password" v:"required|length:6,50" dc:"新密码"`
}

// ChangePasswordRes 修改密码响应
type ChangePasswordRes struct {
	Success bool `json:"success" dc:"是否成功"`
}
