// Package controller 目标追踪模块 - 路由和控制器
// 绑定目标追踪 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	goalApi "your-finance/allfi/api/v1/goal"
	"your-finance/allfi/internal/app/goal/service"
)

// Controller 目标追踪控制器
type Controller struct{}

// Register 注册目标追踪模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// List 获取目标列表
func (c *Controller) List(ctx context.Context, req *goalApi.ListReq) (res *goalApi.ListRes, err error) {
	goals, err := service.Goal().GetGoals(ctx)
	if err != nil {
		return nil, err
	}

	return &goalApi.ListRes{
		Goals: goals,
	}, nil
}

// Create 创建目标
func (c *Controller) Create(ctx context.Context, req *goalApi.CreateReq) (res *goalApi.CreateRes, err error) {
	goal, err := service.Goal().CreateGoal(ctx, req)
	if err != nil {
		return nil, err
	}

	return &goalApi.CreateRes{
		Goal: goal,
	}, nil
}

// Update 更新目标
func (c *Controller) Update(ctx context.Context, req *goalApi.UpdateReq) (res *goalApi.UpdateRes, err error) {
	goal, err := service.Goal().UpdateGoal(ctx, req)
	if err != nil {
		return nil, err
	}

	return &goalApi.UpdateRes{
		Goal: goal,
	}, nil
}

// Delete 删除目标
func (c *Controller) Delete(ctx context.Context, req *goalApi.DeleteReq) (res *goalApi.DeleteRes, err error) {
	err = service.Goal().DeleteGoal(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &goalApi.DeleteRes{}, nil
}
