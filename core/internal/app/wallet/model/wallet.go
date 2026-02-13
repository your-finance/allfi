// Package model 钱包模块 - 业务数据传输对象
// 定义钱包模块内部使用的 DTO
package model

// CreateWalletInput 创建钱包的内部输入参数
type CreateWalletInput struct {
	UserID     int    // 用户 ID
	Blockchain string // 区块链网络
	Address    string // 钱包地址
	Label      string // 地址标签
}

// BatchImportInput 批量导入的内部输入参数
type BatchImportInput struct {
	UserID    int              // 用户 ID
	Addresses []BatchAddressItem // 地址列表
}

// BatchAddressItem 批量导入的单个地址
type BatchAddressItem struct {
	Blockchain string // 区块链网络
	Address    string // 钱包地址
	Label      string // 地址标签
}
