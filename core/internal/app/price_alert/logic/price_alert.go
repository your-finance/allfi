// Package logic 价格预警业务逻辑
// 实现价格预警的增删改查和触发检查
package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	priceAlertApi "your-finance/allfi/api/v1/price_alert"
	alertModel "your-finance/allfi/internal/app/price_alert/model"
	notificationService "your-finance/allfi/internal/app/notification/service"
	"your-finance/allfi/internal/app/price_alert/dao"
	"your-finance/allfi/internal/app/price_alert/service"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/model/entity"
)

// sPriceAlert 价格预警服务实现
type sPriceAlert struct{}

// New 创建价格预警服务实例
func New() service.IPriceAlert {
	return &sPriceAlert{}
}

// CreateAlert 创建价格预警
func (s *sPriceAlert) CreateAlert(ctx context.Context, userID int, req *priceAlertApi.CreateReq) (*priceAlertApi.AlertItem, error) {
	// 验证条件类型
	if !alertModel.IsValidCondition(req.Condition) {
		return nil, gerror.Newf("无效的预警条件: %s，仅支持 above/below", req.Condition)
	}

	// 目标价格必须大于 0
	if req.TargetPrice <= 0 {
		return nil, gerror.New("目标价格必须大于 0")
	}

	// 将币种转为大写
	symbol := strings.ToUpper(req.Symbol)

	// 插入数据库
	now := gtime.Now()
	result, err := dao.PriceAlerts.Ctx(ctx).Data(g.Map{
		dao.PriceAlerts.Columns().UserId:      userID,
		dao.PriceAlerts.Columns().Symbol:      symbol,
		dao.PriceAlerts.Columns().Condition:   req.Condition,
		dao.PriceAlerts.Columns().TargetPrice: req.TargetPrice,
		dao.PriceAlerts.Columns().Note:        req.Note,
		dao.PriceAlerts.Columns().IsActive:    1,
		dao.PriceAlerts.Columns().Triggered:   0,
		dao.PriceAlerts.Columns().CreatedAt:   now,
		dao.PriceAlerts.Columns().UpdatedAt:   now,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建价格预警失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取预警ID失败")
	}

	g.Log().Info(ctx, "创建价格预警成功",
		"alertId", id,
		"symbol", symbol,
		"condition", req.Condition,
		"targetPrice", req.TargetPrice,
	)

	return s.getAlertByID(ctx, int(id))
}

// GetAlerts 获取预警列表
func (s *sPriceAlert) GetAlerts(ctx context.Context, userID int) ([]priceAlertApi.AlertItem, error) {
	var alerts []entity.PriceAlerts
	err := dao.PriceAlerts.Ctx(ctx).
		Where(dao.PriceAlerts.Columns().UserId, userID).
		WhereNull(dao.PriceAlerts.Columns().DeletedAt).
		OrderDesc(dao.PriceAlerts.Columns().CreatedAt).
		Scan(&alerts)
	if err != nil {
		return nil, gerror.Wrap(err, "查询预警列表失败")
	}

	items := make([]priceAlertApi.AlertItem, 0, len(alerts))
	for _, a := range alerts {
		items = append(items, s.toAlertItem(&a))
	}
	return items, nil
}

// UpdateAlert 更新预警（暂停/恢复）
func (s *sPriceAlert) UpdateAlert(ctx context.Context, req *priceAlertApi.UpdateReq) (*priceAlertApi.AlertItem, error) {
	// 检查是否存在
	var existing entity.PriceAlerts
	err := dao.PriceAlerts.Ctx(ctx).
		Where(dao.PriceAlerts.Columns().Id, req.Id).
		WhereNull(dao.PriceAlerts.Columns().DeletedAt).
		Scan(&existing)
	if err != nil {
		return nil, gerror.Wrap(err, "查询预警信息失败")
	}
	if existing.Id == 0 {
		return nil, gerror.Newf("预警不存在: %d", req.Id)
	}

	// 构建更新数据
	updateData := g.Map{
		dao.PriceAlerts.Columns().UpdatedAt: gtime.Now(),
	}
	if req.Symbol != "" {
		updateData[dao.PriceAlerts.Columns().Symbol] = strings.ToUpper(req.Symbol)
	}
	if req.Condition != "" {
		if !alertModel.IsValidCondition(req.Condition) {
			return nil, gerror.Newf("无效的预警条件: %s", req.Condition)
		}
		updateData[dao.PriceAlerts.Columns().Condition] = req.Condition
	}
	if req.TargetPrice > 0 {
		updateData[dao.PriceAlerts.Columns().TargetPrice] = req.TargetPrice
	}
	if req.Note != "" {
		updateData[dao.PriceAlerts.Columns().Note] = req.Note
	}
	if req.IsActive != nil {
		if *req.IsActive {
			updateData[dao.PriceAlerts.Columns().IsActive] = 1
		} else {
			updateData[dao.PriceAlerts.Columns().IsActive] = 0
		}
	}

	_, err = dao.PriceAlerts.Ctx(ctx).
		Where(dao.PriceAlerts.Columns().Id, req.Id).
		Data(updateData).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新预警失败")
	}

	g.Log().Info(ctx, "更新预警成功", "alertId", req.Id)

	return s.getAlertByID(ctx, int(req.Id))
}

// DeleteAlert 删除预警（软删除）
func (s *sPriceAlert) DeleteAlert(ctx context.Context, alertID int) error {
	_, err := dao.PriceAlerts.Ctx(ctx).
		Where(dao.PriceAlerts.Columns().Id, alertID).
		Data(g.Map{
			dao.PriceAlerts.Columns().DeletedAt: gtime.Now(),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除预警失败")
	}

	g.Log().Info(ctx, "删除预警成功", "alertId", alertID)
	return nil
}

// CheckAlerts 检查所有活跃预警，满足条件时触发通知
// 由 CronJob 定期调用
func (s *sPriceAlert) CheckAlerts(ctx context.Context) error {
	// 获取所有活跃且未触发的预警
	var alerts []entity.PriceAlerts
	err := dao.PriceAlerts.Ctx(ctx).
		Where(dao.PriceAlerts.Columns().IsActive, 1).
		Where(dao.PriceAlerts.Columns().Triggered, 0).
		WhereNull(dao.PriceAlerts.Columns().DeletedAt).
		Scan(&alerts)
	if err != nil {
		return gerror.Wrap(err, "获取活跃预警失败")
	}

	if len(alerts) == 0 {
		return nil
	}

	// 收集所有预警币种（去重）
	symbolSet := make(map[string]bool)
	for _, a := range alerts {
		symbolSet[strings.ToUpper(a.Symbol)] = true
	}
	symbols := make([]string, 0, len(symbolSet))
	for s := range symbolSet {
		symbols = append(symbols, s)
	}

	// 通过 CoinGecko 批量获取当前价格
	cgClient := coingecko.NewClient("")
	prices, err := cgClient.GetPrices(ctx, symbols)
	if err != nil {
		return gerror.Wrap(err, "获取当前价格失败")
	}

	// 逐条检查预警条件
	for _, alert := range alerts {
		currentPrice, ok := prices[strings.ToUpper(alert.Symbol)]
		if !ok {
			continue
		}

		triggered := false
		switch alert.Condition {
		case alertModel.ConditionAbove:
			triggered = currentPrice >= float64(alert.TargetPrice)
		case alertModel.ConditionBelow:
			triggered = currentPrice <= float64(alert.TargetPrice)
		}

		if triggered {
			// 标记已触发
			_, err := dao.PriceAlerts.Ctx(ctx).
				Where(dao.PriceAlerts.Columns().Id, alert.Id).
				Data(g.Map{
					dao.PriceAlerts.Columns().Triggered:   1,
					dao.PriceAlerts.Columns().TriggeredAt: gtime.Now(),
					dao.PriceAlerts.Columns().IsActive:    0,
				}).Update()
			if err != nil {
				g.Log().Warning(ctx, "更新预警状态失败", "alertId", alert.Id, "error", err)
				continue
			}

			// 发送通知
			condText := "高于"
			if alert.Condition == alertModel.ConditionBelow {
				condText = "低于"
			}
			title := fmt.Sprintf("%s 价格预警触发", alert.Symbol)
			message := fmt.Sprintf("%s 当前价格 $%.2f %s目标价 $%.2f",
				alert.Symbol, currentPrice, condText, float64(alert.TargetPrice))

			if sendErr := notificationService.Notification().Send(ctx, alert.UserId, "price_alert", title, message); sendErr != nil {
				g.Log().Warning(ctx, "发送预警通知失败", "alertId", alert.Id, "error", sendErr)
			}

			g.Log().Info(ctx, "预警触发", "alertId", alert.Id, "symbol", alert.Symbol,
				"condition", alert.Condition, "targetPrice", alert.TargetPrice,
				"currentPrice", currentPrice)
		}
	}

	return nil
}

// getAlertByID 根据 ID 查询预警
func (s *sPriceAlert) getAlertByID(ctx context.Context, alertID int) (*priceAlertApi.AlertItem, error) {
	var alert entity.PriceAlerts
	err := dao.PriceAlerts.Ctx(ctx).
		Where(dao.PriceAlerts.Columns().Id, alertID).
		WhereNull(dao.PriceAlerts.Columns().DeletedAt).
		Scan(&alert)
	if err != nil {
		return nil, gerror.Wrap(err, "查询预警信息失败")
	}
	if alert.Id == 0 {
		return nil, gerror.Newf("预警不存在: %d", alertID)
	}

	item := s.toAlertItem(&alert)
	return &item, nil
}

// toAlertItem 将数据库实体转换为 API 响应格式
func (s *sPriceAlert) toAlertItem(a *entity.PriceAlerts) priceAlertApi.AlertItem {
	item := priceAlertApi.AlertItem{
		ID:          uint(a.Id),
		Symbol:      a.Symbol,
		Condition:   a.Condition,
		TargetPrice: float64(a.TargetPrice),
		Note:        a.Note,
		IsActive:    a.IsActive == 1,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// 触发时间（如有）
	if !a.TriggeredAt.IsZero() {
		item.TriggeredAt = a.TriggeredAt.Format("2006-01-02T15:04:05Z")
	}

	return item
}
