// Package logic 报告业务逻辑
// 报告生成：查快照 -> 计算变化 -> 保存
package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	reportApi "your-finance/allfi/api/v1/report"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/report/dao"
	reportModel "your-finance/allfi/internal/app/report/model"
	"your-finance/allfi/internal/app/report/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/model/entity"
)

// sReport 报告服务实现
type sReport struct{}

// New 创建报告服务实例
func New() service.IReport {
	return &sReport{}
}

// GetReports 获取报告列表
func (s *sReport) GetReports(ctx context.Context, reportType string, limit int) (*reportApi.ListRes, error) {
	if limit <= 0 {
		limit = 20
	}

	// 构建查询
	query := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, consts.GetUserID(ctx))

	// 按类型筛选
	if reportType != "" {
		query = query.Where(dao.Reports.Columns().Type, reportType)
	}

	// 查询
	var reports []entity.Reports
	err := query.
		OrderDesc(dao.Reports.Columns().GeneratedAt).
		Limit(limit).
		Scan(&reports)
	if err != nil {
		return nil, gerror.Wrap(err, "查询报告列表失败")
	}

	// 转换为 API 响应格式
	items := make([]reportApi.ReportSummary, 0, len(reports))
	for _, r := range reports {
		items = append(items, toReportSummary(&r))
	}

	return &reportApi.ListRes{Reports: items}, nil
}

// GetReport 获取单个报告详情
func (s *sReport) GetReport(ctx context.Context, id uint) (*reportApi.GetRes, error) {
	var report entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().Id, id).
		Scan(&report)
	if err != nil {
		return nil, gerror.Wrap(err, "查询报告失败")
	}
	if report.Id == 0 {
		return nil, gerror.Newf("报告不存在: %d", id)
	}

	// 解析报告内容
	var content interface{}
	if report.Content != "" {
		var parsed map[string]interface{}
		if json.Unmarshal([]byte(report.Content), &parsed) == nil {
			content = parsed
		} else {
			content = report.Content
		}
	} else {
		// 构建默认内容
		content = map[string]interface{}{
			"id":             report.Id,
			"type":           report.Type,
			"period":         report.Period,
			"total_value":    report.TotalValue,
			"change":         report.Change,
			"change_percent": report.ChangePercent,
			"btc_benchmark":  report.BtcBenchmark,
			"eth_benchmark":  report.EthBenchmark,
			"top_gainers":    report.TopGainers,
			"top_losers":     report.TopLosers,
			"generated_at":   report.GeneratedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &reportApi.GetRes{Report: content}, nil
}

// GenerateReport 生成报告
func (s *sReport) GenerateReport(ctx context.Context, reportType string) (*reportApi.GenerateRes, error) {
	userID := consts.GetUserID(ctx)

	var report *entity.Reports
	var err error

	switch reportType {
	case reportModel.ReportTypeDaily:
		report, err = s.generateDailyReport(ctx, userID)
	case reportModel.ReportTypeWeekly:
		report, err = s.generateWeeklyReport(ctx, userID)
	case reportModel.ReportTypeMonthly:
		month := time.Now().Format("2006-01")
		report, err = s.generateMonthlyReport(ctx, userID, month)
	case reportModel.ReportTypeAnnual:
		year := time.Now().Format("2006")
		report, err = s.generateAnnualReport(ctx, userID, year)
	default:
		return nil, gerror.Newf("不支持的报告类型: %s", reportType)
	}

	if err != nil {
		return nil, err
	}

	summary := toReportSummary(report)
	return &reportApi.GenerateRes{Report: &summary}, nil
}

// GetMonthlyReport 获取/生成月度报告
func (s *sReport) GetMonthlyReport(ctx context.Context, month string) (*reportApi.GetMonthlyRes, error) {
	userID := consts.GetUserID(ctx)

	// 先查找已有报告
	var existing entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, userID).
		Where(dao.Reports.Columns().Type, reportModel.ReportTypeMonthly).
		Where(dao.Reports.Columns().Period, month).
		Scan(&existing)
	if err == nil && existing.Id > 0 {
		// 解析内容
		var content interface{}
		if existing.Content != "" {
			var parsed map[string]interface{}
			if json.Unmarshal([]byte(existing.Content), &parsed) == nil {
				content = parsed
			}
		}
		if content == nil {
			content = map[string]interface{}{
				"id":             existing.Id,
				"type":           existing.Type,
				"period":         existing.Period,
				"total_value":    existing.TotalValue,
				"change":         existing.Change,
				"change_percent": existing.ChangePercent,
			}
		}
		return &reportApi.GetMonthlyRes{Report: content}, nil
	}

	// 生成月报
	report, err := s.generateMonthlyReport(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	var content interface{}
	if report.Content != "" {
		var parsed map[string]interface{}
		if json.Unmarshal([]byte(report.Content), &parsed) == nil {
			content = parsed
		}
	}
	if content == nil {
		content = toReportSummary(report)
	}

	return &reportApi.GetMonthlyRes{Report: content}, nil
}

// GetAnnualReport 获取/生成年度报告
func (s *sReport) GetAnnualReport(ctx context.Context, year string) (*reportApi.GetAnnualRes, error) {
	userID := consts.GetUserID(ctx)

	// 先查找已有报告
	var existing entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, userID).
		Where(dao.Reports.Columns().Type, reportModel.ReportTypeAnnual).
		Where(dao.Reports.Columns().Period, year).
		Scan(&existing)
	if err == nil && existing.Id > 0 {
		var content interface{}
		if existing.Content != "" {
			var parsed map[string]interface{}
			if json.Unmarshal([]byte(existing.Content), &parsed) == nil {
				content = parsed
			}
		}
		if content == nil {
			content = map[string]interface{}{
				"id":             existing.Id,
				"type":           existing.Type,
				"period":         existing.Period,
				"total_value":    existing.TotalValue,
				"change":         existing.Change,
				"change_percent": existing.ChangePercent,
			}
		}
		return &reportApi.GetAnnualRes{Report: content}, nil
	}

	// 生成年报
	report, err := s.generateAnnualReport(ctx, userID, year)
	if err != nil {
		return nil, err
	}

	var content interface{}
	if report.Content != "" {
		var parsed map[string]interface{}
		if json.Unmarshal([]byte(report.Content), &parsed) == nil {
			content = parsed
		}
	}
	if content == nil {
		content = toReportSummary(report)
	}

	return &reportApi.GetAnnualRes{Report: content}, nil
}

// generateDailyReport 生成日报
func (s *sReport) generateDailyReport(ctx context.Context, userID int) (*entity.Reports, error) {
	now := time.Now()
	period := now.Format("2006-01-02")

	// 检查是否已存在
	var existing entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, userID).
		Where(dao.Reports.Columns().Type, reportModel.ReportTypeDaily).
		Where(dao.Reports.Columns().Period, period).
		Scan(&existing)
	if err == nil && existing.Id > 0 {
		return &existing, nil
	}

	// 获取今日快照
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var todaySnapshots []entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, todayStart).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&todaySnapshots)
	if err != nil || len(todaySnapshots) == 0 {
		return nil, gerror.New("未找到今日快照数据")
	}
	latest := todaySnapshots[0]

	// 获取昨日快照
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	var yesterdaySnapshots []entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, yesterdayStart).
		WhereLT(assetDao.AssetSnapshots.Columns().SnapshotTime, todayStart).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&yesterdaySnapshots)

	var change, changePercent float32
	if len(yesterdaySnapshots) > 0 {
		yesterdayValue := float32(yesterdaySnapshots[0].TotalValueUsd)
		todayValue := float32(latest.TotalValueUsd)
		if yesterdayValue > 0 {
			change = todayValue - yesterdayValue
			changePercent = (change / yesterdayValue) * 100
		}
	}

	report := &entity.Reports{
		UserId:        userID,
		Type:          reportModel.ReportTypeDaily,
		Period:        period,
		TotalValue:    float32(latest.TotalValueUsd),
		Change:        change,
		ChangePercent: changePercent,
		GeneratedAt:   gtime.Now().Time,
	}

	result, err := dao.Reports.Ctx(ctx).Insert(report)
	if err != nil {
		return nil, gerror.Wrap(err, "保存日报失败")
	}
	id, _ := result.LastInsertId()
	report.Id = int(id)

	g.Log().Info(ctx, "日报已生成", "period", period, "totalValue", latest.TotalValueUsd)
	return report, nil
}

// generateWeeklyReport 生成周报
func (s *sReport) generateWeeklyReport(ctx context.Context, userID int) (*entity.Reports, error) {
	now := time.Now()
	year, week := now.ISOWeek()
	period := fmt.Sprintf("%d-W%02d", year, week)

	// 检查是否已存在
	var existing entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, userID).
		Where(dao.Reports.Columns().Type, reportModel.ReportTypeWeekly).
		Where(dao.Reports.Columns().Period, period).
		Scan(&existing)
	if err == nil && existing.Id > 0 {
		return &existing, nil
	}

	// 获取本周一
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	thisMonday := time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, now.Location())
	lastMonday := thisMonday.AddDate(0, 0, -7)

	// 获取本周快照
	var thisWeekSnapshots []entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, thisMonday).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&thisWeekSnapshots)
	if err != nil || len(thisWeekSnapshots) == 0 {
		return nil, gerror.New("未找到本周快照数据")
	}
	latest := thisWeekSnapshots[0]

	// 获取上周快照
	var lastWeekSnapshots []entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, lastMonday).
		WhereLT(assetDao.AssetSnapshots.Columns().SnapshotTime, thisMonday).
		OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&lastWeekSnapshots)

	var change, changePercent float32
	if len(lastWeekSnapshots) > 0 {
		lastValue := float32(lastWeekSnapshots[0].TotalValueUsd)
		thisValue := float32(latest.TotalValueUsd)
		if lastValue > 0 {
			change = thisValue - lastValue
			changePercent = (change / lastValue) * 100
		}
	}

	report := &entity.Reports{
		UserId:        userID,
		Type:          reportModel.ReportTypeWeekly,
		Period:        period,
		TotalValue:    float32(latest.TotalValueUsd),
		Change:        change,
		ChangePercent: changePercent,
		GeneratedAt:   gtime.Now().Time,
	}

	result, err := dao.Reports.Ctx(ctx).Insert(report)
	if err != nil {
		return nil, gerror.Wrap(err, "保存周报失败")
	}
	id, _ := result.LastInsertId()
	report.Id = int(id)

	g.Log().Info(ctx, "周报已生成", "period", period)
	return report, nil
}

// generateMonthlyReport 生成月报
func (s *sReport) generateMonthlyReport(ctx context.Context, userID int, month string) (*entity.Reports, error) {
	// 检查是否已存在
	var existing entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, userID).
		Where(dao.Reports.Columns().Type, reportModel.ReportTypeMonthly).
		Where(dao.Reports.Columns().Period, month).
		Scan(&existing)
	if err == nil && existing.Id > 0 {
		return &existing, nil
	}

	// 解析月份
	monthStart, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, gerror.Newf("月份格式错误: %s", month)
	}
	monthEnd := monthStart.AddDate(0, 1, 0)
	prevMonthStart := monthStart.AddDate(0, -1, 0)

	// 获取本月快照
	var thisMonthSnapshots []entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, monthStart).
		WhereLT(assetDao.AssetSnapshots.Columns().SnapshotTime, monthEnd).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&thisMonthSnapshots)
	if err != nil || len(thisMonthSnapshots) == 0 {
		return nil, gerror.Newf("未找到 %s 月的快照数据", month)
	}
	latest := thisMonthSnapshots[0]

	// 获取上月快照
	var prevMonthSnapshots []entity.AssetSnapshots
	_ = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, prevMonthStart).
		WhereLT(assetDao.AssetSnapshots.Columns().SnapshotTime, monthStart).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&prevMonthSnapshots)

	var change, changePercent float32
	if len(prevMonthSnapshots) > 0 {
		prevValue := float32(prevMonthSnapshots[0].TotalValueUsd)
		thisValue := float32(latest.TotalValueUsd)
		if prevValue > 0 {
			change = thisValue - prevValue
			changePercent = (change / prevValue) * 100
		}
	}

	// 计算 TopGainers/TopLosers
	gainers, losers := s.getTopGainersLosers(ctx, userID, 5)
	gainersJSON, _ := json.Marshal(gainers)
	losersJSON, _ := json.Marshal(losers)

	// 构建月报内容
	content := reportModel.ReportContent{
		Type:          reportModel.ReportTypeMonthly,
		TotalValue:    float64(float32(latest.TotalValueUsd)),
		Change:        float64(change),
		ChangePercent: float64(changePercent),
		SnapshotCount: len(thisMonthSnapshots),
		TopGainers:    gainers,
		TopLosers:     losers,
		Summary:       fmt.Sprintf("本月资产总值 $%.2f，较上月变化 %.2f%%", latest.TotalValueUsd, changePercent),
	}
	contentJSON, _ := json.Marshal(content)

	report := &entity.Reports{
		UserId:        userID,
		Type:          reportModel.ReportTypeMonthly,
		Period:        month,
		TotalValue:    float32(latest.TotalValueUsd),
		Change:        change,
		ChangePercent: changePercent,
		TopGainers:    string(gainersJSON),
		TopLosers:     string(losersJSON),
		Content:       string(contentJSON),
		GeneratedAt:   gtime.Now().Time,
	}

	result, err := dao.Reports.Ctx(ctx).Insert(report)
	if err != nil {
		return nil, gerror.Wrap(err, "保存月报失败")
	}
	id, _ := result.LastInsertId()
	report.Id = int(id)

	g.Log().Info(ctx, "月报已生成", "period", month)
	return report, nil
}

// generateAnnualReport 生成年报
func (s *sReport) generateAnnualReport(ctx context.Context, userID int, year string) (*entity.Reports, error) {
	// 检查是否已存在
	var existing entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().UserId, userID).
		Where(dao.Reports.Columns().Type, reportModel.ReportTypeAnnual).
		Where(dao.Reports.Columns().Period, year).
		Scan(&existing)
	if err == nil && existing.Id > 0 {
		return &existing, nil
	}

	// 解析年份
	yearStart, err := time.Parse("2006", year)
	if err != nil {
		return nil, gerror.Newf("年份格式错误: %s", year)
	}
	yearEnd := yearStart.AddDate(1, 0, 0)
	prevYearStart := yearStart.AddDate(-1, 0, 0)

	// 获取本年快照
	var thisYearSnapshots []entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, yearStart).
		WhereLT(assetDao.AssetSnapshots.Columns().SnapshotTime, yearEnd).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&thisYearSnapshots)
	if err != nil || len(thisYearSnapshots) == 0 {
		return nil, gerror.Newf("未找到 %s 年的快照数据", year)
	}
	latest := thisYearSnapshots[0]

	// 获取上年快照
	var prevYearSnapshots []entity.AssetSnapshots
	_ = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, prevYearStart).
		WhereLT(assetDao.AssetSnapshots.Columns().SnapshotTime, yearStart).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&prevYearSnapshots)

	var change, changePercent float32
	if len(prevYearSnapshots) > 0 {
		prevValue := float32(prevYearSnapshots[0].TotalValueUsd)
		thisValue := float32(latest.TotalValueUsd)
		if prevValue > 0 {
			change = thisValue - prevValue
			changePercent = (change / prevValue) * 100
		}
	}

	// 计算 TopGainers/TopLosers
	gainers, losers := s.getTopGainersLosers(ctx, userID, 5)
	gainersJSON, _ := json.Marshal(gainers)
	losersJSON, _ := json.Marshal(losers)

	// 构建年报内容
	content := reportModel.ReportContent{
		Type:          reportModel.ReportTypeAnnual,
		TotalValue:    float64(float32(latest.TotalValueUsd)),
		Change:        float64(change),
		ChangePercent: float64(changePercent),
		SnapshotCount: len(thisYearSnapshots),
		TopGainers:    gainers,
		TopLosers:     losers,
		Summary:       fmt.Sprintf("年度资产总值 $%.2f，年度收益率 %.2f%%", latest.TotalValueUsd, changePercent),
	}
	contentJSON, _ := json.Marshal(content)

	report := &entity.Reports{
		UserId:        userID,
		Type:          reportModel.ReportTypeAnnual,
		Period:        year,
		TotalValue:    float32(latest.TotalValueUsd),
		Change:        change,
		ChangePercent: changePercent,
		TopGainers:    string(gainersJSON),
		TopLosers:     string(losersJSON),
		Content:       string(contentJSON),
		GeneratedAt:   gtime.Now().Time,
	}

	result, err := dao.Reports.Ctx(ctx).Insert(report)
	if err != nil {
		return nil, gerror.Wrap(err, "保存年报失败")
	}
	id, _ := result.LastInsertId()
	report.Id = int(id)

	g.Log().Info(ctx, "年报已生成", "period", year)
	return report, nil
}

// getTopGainersLosers 获取涨跌幅 Top 资产
// 按当前持仓价值排序，前 N 为 TopGainers，后 N 为 TopLosers
func (s *sReport) getTopGainersLosers(ctx context.Context, userID int, topN int) (gainers, losers []reportModel.GainerLoser) {
	// 获取当前资产明细
	var details []entity.AssetDetails
	err := assetDao.AssetDetails.Ctx(ctx).
		Where(assetDao.AssetDetails.Columns().UserId, userID).
		Scan(&details)
	if err != nil || len(details) == 0 {
		return nil, nil
	}

	// 按币种聚合
	type assetAgg struct {
		Symbol   string
		TotalVal float64
	}
	aggMap := make(map[string]*assetAgg)
	for _, d := range details {
		if d.ValueUsd <= 0 {
			continue
		}
		if agg, ok := aggMap[d.AssetSymbol]; ok {
			agg.TotalVal += d.ValueUsd
		} else {
			aggMap[d.AssetSymbol] = &assetAgg{Symbol: d.AssetSymbol, TotalVal: d.ValueUsd}
		}
	}

	// 转为切片并排序
	assets := make([]assetAgg, 0, len(aggMap))
	for _, a := range aggMap {
		assets = append(assets, *a)
	}
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].TotalVal > assets[j].TotalVal
	})

	// 前 topN 为 TopGainers
	for i := 0; i < topN && i < len(assets); i++ {
		gainers = append(gainers, reportModel.GainerLoser{
			Symbol: assets[i].Symbol,
			Value:  assets[i].TotalVal,
		})
	}

	// 后 topN 为 TopLosers
	for i := len(assets) - 1; i >= 0 && len(losers) < topN; i-- {
		if i < topN {
			break // 避免与 gainers 重复
		}
		losers = append(losers, reportModel.GainerLoser{
			Symbol: assets[i].Symbol,
			Value:  assets[i].TotalVal,
		})
	}

	return gainers, losers
}

// Compare 对比两份报告
// 查询两份报告，提取关键数据，计算差异
func (s *sReport) Compare(ctx context.Context, reportID1, reportID2 int) (*reportApi.CompareRes, error) {
	// 查询报告 1
	var report1 entity.Reports
	err := dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().Id, reportID1).
		Scan(&report1)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询报告 %d 失败", reportID1)
	}
	if report1.Id == 0 {
		return nil, gerror.Newf("报告不存在: %d", reportID1)
	}

	// 查询报告 2
	var report2 entity.Reports
	err = dao.Reports.Ctx(ctx).
		Where(dao.Reports.Columns().Id, reportID2).
		Scan(&report2)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询报告 %d 失败", reportID2)
	}
	if report2.Id == 0 {
		return nil, gerror.Newf("报告不存在: %d", reportID2)
	}

	// 提取报告摘要（优先从 JSON content 中解析，回退到表字段）
	summary1 := extractCompareSummary(&report1)
	summary2 := extractCompareSummary(&report2)

	// 计算差异
	valueDiff := summary2.TotalValue - summary1.TotalValue
	changeDiff := summary2.Change - summary1.Change

	return &reportApi.CompareRes{
		Report1:    summary1,
		Report2:    summary2,
		ValueDiff:  valueDiff,
		ChangeDiff: changeDiff,
	}, nil
}

// extractCompareSummary 从报告实体中提取对比摘要
// 优先解析 JSON content 字段，回退使用表的直接字段
func extractCompareSummary(r *entity.Reports) *reportApi.CompareReportSummary {
	summary := &reportApi.CompareReportSummary{
		ID:          uint(r.Id),
		Type:        r.Type,
		Period:      r.Period,
		TotalValue:  float64(r.TotalValue),
		Change:      float64(r.Change),
		ChangePct:   float64(r.ChangePercent),
		GeneratedAt: r.GeneratedAt.Format("2006-01-02 15:04:05"),
	}

	// 尝试从 JSON content 中提取更精确的数据
	if r.Content != "" {
		var content reportModel.ReportContent
		if json.Unmarshal([]byte(r.Content), &content) == nil {
			if content.TotalValue > 0 {
				summary.TotalValue = content.TotalValue
			}
			if content.Change != 0 {
				summary.Change = content.Change
			}
			if content.ChangePercent != 0 {
				summary.ChangePct = content.ChangePercent
			}
		}
	}

	return summary
}

// toReportSummary 将实体转换为 API 摘要
func toReportSummary(r *entity.Reports) reportApi.ReportSummary {
	title := ""
	switch r.Type {
	case reportModel.ReportTypeDaily:
		title = fmt.Sprintf("每日报告 - %s", r.Period)
	case reportModel.ReportTypeWeekly:
		title = fmt.Sprintf("每周报告 - %s", r.Period)
	case reportModel.ReportTypeMonthly:
		title = fmt.Sprintf("月度报告 - %s", r.Period)
	case reportModel.ReportTypeAnnual:
		title = fmt.Sprintf("年度报告 - %s", r.Period)
	default:
		title = r.Period
	}

	return reportApi.ReportSummary{
		ID:         uint(r.Id),
		Type:       r.Type,
		Title:      title,
		TotalValue: float64(r.TotalValue),
		PnL:        float64(r.Change),
		PnLPercent: float64(r.ChangePercent),
		CreatedAt:  r.GeneratedAt.Format("2006-01-02 15:04:05"),
	}
}
