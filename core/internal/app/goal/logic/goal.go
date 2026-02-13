// Package logic 目标追踪业务逻辑
// 实现投资目标的增删改查和进度计算
package logic

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	goalApi "your-finance/allfi/api/v1/goal"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	assetService "your-finance/allfi/internal/app/asset/service"
	"your-finance/allfi/internal/app/goal/dao"
	goalModel "your-finance/allfi/internal/app/goal/model"
	"your-finance/allfi/internal/app/goal/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/model/entity"
)

// sGoal 目标追踪服务实现
type sGoal struct{}

// New 创建目标追踪服务实例
func New() service.IGoal {
	return &sGoal{}
}

// GetGoals 获取目标列表（带进度百分比）
// 通过资产服务获取当前值，计算每个目标的进度
func (s *sGoal) GetGoals(ctx context.Context) ([]goalApi.GoalItem, error) {
	var goals []entity.Goals
	err := dao.Goals.Ctx(ctx).
		WhereNull(dao.Goals.Columns().DeletedAt).
		OrderDesc(dao.Goals.Columns().CreatedAt).
		Scan(&goals)
	if err != nil {
		return nil, gerror.Wrap(err, "查询目标列表失败")
	}

	items := make([]goalApi.GoalItem, 0, len(goals))
	for _, goal := range goals {
		item := s.toGoalItem(&goal)

		// 计算当前进度
		currentValue := s.calculateCurrentValue(ctx, &goal)
		item.CurrentValue = math.Round(currentValue*100) / 100

		if float64(goal.TargetValue) > 0 {
			progress := (currentValue / float64(goal.TargetValue)) * 100
			if progress > 100 {
				progress = 100
			}
			item.Progress = math.Round(progress*100) / 100
		}

		// 判断是否已完成
		item.IsCompleted = item.Progress >= 100

		items = append(items, item)
	}

	return items, nil
}

// CreateGoal 创建目标
func (s *sGoal) CreateGoal(ctx context.Context, req *goalApi.CreateReq) (*goalApi.GoalItem, error) {
	// 验证目标类型
	if !goalModel.IsValidGoalType(req.Type) {
		return nil, gerror.Newf("不支持的目标类型: %s", req.Type)
	}

	// 目标值必须大于 0
	if req.TargetValue <= 0 {
		return nil, gerror.New("目标值必须大于 0")
	}

	// 默认货币
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	// 构建插入数据
	now := gtime.Now()
	data := g.Map{
		dao.Goals.Columns().Title:       req.Title,
		dao.Goals.Columns().Type:        req.Type,
		dao.Goals.Columns().TargetValue: req.TargetValue,
		dao.Goals.Columns().Currency:    currency,
		dao.Goals.Columns().CreatedAt:   now,
		dao.Goals.Columns().UpdatedAt:   now,
	}

	// 设置截止日期（可选）
	if req.Deadline != "" {
		deadline := gtime.NewFromStr(req.Deadline)
		if deadline == nil {
			return nil, gerror.New("截止日期格式错误，请使用 ISO 8601 格式")
		}
		data[dao.Goals.Columns().Deadline] = deadline
	}

	result, err := dao.Goals.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建目标失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取目标ID失败")
	}

	g.Log().Info(ctx, "创建目标成功",
		"goalId", id,
		"title", req.Title,
		"type", req.Type,
		"targetValue", req.TargetValue,
	)

	return s.getGoalByID(ctx, int(id))
}

// UpdateGoal 更新目标
func (s *sGoal) UpdateGoal(ctx context.Context, req *goalApi.UpdateReq) (*goalApi.GoalItem, error) {
	// 检查是否存在
	var existing entity.Goals
	err := dao.Goals.Ctx(ctx).
		Where(dao.Goals.Columns().Id, req.Id).
		WhereNull(dao.Goals.Columns().DeletedAt).
		Scan(&existing)
	if err != nil {
		return nil, gerror.Wrap(err, "查询目标信息失败")
	}
	if existing.Id == 0 {
		return nil, gerror.Newf("目标不存在: %d", req.Id)
	}

	// 验证目标类型
	if req.Type != "" && !goalModel.IsValidGoalType(req.Type) {
		return nil, gerror.Newf("不支持的目标类型: %s", req.Type)
	}

	// 构建更新数据
	updateData := g.Map{
		dao.Goals.Columns().UpdatedAt: gtime.Now(),
	}
	if req.Title != "" {
		updateData[dao.Goals.Columns().Title] = req.Title
	}
	if req.Type != "" {
		updateData[dao.Goals.Columns().Type] = req.Type
	}
	if req.TargetValue > 0 {
		updateData[dao.Goals.Columns().TargetValue] = req.TargetValue
	}
	if req.Currency != "" {
		updateData[dao.Goals.Columns().Currency] = req.Currency
	}
	if req.Deadline != "" {
		deadline := gtime.NewFromStr(req.Deadline)
		if deadline == nil {
			return nil, gerror.New("截止日期格式错误")
		}
		updateData[dao.Goals.Columns().Deadline] = deadline
	}

	_, err = dao.Goals.Ctx(ctx).
		Where(dao.Goals.Columns().Id, req.Id).
		Data(updateData).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新目标失败")
	}

	g.Log().Info(ctx, "更新目标成功", "goalId", req.Id)

	return s.getGoalByID(ctx, int(req.Id))
}

// DeleteGoal 删除目标（软删除）
func (s *sGoal) DeleteGoal(ctx context.Context, goalID int) error {
	_, err := dao.Goals.Ctx(ctx).
		Where(dao.Goals.Columns().Id, goalID).
		Data(g.Map{
			dao.Goals.Columns().DeletedAt: gtime.Now(),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除目标失败")
	}

	g.Log().Info(ctx, "删除目标成功", "goalId", goalID)
	return nil
}

// calculateCurrentValue 根据目标类型计算当前值
// 通过资产服务和 DAO 查询获取实时数据
func (s *sGoal) calculateCurrentValue(ctx context.Context, goal *entity.Goals) float64 {
	switch goal.Type {
	case goalModel.GoalTypeAssetValue:
		// 调用资产服务获取当前总资产值（按目标指定的货币计价）
		res, err := assetService.Asset().GetSummary(ctx, goal.Currency)
		if err != nil {
			g.Log().Warning(ctx, "获取资产概览失败，目标进度返回0", "goalId", goal.Id, "error", err)
			return 0
		}
		return res.TotalValue

	case goalModel.GoalTypeReturnRate:
		// 从快照数据计算收益率：(最新 - 最旧) / 最旧 * 100
		var snapshots []entity.AssetSnapshots
		err := assetDao.AssetSnapshots.Ctx(ctx).
			Where(assetDao.AssetSnapshots.Columns().UserId, consts.GetUserID(ctx)).
			OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
			Scan(&snapshots)
		if err != nil || len(snapshots) < 2 {
			g.Log().Warning(ctx, "快照数据不足，收益率返回0", "goalId", goal.Id, "error", err)
			return 0
		}
		oldest := snapshots[0]
		latest := snapshots[len(snapshots)-1]
		if oldest.TotalValueUsd <= 0 {
			return 0
		}
		return (latest.TotalValueUsd - oldest.TotalValueUsd) / oldest.TotalValueUsd * 100

	case goalModel.GoalTypeHoldingAmount:
		// 从资产明细表查询指定 symbol 的总持仓数量
		// goal.Currency 字段在此类型中存储的是 symbol（如 BTC、ETH）
		balance, err := assetDao.AssetDetails.Ctx(ctx).
			Where(assetDao.AssetDetails.Columns().UserId, consts.GetUserID(ctx)).
			Where(assetDao.AssetDetails.Columns().AssetSymbol, goal.Currency).
			Sum(assetDao.AssetDetails.Columns().Balance)
		if err != nil {
			g.Log().Warning(ctx, "查询持仓数量失败，返回0", "goalId", goal.Id, "symbol", goal.Currency, "error", err)
			return 0
		}
		return balance

	default:
		return 0
	}
}

// getGoalByID 根据 ID 查询目标
func (s *sGoal) getGoalByID(ctx context.Context, goalID int) (*goalApi.GoalItem, error) {
	var goal entity.Goals
	err := dao.Goals.Ctx(ctx).
		Where(dao.Goals.Columns().Id, goalID).
		WhereNull(dao.Goals.Columns().DeletedAt).
		Scan(&goal)
	if err != nil {
		return nil, gerror.Wrap(err, "查询目标信息失败")
	}
	if goal.Id == 0 {
		return nil, gerror.Newf("目标不存在: %d", goalID)
	}

	item := s.toGoalItem(&goal)
	return &item, nil
}

// toGoalItem 将数据库实体转换为 API 响应格式
func (s *sGoal) toGoalItem(g *entity.Goals) goalApi.GoalItem {
	item := goalApi.GoalItem{
		ID:          uint(g.Id),
		Title:       g.Title,
		Type:        g.Type,
		TargetValue: float64(g.TargetValue),
		Currency:    g.Currency,
		CreatedAt:   g.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   g.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// 设置截止日期（如有）
	if !g.Deadline.IsZero() {
		item.Deadline = g.Deadline.Format("2006-01-02T15:04:05Z")
	}

	return item
}
