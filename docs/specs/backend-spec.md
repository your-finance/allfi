# AllFi 后端需求规格文档

> 版本：v2.0 | 更新时间：2026-02-11 | 状态：与代码实现对齐

---

## 文档概述

本文档描述 AllFi 后端的实际技术架构、API 接口、数据模型、服务层设计和第三方集成方案。所有内容均与当前代码实现对齐。

---

## 1. 技术栈

### 1.1 核心依赖

| 层 | 技术 | 说明 |
|----|------|------|
| 语言 | Go 1.24 | go.mod 固定 `go 1.24.11` |
| HTTP 路由 | `net/http` (标准库) | Go 1.22+ 增强路由（`METHOD /path/{param}`） |
| 配置管理 | GoFrame `gcfg` | YAML 配置 + 环境变量覆盖 |
| ORM | GORM v2 | 支持 SQLite3 / MySQL 双驱动 |
| 数据库 | SQLite3（默认） / MySQL（可选） | 自动迁移，无需手动建表 |
| 加密 | AES-256-GCM | API Key 等敏感数据加密存储 |
| 认证 | PIN + JWT | 单用户 PIN 锁 + Bearer Token |
| JSON | json-iterator | 高性能 JSON 序列化 |
| 精确计算 | shopspring/decimal | 金额使用 Decimal 类型 |
| 监控 | Prometheus | `/metrics` 端点 |
| 推送 | webpush-go | Web Push 通知 |
| 密码学 | golang.org/x/crypto (bcrypt) | PIN 码哈希 |

### 1.2 go.mod 直接依赖

```
github.com/adshao/go-binance/v2    # Binance SDK
github.com/ccxt/ccxt/go/v4          # CCXT 多交易所统一 SDK
github.com/gogf/gf/v2               # GoFrame（仅用于配置管理）
github.com/golang-jwt/jwt/v5        # JWT Token
github.com/json-iterator/go         # 高性能 JSON
github.com/prometheus/client_golang  # Prometheus 监控
github.com/shopspring/decimal        # 精确数值计算
github.com/stretchr/testify          # 测试断言
golang.org/x/crypto                  # bcrypt 密码哈希
gorm.io/gorm                         # ORM
gorm.io/driver/sqlite                # SQLite3 驱动
gorm.io/driver/mysql                 # MySQL 驱动
```

---

## 2. 项目结构

```
core/
├── cmd/
│   └── server/main.go               # 服务入口
├── config/
│   ├── config.yaml                   # 配置文件
│   └── config.example.yaml           # 配置模板
├── internal/
│   ├── api/
│   │   ├── router.go                 # 路由定义（75 条路由）
│   │   ├── handlers/                 # HTTP 处理器（24 个文件）
│   │   │   ├── common.go            # 通用工具函数
│   │   │   ├── auth.go              # 认证处理器
│   │   │   ├── asset.go             # 资产处理器
│   │   │   ├── exchange.go          # 交易所处理器
│   │   │   ├── wallet.go            # 钱包处理器
│   │   │   ├── rate.go              # 汇率处理器
│   │   │   ├── notification.go      # 通知处理器
│   │   │   ├── price_alert.go       # 价格预警处理器
│   │   │   ├── report.go            # 报告处理器
│   │   │   ├── defi.go              # DeFi 处理器
│   │   │   ├── nft.go               # NFT 处理器
│   │   │   ├── transaction.go       # 交易记录处理器
│   │   │   ├── webpush.go           # WebPush 处理器
│   │   │   ├── fee.go               # 费用分析处理器
│   │   │   ├── strategy.go          # 策略处理器
│   │   │   ├── achievement.go       # 成就处理器
│   │   │   ├── benchmark.go         # 基准对比处理器
│   │   │   ├── market.go            # 市场数据处理器
│   │   │   ├── goal.go              # 目标追踪处理器
│   │   │   ├── health.go            # 健康检查处理器
│   │   │   ├── health_score.go      # 健康评分处理器
│   │   │   ├── pnl.go               # 盈亏分析处理器
│   │   │   ├── attribution.go       # 归因分析处理器
│   │   │   └── forecast.go          # 趋势预测处理器
│   │   └── middleware/               # 中间件（4 个）
│   │       ├── cors.go              # CORS 跨域
│   │       ├── logger.go            # 请求日志
│   │       ├── recovery.go          # Panic 恢复
│   │       └── auth.go              # JWT 认证
│   ├── service/                      # 业务服务层（20+ 服务）
│   │   ├── interfaces.go            # 服务接口定义
│   │   ├── container.go             # 依赖注入容器
│   │   ├── auth_service.go          # 认证服务
│   │   ├── asset_service.go         # 资产聚合服务
│   │   ├── exchange_service.go      # 交易所服务
│   │   ├── blockchain_service.go    # 区块链服务
│   │   ├── price_service.go         # 价格服务
│   │   ├── snapshot_service.go      # 快照服务
│   │   ├── notification_service.go  # 通知服务
│   │   ├── price_alert_service.go   # 价格预警服务
│   │   ├── report_service.go        # 报告服务
│   │   ├── defi_service.go          # DeFi 服务
│   │   ├── nft_service.go           # NFT 服务
│   │   ├── transaction_service.go   # 交易记录服务
│   │   ├── fee_service.go           # 费用分析服务
│   │   ├── strategy_service.go      # 策略引擎服务
│   │   ├── achievement_service.go   # 成就系统服务
│   │   ├── benchmark_service.go     # 基准对比服务
│   │   ├── goal_service.go          # 目标追踪服务
│   │   ├── health_score_service.go  # 健康评分服务
│   │   ├── risk_alert_service.go    # 风险预警服务
│   │   └── webpush.go               # WebPush 推送服务
│   ├── repository/                   # 数据访问层（18 个）
│   │   ├── base_repository.go       # 基础仓库
│   │   ├── exchange_account_repo.go
│   │   ├── wallet_address_repo.go
│   │   ├── manual_asset_repo.go
│   │   ├── asset_detail_repo.go
│   │   ├── asset_snapshot_repo.go
│   │   ├── exchange_rate_repo.go
│   │   ├── system_config_repo.go
│   │   ├── notification_repo.go
│   │   ├── price_alert_repo.go
│   │   ├── report_repo.go
│   │   ├── nft_repository.go
│   │   ├── strategy_repository.go
│   │   ├── achievement_repository.go
│   │   ├── transaction_repository.go
│   │   ├── daily_summary_repository.go
│   │   ├── sync_metadata_repository.go
│   │   └── goal_repo.go
│   ├── models/                       # 数据模型（18 个表）
│   │   ├── base.go                  # BaseModel + 常量定义
│   │   ├── exchange_account.go
│   │   ├── wallet_address.go
│   │   ├── manual_asset.go
│   │   ├── asset_detail.go
│   │   ├── asset_snapshot.go
│   │   ├── exchange_rate.go
│   │   ├── system_config.go
│   │   ├── notification.go
│   │   ├── price_alert.go
│   │   ├── report.go
│   │   ├── nft.go
│   │   ├── strategy.go
│   │   ├── achievement.go
│   │   ├── transaction.go
│   │   ├── sync_metadata.go
│   │   └── goal.go
│   ├── integrations/                 # 第三方集成（8 个模块）
│   │   ├── base_client.go           # HTTP 基础客户端
│   │   ├── binance/                 # Binance API（go-binance SDK）
│   │   ├── okx/                     # OKX API（CCXT）
│   │   ├── coinbase/                # Coinbase API（CCXT）
│   │   ├── etherscan/               # Etherscan 兼容 API（6 链）
│   │   ├── coingecko/               # CoinGecko 价格 API
│   │   ├── yahoo/                   # Yahoo Finance（法币汇率）
│   │   ├── alchemy/                 # Alchemy NFT API
│   │   └── defi/                    # DeFi 协议（7 个）
│   │       ├── interface.go         # 协议接口定义
│   │       ├── registry.go          # 协议注册表
│   │       ├── position.go          # 仓位数据结构
│   │       ├── lido.go              # Lido（ETH 质押）
│   │       ├── rocketpool.go        # Rocket Pool（ETH 质押）
│   │       ├── aave.go              # Aave（借贷）
│   │       ├── compound.go          # Compound（借贷）
│   │       ├── uniswap_v2.go        # Uniswap V2（DEX）
│   │       ├── uniswap_v3.go        # Uniswap V3（DEX）
│   │       └── curve.go             # Curve（稳定币 DEX）
│   ├── cron/                         # 定时任务（6 个）
│   │   ├── snapshot_job.go          # 资产快照 + CronManager
│   │   ├── notification_job.go      # 通知摘要
│   │   ├── price_alert_job.go       # 价格预警检查
│   │   ├── report_job.go            # 自动报告生成
│   │   ├── strategy_job.go          # 策略监控
│   │   └── risk_alert_job.go        # 风险预警
│   ├── model/                        # 业务传输对象（DTO）
│   └── utils/                        # 工具库
│       ├── config.go                # 配置加载（GoFrame gcfg 封装）
│       ├── crypto.go                # AES-256-GCM 加密/解密
│       ├── response.go              # 统一响应格式 + 错误码
│       ├── cache.go                 # 内存缓存
│       ├── retry.go                 # 重试工具
│       └── ratelimit.go             # 速率限制
├── docs/
│   ├── swagger.yaml                  # OpenAPI 3.0 规范
│   └── swagger.go                    # Swagger UI 处理器
└── data/
    └── allfi.db                      # SQLite 数据库文件
```

---

## 3. 认证系统

### 3.1 单用户 PIN 锁

AllFi 是自部署的单用户应用，采用 PIN 码（4-8 位数字）作为认证方式：

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/auth/status` | GET | 查询 PIN 是否已设置 |
| `/api/v1/auth/setup` | POST | 首次设置 PIN（返回 JWT Token） |
| `/api/v1/auth/login` | POST | 验证 PIN（返回 JWT Token） |
| `/api/v1/auth/change` | POST | 修改 PIN（需验证旧 PIN） |

### 3.2 认证流程

```
首次使用 → PIN 未设置 → 所有接口开放（无需认证）
         → 用户调用 /auth/setup → PIN 存入数据库（bcrypt 哈希）
         → 后续访问需在 Header 中携带 Bearer Token

已设置 PIN → 调用 /auth/login → 返回 JWT Token
           → 请求头: Authorization: Bearer <token>
           → 中间件验证 Token → 放行或拒绝
```

### 3.3 白名单路由

以下路由无需认证即可访问：
- `/api/v1/health` — 健康检查
- `/api/v1/auth/*` — 认证相关接口

### 3.4 安全措施

- PIN 使用 `bcrypt` 哈希存储（非明文）
- JWT Token 有过期时间
- 连续错误登录触发账户锁定（`ErrAccountLocked`）

---

## 4. API 设计

### 4.1 统一响应格式

```json
{
  "code": 0,
  "message": "成功",
  "data": {},
  "timestamp": 1707382800
}
```

### 4.2 错误码体系

| 范围 | 类别 | 示例 |
|------|------|------|
| 0 | 成功 | `0` 成功 |
| 1001-1999 | 客户端错误 | `1001` 参数错误、`1003` 资源不存在、`1005` 未授权 |
| 2001-2999 | 服务器错误 | `2001` 内部错误、`2002` 数据库错误、`2003` 外部 API 错误 |
| 3001-3999 | 业务错误 | `3001` 交易所 API 错误、`3004` 地址无效、`3006` 汇率获取失败 |

### 4.3 完整路由表（75 条）

#### 认证（4 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/status` | 获取认证状态 |
| POST | `/api/v1/auth/setup` | 首次设置 PIN |
| POST | `/api/v1/auth/login` | 验证 PIN 登录 |
| POST | `/api/v1/auth/change` | 修改 PIN |

#### 资产（8 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/assets/summary` | 资产总览（支持 `?currency=` 参数） |
| GET | `/api/v1/assets/details` | 资产明细列表 |
| GET | `/api/v1/assets/history` | 历史资产数据（支持 `?days=` 参数） |
| POST | `/api/v1/assets/refresh` | 刷新所有资产数据 |
| GET | `/api/v1/assets/manual` | 获取传统资产列表 |
| POST | `/api/v1/assets/manual` | 添加传统资产 |
| PUT | `/api/v1/assets/manual/{id}` | 更新传统资产 |
| DELETE | `/api/v1/assets/manual/{id}` | 删除传统资产 |

#### 交易所（7 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/exchanges/accounts` | 获取所有交易所账户 |
| POST | `/api/v1/exchanges/accounts` | 添加交易所账户 |
| GET | `/api/v1/exchanges/accounts/{id}` | 获取单个账户详情 |
| PUT | `/api/v1/exchanges/accounts/{id}` | 更新账户 |
| DELETE | `/api/v1/exchanges/accounts/{id}` | 删除账户 |
| POST | `/api/v1/exchanges/accounts/{id}/test` | 测试 API 连接 |
| GET | `/api/v1/exchanges/accounts/{id}/balances` | 获取账户余额 |

#### 钱包（8 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/wallets/addresses` | 获取所有钱包地址 |
| POST | `/api/v1/wallets/addresses` | 添加钱包地址 |
| GET | `/api/v1/wallets/addresses/{id}` | 获取单个钱包详情 |
| PUT | `/api/v1/wallets/addresses/{id}` | 更新钱包 |
| DELETE | `/api/v1/wallets/addresses/{id}` | 删除钱包 |
| GET | `/api/v1/wallets/addresses/{id}/balances` | 获取钱包余额 |
| POST | `/api/v1/wallets/batch` | 批量导入钱包地址 |

#### 汇率（3 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/rates/current` | 获取当前汇率 |
| GET | `/api/v1/rates/prices` | 批量获取价格 |
| POST | `/api/v1/rates/refresh` | 刷新汇率缓存 |

#### 通知（9 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/notifications` | 获取通知列表（分页） |
| GET | `/api/v1/notifications/unread-count` | 获取未读数量 |
| POST | `/api/v1/notifications/{id}/read` | 标记已读 |
| POST | `/api/v1/notifications/read-all` | 全部标记已读 |
| GET | `/api/v1/notifications/preferences` | 获取通知偏好 |
| PUT | `/api/v1/notifications/preferences` | 更新通知偏好 |
| GET | `/api/v1/notifications/webpush/vapid` | 获取 VAPID 公钥 |
| POST | `/api/v1/notifications/webpush/subscribe` | 订阅 WebPush |
| POST | `/api/v1/notifications/webpush/unsubscribe` | 取消订阅 |

#### 价格预警（4 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/alerts` | 创建预警 |
| GET | `/api/v1/alerts` | 获取预警列表 |
| PUT | `/api/v1/alerts/{id}` | 更新预警 |
| DELETE | `/api/v1/alerts/{id}` | 删除预警 |

#### 报告（3 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/reports` | 获取报告列表 |
| GET | `/api/v1/reports/{id}` | 获取单个报告 |
| POST | `/api/v1/reports/generate` | 生成报告（日报/周报/月报/年报） |

#### 交易记录（5 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/transactions` | 分页查询交易记录（支持游标分页） |
| POST | `/api/v1/transactions/sync` | 同步交易记录 |
| GET | `/api/v1/transactions/stats` | 获取交易统计 |
| GET | `/api/v1/settings/tx-sync` | 获取同步设置 |
| PUT | `/api/v1/settings/tx-sync` | 更新同步设置 |

#### 策略引擎（5 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/strategies` | 获取策略列表 |
| POST | `/api/v1/strategies` | 创建策略 |
| PUT | `/api/v1/strategies/{id}` | 更新策略 |
| DELETE | `/api/v1/strategies/{id}` | 删除策略 |
| GET | `/api/v1/strategies/{id}/analysis` | 分析再平衡需求 |

#### DeFi（2 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/defi/positions` | 获取 DeFi 仓位（支持链/协议过滤） |
| GET | `/api/v1/defi/protocols` | 获取支持的协议列表 |

#### NFT（1 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/nft/assets` | 获取 NFT 资产列表 |

#### 成就系统（2 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/achievements` | 获取成就列表（含解锁状态） |
| POST | `/api/v1/achievements/check` | 检查并解锁新成就 |

#### 目标追踪（4 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/goals` | 获取目标列表（附进度） |
| POST | `/api/v1/goals` | 创建目标 |
| PUT | `/api/v1/goals/{id}` | 更新目标 |
| DELETE | `/api/v1/goals/{id}` | 删除目标 |

#### 数据分析（5 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/analytics/fees` | 费用分析 |
| GET | `/api/v1/analytics/pnl/daily` | 每日盈亏 |
| GET | `/api/v1/analytics/pnl/summary` | 盈亏汇总 |
| GET | `/api/v1/analytics/attribution` | 资产归因分析 |
| GET | `/api/v1/analytics/forecast` | 趋势预测 |

#### 其他（7 条）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health` | 健康检查 |
| GET | `/api/v1/benchmark` | 基准对比（vs BTC/ETH/S&P500） |
| GET | `/api/v1/market/gas` | Gas 费查询 |
| GET | `/api/v1/portfolio/health` | 资产健康评分 |
| GET | `/api/v1/docs` | Swagger UI |
| GET | `/api/v1/docs/swagger.yaml` | OpenAPI 规范文件 |

---

## 5. 数据模型（18 张表）

所有表继承 `BaseModel`，包含 `id`（主键）、`created_at`、`updated_at`、`deleted_at`（软删除）字段。

### 5.1 核心业务表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `exchange_accounts` | 交易所账户 | exchange, label, api_key_encrypted, api_secret_encrypted, status |
| `wallet_addresses` | 钱包地址 | blockchain, address, label, status |
| `manual_assets` | 传统资产 | asset_type(bank/cash/stock/fund), asset_name, institution, amount(Decimal), currency |
| `asset_details` | 资产明细 | source_type, symbol, name, balance, price_usd, value_usd |
| `asset_snapshots` | 资产快照 | total_value_usd, exchange_value, blockchain_value, manual_value |
| `exchange_rates` | 汇率缓存 | from_currency, to_currency, rate, source, expires_at |
| `system_config` | 系统配置 | config_key, config_value |

### 5.2 功能扩展表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| `notifications` | 通知消息 | user_id, type, title, content, is_read |
| `price_alerts` | 价格预警 | symbol, condition(above/below), target_price, is_active |
| `reports` | 资产报告 | user_id, report_type(daily/weekly/monthly/annual), content |
| `unified_transactions` | 统一交易记录 | tx_type, source, source_id, from_asset, to_asset, value_usd, tx_hash |
| `transaction_daily_summaries` | 交易日汇总 | date, buy_count, sell_count, total_fee_usd, net_flow_usd |
| `strategies` | 投资策略 | name, target_allocations, rebalance_threshold |
| `achievements` | 成就记录 | achievement_id, unlocked_at |
| `nfts` | NFT 资产 | chain, contract_address, token_id, name, collection |
| `goals` | 目标追踪 | name, target_amount, current_amount, deadline |
| `sync_metadata` | 同步元数据 | source, last_sync_at, sync_status |

### 5.3 支持的常量

```go
// 交易所: binance, okx, coinbase
// 区块链: ethereum, bsc, polygon
// 资产来源: cex, blockchain, manual
// 传统资产类型: bank, cash, stock, fund
// 计价货币: USDC, USD, BTC, ETH, CNY
// 交易类型: buy, sell, swap, transfer, deposit, withdraw
```

---

## 6. 服务层架构

### 6.1 依赖注入

通过 `ServiceContainer` 管理所有服务实例，构造时注入依赖：

```go
type ServiceContainer struct {
    ExchangeService     ExchangeServiceInterface
    BlockchainService   BlockchainServiceInterface
    PriceService        PriceServiceInterface
    AssetService        AssetServiceInterface
    SnapshotService     SnapshotServiceInterface
    NotificationService NotificationServiceInterface
    PriceAlertService   PriceAlertServiceInterface
    ReportService       ReportServiceInterface
    DeFiService         DeFiServiceInterface
    NFTService          NFTServiceInterface
    TransactionService  TransactionServiceInterface
    FeeService          FeeServiceInterface
    StrategyService     StrategyServiceInterface
    AchievementService  AchievementServiceInterface
    BenchmarkService    BenchmarkServiceInterface
}
```

### 6.2 服务接口总览

| 服务 | 接口 | 核心方法 |
|------|------|---------|
| AuthService | `AuthServiceInterface` | SetupPIN, VerifyPIN, ChangePIN, ValidateToken |
| AssetService | `AssetServiceInterface` | GetSummary, GetDetails, GetHistory, RefreshAll, 手动资产 CRUD |
| ExchangeService | `ExchangeServiceInterface` | CRUD, TestConnection, GetBalances, RefreshAllBalances |
| BlockchainService | `BlockchainServiceInterface` | CRUD, GetBalances, RefreshAllBalances, BatchImportWallets |
| PriceService | `PriceServiceInterface` | GetPrice, GetPrices, GetExchangeRate, ConvertValue, RefreshRates |
| SnapshotService | `SnapshotServiceInterface` | CreateSnapshot, GetSnapshots, CleanupOldSnapshots |
| NotificationService | `NotificationServiceInterface` | Send, List, MarkRead, MarkAllRead, GetPreference, GenerateDailyDigest |
| PriceAlertService | `PriceAlertServiceInterface` | CRUD, CheckAlerts |
| ReportService | `ReportServiceInterface` | GetReports, GenerateDaily/Weekly/Monthly/AnnualReport |
| DeFiService | `DeFiServiceInterface` | GetPositions, GetTotalValue, GetSupportedProtocols |
| NFTService | `NFTServiceInterface` | GetNFTs, GetCollections, GetTotalValue |
| TransactionService | `TransactionServiceInterface` | GetTransactions(游标分页), SyncTransactions, GetStats, GetSyncSettings |
| FeeService | `FeeServiceInterface` | GetFeeAnalytics |
| StrategyService | `StrategyServiceInterface` | CRUD, AnalyzeRebalance, CheckStrategies |
| AchievementService | `AchievementServiceInterface` | GetAchievements, CheckAndUnlock |
| BenchmarkService | `BenchmarkServiceInterface` | GetBenchmark（7d/30d/90d/1y） |
| GoalService | `GoalServiceInterface` | CRUD + GetGoals(附带进度百分比) |
| HealthScoreService | `HealthScoreServiceInterface` | GetHealthScore |
| RiskAlertService | — | 风险预警后台检查 |
| WebPushService | — | 浏览器推送通知 |

---

## 7. 第三方集成

### 7.1 交易所集成

| 交易所 | SDK | 说明 |
|--------|-----|------|
| Binance | `go-binance/v2` | 专用 SDK，现货余额查询 |
| OKX | `ccxt/go/v4` | 通过 CCXT 统一接口 |
| Coinbase | `ccxt/go/v4` | 通过 CCXT 统一接口 |

API Key 使用 AES-256-GCM 加密存储，添加账户时验证连接。

### 7.2 区块链集成

| 服务 | 说明 |
|------|------|
| Etherscan | 统一客户端，支持 6 条 EVM 链（Ethereum/BSC/Polygon/Arbitrum/Optimism/Base） |
| Alchemy | NFT 数据查询（获取 NFT 列表、收藏集、估值） |

### 7.3 DeFi 协议集成

通过 `defi/registry.go` 统一注册，所有协议实现相同接口：

| 协议 | 文件 | 类型 |
|------|------|------|
| Lido | `lido.go` | ETH 质押 |
| Rocket Pool | `rocketpool.go` | ETH 质押 |
| Aave | `aave.go` | 借贷 |
| Compound | `compound.go` | 借贷 |
| Uniswap V2 | `uniswap_v2.go` | DEX 流动性 |
| Uniswap V3 | `uniswap_v3.go` | DEX 流动性 |
| Curve | `curve.go` | 稳定币 DEX |

### 7.4 价格数据

| 数据源 | 用途 |
|--------|------|
| CoinGecko | 加密货币价格（BTC/ETH/USDC 等） |
| Yahoo Finance | 法币汇率（USD/CNY 等）、传统资产价格 |

---

## 8. 定时任务

由 `CronManager` 统一管理，启动时开始执行，支持优雅停止：

| 任务 | 文件 | 默认间隔 | 说明 |
|------|------|---------|------|
| 资产快照 | `snapshot_job.go` | 1 小时 | 刷新价格 → 创建快照 → 清理旧数据（90 天） |
| 通知摘要 | `notification_job.go` | — | 为启用摘要的用户生成每日通知 |
| 价格预警 | `price_alert_job.go` | — | 检查活跃预警条件并触发通知 |
| 报告生成 | `report_job.go` | — | 自动生成日报/周报/月报 |
| 策略监控 | `strategy_job.go` | — | 检查策略偏移，触发再平衡提醒 |
| 风险预警 | `risk_alert_job.go` | — | 检查资产集中度、波动率等风险指标 |

---

## 9. 安全设计

### 9.1 API Key 加密

```go
// AES-256-GCM 加密
// 主密钥: 32 字节，从环境变量 ALLFI_MASTER_KEY 或配置文件读取
// 存储格式: Base64(nonce + ciphertext)
// 全局实例: utils.EncryptAPIKey() / utils.DecryptAPIKey()
```

### 9.2 PIN 码存储

- 使用 `golang.org/x/crypto/bcrypt` 哈希
- 支持 4-8 位纯数字
- 连续错误登录触发锁定

### 9.3 中间件链

请求经过的中间件顺序（从外到内）：

```
Recovery → CORS → Logger → Auth → Handler
```

---

## 10. 配置管理

使用 GoFrame `gcfg` 管理配置，支持环境变量覆盖：

| 配置项 | 环境变量 | 默认值 | 说明 |
|--------|---------|--------|------|
| `server.port` | `ALLFI_PORT` | 8080 | 服务端口 |
| `server.mode` | `ALLFI_MODE` | development | 运行模式 |
| `database.type` | `ALLFI_DB_TYPE` | sqlite | 数据库类型 |
| `database.sqlite.path` | `ALLFI_DB_PATH` | ../data/allfi.db | SQLite 路径 |
| `security.master_key` | `ALLFI_MASTER_KEY` | — | 加密主密钥（必填） |
| `external_apis.etherscan.api_key` | `ETHERSCAN_API_KEY` | — | Etherscan API Key |
| `external_apis.bscscan.api_key` | `BSCSCAN_API_KEY` | — | BscScan API Key |
| `external_apis.coingecko.api_key` | `COINGECKO_API_KEY` | — | CoinGecko API Key |
| `cron.snapshot_interval` | — | 3600 | 快照间隔（秒） |
| `cron.price_cache_ttl` | — | 300 | 价格缓存时间（秒） |
| `defaults.currency` | — | USDC | 默认计价货币 |
| `defaults.history_retention_days` | — | 90 | 历史数据保留天数 |

---

## 11. API 文档

内置 Swagger UI，启动后访问：

- **Swagger UI**: `http://localhost:8080/api/v1/docs`
- **OpenAPI YAML**: `http://localhost:8080/api/v1/docs/swagger.yaml`

---

## 12. 数据统计

| 维度 | 数量 |
|------|------|
| API 路由 | 75 条 |
| Handler 文件 | 24 个 |
| 服务接口 | 20+ 个 |
| Repository | 18 个 |
| 数据模型 | 18 张表 |
| 第三方集成 | 8 个模块 |
| DeFi 协议 | 7 个 |
| 定时任务 | 6 个 |
| 中间件 | 4 个 |
| 支持交易所 | 3 个（Binance/OKX/Coinbase） |
| 支持区块链 | 3 条（Ethereum/BSC/Polygon） |
| Etherscan 兼容链 | 6 条（+Arbitrum/Optimism/Base） |
| 计价货币 | 5 种（USDC/USD/BTC/ETH/CNY） |

---

**文档维护者**: @allfi
**最后更新**: 2026-02-11
