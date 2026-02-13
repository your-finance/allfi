// Package notification 通知 API 定义
// 提供通知列表、未读计数、标记已读、偏好设置接口
package notification

import "github.com/gogf/gf/v2/frame/g"

// ListReq 获取通知列表请求
type ListReq struct {
	g.Meta   `path:"/notifications" method:"get" summary:"获取通知列表" tags:"通知"`
	Page     int `json:"page" in:"query" d:"1" dc:"页码"`
	PageSize int `json:"page_size" in:"query" d:"20" dc:"每页数量"`
}

// NotificationItem 通知条目
type NotificationItem struct {
	ID        uint   `json:"id" dc:"通知 ID"`
	Type      string `json:"type" dc:"通知类型"`
	Title     string `json:"title" dc:"通知标题"`
	Message   string `json:"message" dc:"通知内容"`
	IsRead    bool   `json:"is_read" dc:"是否已读"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page" dc:"当前页码"`
	PageSize   int   `json:"page_size" dc:"每页数量"`
	Total      int64 `json:"total" dc:"总记录数"`
	TotalPages int   `json:"total_pages" dc:"总页数"`
	HasNext    bool  `json:"has_next" dc:"是否有下一页"`
	HasPrev    bool  `json:"has_prev" dc:"是否有上一页"`
}

// ListRes 获取通知列表响应
type ListRes struct {
	List       []NotificationItem `json:"list" dc:"通知列表"`
	Pagination *PaginationInfo    `json:"pagination" dc:"分页信息"`
}

// GetUnreadCountReq 获取未读通知数量请求
type GetUnreadCountReq struct {
	g.Meta `path:"/notifications/unread-count" method:"get" summary:"获取未读通知数量" tags:"通知"`
}

// GetUnreadCountRes 获取未读通知数量响应
type GetUnreadCountRes struct {
	Count int64 `json:"count" dc:"未读数量"`
}

// MarkReadReq 标记通知已读请求
type MarkReadReq struct {
	g.Meta `path:"/notifications/{id}/read" method:"post" summary:"标记通知为已读" tags:"通知"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"通知 ID"`
}

// MarkReadRes 标记通知已读响应
type MarkReadRes struct{}

// MarkAllReadReq 标记所有通知已读请求
type MarkAllReadReq struct {
	g.Meta `path:"/notifications/read-all" method:"post" summary:"标记所有通知为已读" tags:"通知"`
}

// MarkAllReadRes 标记所有通知已读响应
type MarkAllReadRes struct{}

// GetPreferencesReq 获取通知偏好设置请求
type GetPreferencesReq struct {
	g.Meta `path:"/notifications/preferences" method:"get" summary:"获取通知偏好设置" tags:"通知"`
}

// PreferenceItem 通知偏好设置
type PreferenceItem struct {
	EmailEnabled   bool `json:"email_enabled" dc:"是否启用邮件通知"`
	PushEnabled    bool `json:"push_enabled" dc:"是否启用推送通知"`
	PriceAlert     bool `json:"price_alert" dc:"价格预警通知"`
	PortfolioAlert bool `json:"portfolio_alert" dc:"资产变动通知"`
	SystemNotice   bool `json:"system_notice" dc:"系统通知"`
}

// GetPreferencesRes 获取通知偏好设置响应
type GetPreferencesRes struct {
	Preferences *PreferenceItem `json:"preferences" dc:"偏好设置"`
}

// UpdatePreferencesReq 更新通知偏好设置请求
type UpdatePreferencesReq struct {
	g.Meta         `path:"/notifications/preferences" method:"put" summary:"更新通知偏好设置" tags:"通知"`
	EmailEnabled   *bool `json:"email_enabled" dc:"是否启用邮件通知"`
	PushEnabled    *bool `json:"push_enabled" dc:"是否启用推送通知"`
	PriceAlert     *bool `json:"price_alert" dc:"价格预警通知"`
	PortfolioAlert *bool `json:"portfolio_alert" dc:"资产变动通知"`
	SystemNotice   *bool `json:"system_notice" dc:"系统通知"`
}

// UpdatePreferencesRes 更新通知偏好设置响应
type UpdatePreferencesRes struct {
	Preferences *PreferenceItem `json:"preferences" dc:"更新后的偏好设置"`
}
