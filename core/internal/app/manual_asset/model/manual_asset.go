// Package model 手动资产模块 - 业务数据传输对象
// 定义手动资产模块内部使用的 DTO
package model

// CreateManualAssetInput 创建手动资产的内部输入参数
type CreateManualAssetInput struct {
	UserID      int     // 用户 ID
	AssetType   string  // 资产类型（cash/bank/stock/fund）
	AssetName   string  // 资产名称
	Amount      float64 // 数量
	Currency    string  // 货币
	Notes       string  // 备注
	Institution string  // 机构名称
}

// UpdateManualAssetInput 更新手动资产的内部输入参数
type UpdateManualAssetInput struct {
	AssetID     int     // 资产 ID
	AssetType   string  // 资产类型
	AssetName   string  // 资产名称
	Amount      float64 // 数量
	Currency    string  // 货币
	Notes       string  // 备注
	Institution string  // 机构名称
}
