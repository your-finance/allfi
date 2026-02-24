# CLAUDE.md — AllFi 后端开发指南

本文件为 Claude Code 在 AllFi 后端（`core/`）目录中工作时提供指导。

---

## 中文优先原则

- **代码注释必须使用中文** — 函数、方法、结构体、类型的注释均用中文
- **文档内容必须使用中文** — Markdown 文档正文用中文
- **例外** — 代码标识符使用英文；技术术语可保留英文（如 API、HTTP、JSON）
- **编码规范**: 正确使用 UTF-8 编码，确保中文显示正常
- **提交代码**: 每次完成一个任务，提交一次代码
- **禁止操作**:
  - ❌ **禁止运行** `gf gen ctrl` 或 `make ctrl` 命令
  - ❌ **禁止自动生成** 控制器文件（会破坏现有的模块化控制器结构）
- **Go 版本要求**:
  - 项目使用 **Go 1.24**
  - go.mod 中保持 `go 1.24`
  - 所有依赖包必须与 Go 1.24 兼容
  - 编译和构建时使用 Go 1.24 工具链
- **效率最高**: 高效率代码
  - 第三方库用最优的，如 `encoding/json` 换成 `github.com/goccy/go-json`

---

## 角色与思考

做为资深后端开发工程师，精通 GoFrame 设计和金融量化交易系统架构，从设计角度出发，考虑模块化、分层管理、性能优化和金融业务特点。
对于每次要求/需求，始终基于现有代码生成 TODO list，并优先考虑最佳实践。

---

## 技术栈

| 分类 | 技术 |
|------|------|
| 语言 | Go 1.24 |
| 框架 | GoFrame v2（HTTP 路由 + ORM + 配置） |
| ORM | GoFrame 原生 ORM（SQLite3 默认 / MySQL 可选） |
| 数据库驱动 | `github.com/gogf/gf/contrib/drivers/sqlite/v2`（纯 Go 实现，无需 CGO） |
| 配置 | GoFrame `gcfg`（YAML 配置 + 环境变量覆盖） |
| 加密 | AES-256-GCM（API Key 加密存储） |
| 认证 | PIN 码 bcrypt + JWT Bearer Token |

### 数据库说明

- **SQLite 驱动**：采用 `glebarez/go-sqlite`（封装 `modernc.org/sqlite`），**纯 Go 实现**，编译时无需 CGO
- **编译命令**：`CGO_ENABLED=0 go build ./...` 即可完成静态编译
- **已清除 GORM**：项目已完全移除 `gorm.io/*` 依赖，统一使用 GoFrame ORM
- **配置方式**：`manifest/config/config.yaml` 中 `database.default.link` 配置连接串

---

### GoFrame 最佳实践

参考官方文档进行规范化开发：
- [模块化设计](https://goframe.org/docs/design/modular) - 保持业务模块独立性
- [统一框架设计](https://goframe.org/docs/design/unified) - 遵循框架设计理念
- [数据规范-gf gen dao](https://goframe.org/docs/cli/gen-dao) - 数据访问层代码生成
- [模块规范-gf gen service](https://goframe.org/docs/cli/gen-service) - 服务接口自动生成

### 分层架构原则

**严格遵循**: Controller → Service → Logic → DAO → Database

#### 开发规范
- **Controller 层**（位于 `controller/controller.go`）:
  - 接收请求参数，调用 Service 接口
  - 函数返回必须包含返回参数名称与类型
  - 通过 GoFrame `g.Meta` 标签定义路由
- **Service 层**（位于 `service/{module}.go`）:
  - 基于接口设计（`I{Module}`），便于测试和扩展
  - `Register{Module}()` + `{Module}()` 注册和获取模式
- **Logic 层**（位于 `logic/{module}.go`）:
  - 专注业务逻辑实现
  - `logic/logic.go` 中通过 `init()` 自动注册服务
  - 新模块必须在 `cmd/server/main.go` 中添加 blank import
- **DAO 层**（全局 `internal/dao/` + 模块级 `app/{module}/dao/`）:
  - 全局 `internal/dao/` 由 `gf gen dao` 自动生成
  - `dao/internal/` 为自动生成代码，不可手动修改
  - 每个拥有表的模块在 `app/{module}/dao/` 中创建引用文件
  - 模块级 DAO 引用模式：`var Table = &globalDao.Table`
  - 跨模块访问使用别名导入：`assetDao "your-finance/allfi/internal/app/asset/dao"`

---

## 项目结构

```
core/
├── cmd/server/main.go              # GoFrame 服务入口（g.Server + init 自动注册）
├── api/v1/                          # API 请求/响应定义（25 个模块）
│   ├── achievement/achievement.go   # Req/Res + g.Meta 路由标签
│   ├── asset/asset.go
│   ├── attribution/attribution.go
│   ├── auth/auth.go
│   ├── benchmark/benchmark.go
│   ├── defi/defi.go
│   ├── exchange/exchange.go
│   ├── exchange_rate/exchange_rate.go
│   ├── fee/fee.go
│   ├── forecast/forecast.go
│   ├── goal/goal.go
│   ├── health/health.go
│   ├── health_score/health_score.go
│   ├── manual_asset/manual_asset.go
│   ├── market/market.go
│   ├── nft/nft.go
│   ├── notification/notification.go
│   ├── pnl/pnl.go
│   ├── price_alert/price_alert.go
│   ├── report/report.go
│   ├── strategy/strategy.go
│   ├── transaction/transaction.go
│   ├── user/user.go
│   ├── wallet/wallet.go
│   └── webpush/webpush.go
├── hack/config.yaml                 # GoFrame CLI 配置（gf gen dao）
├── manifest/
│   └── config/config.yaml           # 应用配置
├── internal/
│   ├── app/                         # 25 个业务模块（每个模块自闭环）
│   │   ├── {module}/
│   │   │   ├── service/{module}.go  # I{Module} 接口 + Register/Get
│   │   │   ├── logic/
│   │   │   │   ├── logic.go         # init() 自动注册服务
│   │   │   │   └── {module}.go      # s{Module} 实现
│   │   │   ├── dao/{table}.go        # 模块级 DAO 引用（引用全局 DAO）
│   │   │   ├── model/{module}.go    # 业务 DTO + 常量
│   │   │   └── controller/controller.go  # Controller + Register(group)
│   │   ├── achievement/             # 成就系统
│   │   ├── asset/                   # 资产聚合
│   │   ├── attribution/             # 归因分析
│   │   ├── auth/                    # 认证（PIN + JWT）
│   │   ├── benchmark/               # 基准对比
│   │   ├── defi/                    # DeFi 仓位
│   │   ├── exchange/                # 交易所账户
│   │   ├── exchange_rate/           # 汇率
│   │   ├── fee/                     # 费用分析
│   │   ├── forecast/                # 趋势预测
│   │   ├── goal/                    # 目标追踪
│   │   ├── health/                  # 健康检查
│   │   ├── health_score/            # 健康评分
│   │   ├── manual_asset/            # 手动资产
│   │   ├── market/                  # 市场数据
│   │   ├── nft/                     # NFT 资产
│   │   ├── notification/            # 通知系统
│   │   ├── pnl/                     # 盈亏分析
│   │   ├── price_alert/             # 价格预警
│   │   ├── report/                  # 资产报告
│   │   ├── strategy/                # 策略引擎
│   │   ├── transaction/             # 交易记录
│   │   ├── user/                    # 用户设置
│   │   ├── wallet/                  # 钱包地址
│   │   └── webpush/                 # WebPush 推送
│   ├── dao/                         # GoFrame DAO（gf gen dao 生成，19 张表）
│   │   ├── internal/                # 自动生成（不手动修改）
│   │   └── {table}.go              # 自定义查询扩展
│   ├── model/                       # GoFrame 实体模型
│   │   ├── entity/                  # 数据库实体（gf gen dao 生成）
│   │   └── do/                      # 查询对象（gf gen dao 生成）
│   ├── middleware/                   # 全局中间件
│   │   ├── auth.go                  # JWT 认证
│   │   ├── context.go               # 上下文注入
│   │   ├── cors.go                  # 跨域配置
│   │   ├── error_handler.go         # 全局错误处理
│   │   ├── logger.go                # 请求日志
│   │   ├── response.go              # 统一响应包装
│   │   └── middleware.go            # Register(s) 全局注册入口
│   ├── integrations/                # 第三方 API 集成
│   │   ├── binance/                 # Binance 交易所
│   │   ├── okx/                     # OKX 交易所
│   │   ├── coinbase/                # Coinbase 交易所
│   │   ├── etherscan/               # Etherscan
│   │   ├── alchemy/                 # Alchemy（NFT 数据）
│   │   ├── coingecko/               # CoinGecko（价格数据）
│   │   ├── yahoo/                   # Yahoo Finance
│   │   ├── defi/                    # DeFi 协议集成
│   │   └── base_client.go           # 集成基类
│   ├── consts/                      # 全局常量定义
│   ├── cron/                        # 定时任务（待迁移到新 service 接口）
│   └── utils/                       # 工具函数
│       ├── config.go                # 配置加载
│       ├── crypto.go                # AES-256-GCM 加解密
│       ├── cache.go                 # 内存缓存
│       ├── ratelimit.go             # 请求限流
│       └── retry.go                 # 重试机制
└── go.mod
```

---

## 架构分层

项目采用**纯 GoFrame 模块化架构**：

```
api/v1/{module}/ → app/{module}/controller/ → app/{module}/service/ → app/{module}/logic/ → app/{module}/dao/ → Database
```

- **api/v1/**：请求/响应定义，使用 `g.Meta` 标签声明路由
- **controller/**：Controller（接收请求，调用 Service）
- **service/**：接口定义（`I{Module}`），解耦上下层
- **logic/**：核心业务逻辑实现
- **dao/**：数据访问层（GoFrame ORM，由 `gf gen dao` 生成）
- **middleware/**：全局中间件（CORS → Context → Logger → ErrorHandler → ResponseWrapper）

### 模块级 DAO 架构

项目采用**两层 DAO 设计**：全局 DAO（`internal/dao/`）+ 模块级 DAO（`app/{module}/dao/`）。

#### 全局 DAO 层

由 `gf gen dao` 自动生成，位于 `internal/dao/`：
- `internal/dao/internal/` — 自动生成的内部实现（**禁止手动修改**）
- `internal/dao/{table}.go` — 对外暴露的 DAO 变量 + 可添加自定义查询

#### 模块级 DAO 层

每个拥有数据表的模块在 `app/{module}/dao/` 中创建引用文件，指向全局 DAO。

**标准模板**（以 exchange 模块为例）：

```go
// Package dao 交易所模块 DAO 封装
// 对全局 DAO 的模块级引用，供本模块和跨模块调用
package dao

import (
    globalDao "your-finance/allfi/internal/dao"
)

// ExchangeAccounts 交易所账户表访问对象（引用全局 DAO）
var ExchangeAccounts = &globalDao.ExchangeAccounts
```

**Go `internal` 包可见性规则**：`internal/app/{module}/dao/` **不能**直接导入 `internal/dao/internal/`（因为不是其父目录），必须通过 `globalDao` 别名引用全局 `internal/dao/` 包。

#### 表的模块归属

| 模块 | 表名 | DAO 引用文件 |
|------|------|-------------|
| asset | asset_details, asset_snapshots | `app/asset/dao/asset_details.go`, `asset_snapshots.go` |
| exchange | exchange_accounts | `app/exchange/dao/exchange_accounts.go` |
| exchange_rate | exchange_rates | `app/exchange_rate/dao/exchange_rates.go` |
| wallet | wallet_addresses | `app/wallet/dao/wallet_addresses.go` |
| manual_asset | manual_assets | `app/manual_asset/dao/manual_assets.go` |
| goal | goals | `app/goal/dao/goals.go` |
| nft | nft_caches | `app/nft/dao/nft_caches.go` |
| notification | notifications, notification_preferences | `app/notification/dao/notifications.go`, `notification_preferences.go` |
| price_alert | price_alerts | `app/price_alert/dao/price_alerts.go` |
| report | reports | `app/report/dao/reports.go` |
| strategy | strategies | `app/strategy/dao/strategies.go` |
| achievement | user_achievements | `app/achievement/dao/user_achievements.go` |
| transaction | unified_transactions, sync_metadata, transaction_daily_summaries, system_config | `app/transaction/dao/` 下 4 个文件 |
| webpush | web_push_subscriptions | `app/webpush/dao/web_push_subscriptions.go` |

**无表模块**（不需要 `dao/` 目录）：auth, attribution, benchmark, defi, fee, forecast, health, health_score, market, pnl

#### 跨模块 DAO 访问

当模块 A 需要访问模块 B 的数据表时，使用**别名导入**模块 B 的 DAO 包：

```go
import (
    "your-finance/allfi/internal/app/achievement/dao"          // 本模块 DAO
    assetDao "your-finance/allfi/internal/app/asset/dao"       // 跨模块：资产
    exchangeDao "your-finance/allfi/internal/app/exchange/dao" // 跨模块：交易所
)

// 本模块表
dao.UserAchievements.Ctx(ctx).Where(...)

// 跨模块表
assetDao.AssetDetails.Ctx(ctx).Where(...)
exchangeDao.ExchangeAccounts.Ctx(ctx).Where(...)
```

**跨模块依赖关系**：

| 模块 | 依赖的外部 DAO |
|------|---------------|
| achievement | exchangeDao, walletDao, assetDao, strategyDao, priceAlertDao, reportDao |
| health_score | assetDao (AssetDetails) |
| fee | transactionDao (UnifiedTransactions) |
| benchmark | assetDao (AssetSnapshots) |
| attribution | assetDao (AssetSnapshots, AssetDetails) |
| forecast | assetDao (AssetSnapshots) |
| defi | walletDao (WalletAddresses) |
| pnl | assetDao (AssetSnapshots), transactionDao (TransactionDailySummaries) |
| strategy | assetDao (AssetDetails) |
| report | assetDao (AssetSnapshots, AssetDetails) |
| nft | walletDao (WalletAddresses) |
| goal | assetService (GetSummary), assetDao (AssetSnapshots, AssetDetails) |
| user | transactionDao (SystemConfig) |
| asset | exchangeRateDao (ExchangeRates), coingecko |
| manual_asset | exchangeRateDao (ExchangeRates), coingecko |
| notification | webpushService (SendPush) |

### 服务注册模式

每个模块通过 `init()` 自动注册服务：

```go
// logic/logic.go
func init() {
    service.Register{Module}(New())
}
```

`cmd/server/main.go` 通过 blank import 触发注册，同时导入 controller 包注册路由：

```go
// blank import 触发 init() 注册服务
_ "your-finance/allfi/internal/app/{module}/logic"

// 导入 controller 包注册路由
{module}Ctrl "your-finance/allfi/internal/app/{module}/controller"
```

路由注册示例：

```go
s.Group("/api/v1", func(group *ghttp.RouterGroup) {
    healthCtrl.Register(group)       // 免认证
    authCtrl.Register(group)         // 免认证
    group.Middleware(middleware.Auth) // 认证中间件
    assetCtrl.Register(group)        // 需认证
    exchangeCtrl.Register(group)
    // ... 其余模块
})
```

### 统一响应格式

由 `middleware/response.go` 中的 `MiddlewareHandlerResponse` 自动包装：

```json
{"code": 0, "message": "success", "data": {...}, "timestamp": 1707580800}
```

#### 响应中间件规范（参照 GoFrame 官方实现）

`middleware/response.go` 参照 GoFrame v2 官方 `ghttp.MiddlewareHandlerResponse` 实现，在其基础上增加了 `timestamp` 字段。开发时须遵守以下规范：

1. **使用 `r.Response.WriteJson()` 输出响应** — 禁止使用 `json.Marshal` + `WriteOver` 手动序列化，`WriteJson` 会自动设置 `Content-Type` 和处理序列化错误
2. **不要自定义 `defer recover()`** — GoFrame `ghttp.Server` 内部已有 panic recovery 机制，自定义 recover 会与框架冲突
3. **响应已写入检查必须同时使用两个条件**：
   ```go
   if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
       return  // 已有响应，跳过统一包装
   }
   ```
   - `BufferLength()`：检查缓冲区中待发送的字节
   - `BytesWritten()`：检查已直接发送到客户端的字节（如 `WriteJsonExit`）
4. **不要调用 `r.SetError(nil)`** — 保留错误信息以便其他中间件也能获取
5. **错误消息始终使用 `err.Error()`** — 不要用 `code.Message()` 覆盖具体错误信息
6. **无明确错误码时使用 `gcode.CodeInternalError`（code=50）** — 与 GoFrame 官方一致
7. **处理非 200 HTTP 状态码** — 404→`gcode.CodeNotFound`，403→`gcode.CodeNotAuthorized`

#### 业务错误码定义

| 范围 | 含义 | 说明 |
|------|------|------|
| `0` | 成功 | `gcode.CodeOK` |
| `50` | 内部错误 | `gcode.CodeInternalError`（无明确错误码时的默认值） |
| `51` | 参数校验失败 | `gcode.CodeValidationFailed`（GoFrame 自动校验） |
| `401` | 未认证 | auth 中间件直接 `WriteJsonExit` |
| `4001+` | 自定义业务错误 | 使用 `gerror.NewCode(gcode.New(4001, "", nil), "消息")` |

#### 业务逻辑中返回错误的写法

```go
// ✅ 带业务错误码 — 前端可据此做不同处理
return nil, gerror.NewCode(gcode.New(4001, "", nil), "您还没有添加任何资产，请先添加资产后再生成报告")

// ✅ 一般内部错误 — 中间件自动映射为 code=50
return nil, gerror.New("查询数据库失败")

// ✅ 包装底层错误 — 保留原始错误堆栈
return nil, gerror.Wrap(err, "保存报告失败")

// ❌ 禁止直接返回 fmt.Errorf — 丢失 GoFrame 错误码和堆栈
return nil, fmt.Errorf("保存失败: %v", err)
```

---

## 开发规范

### 新增 API 端点步骤

1. 在 `api/v1/{module}/` 中定义 Req/Res（含 `g.Meta` 路由标签）
2. 在 `service/{module}.go` 中添加接口方法
3. 在 `logic/{module}.go` 中实现业务逻辑（使用 `dao.{Table}` 查询）
4. 在 `controller/controller.go` 的 Controller 中添加处理方法
5. 如需新表：运行 `gf gen dao` 生成 DAO/Entity/DO

### 新增模块步骤

1. 创建 `api/v1/{module}/{module}.go` — Req/Res 定义（含 `g.Meta` 路由标签）
2. 创建 `internal/app/{module}/` 五层结构：
   - `controller/controller.go` — Controller + `Register(group)` 函数
   - `service/{module}.go` — `I{Module}` 接口 + `Register{Module}()` / `{Module}()`
   - `logic/logic.go` — `init()` 中调用 `service.Register{Module}(New())`
   - `logic/{module}.go` — `s{Module}` 结构体实现接口方法
   - `model/{module}.go` — 业务 DTO + 常量
   - `dao/{table}.go` — 模块级 DAO 引用（如有自有表）
3. 在 `cmd/server/main.go` 中：
   - 添加 `_ "your-finance/allfi/internal/app/{module}/logic"` blank import
   - 添加 `{module}Ctrl "your-finance/allfi/internal/app/{module}/controller"` 导入
   - 在路由注册中调用 `{module}Ctrl.Register(group)`
4. 如需新表：运行 `gf gen dao` 生成全局 DAO，然后在 `app/{module}/dao/` 创建引用文件

### GoFrame ORM 查询示例

```go
// 查询
dao.ExchangeAccounts.Ctx(ctx).
    Where(dao.ExchangeAccounts.Columns().UserId, userID).
    Scan(&accounts)

// 创建 — 使用 g.Map 方式（推荐，避免零值 ID 问题）
result, err := dao.Reports.Ctx(ctx).Data(g.Map{
    dao.Reports.Columns().UserId:   userID,
    dao.Reports.Columns().Type:     reportType,
    dao.Reports.Columns().Period:   period,
}).Insert()
id, _ := result.LastInsertId()

// 更新
dao.ExchangeAccounts.Ctx(ctx).
    Where(dao.ExchangeAccounts.Columns().Id, id).
    Data(account).Update()

// 分页 + 统计
dao.Table.Ctx(ctx).Page(page, pageSize).ScanAndCount(&list, &total, true)
```

#### ⚠️ Insert 避免零值 ID

**禁止**直接将 entity 结构体传入 `Insert()`，因为结构体的零值 `Id: 0` 会被包含在 SQL 中，可能在 SQLite `AUTOINCREMENT` 列上引发冲突：

```go
// ❌ 错误 — entity.Reports{Id: 0, ...} 会导致 INSERT INTO ... (id, ...) VALUES (0, ...)
result, err := dao.Reports.Ctx(ctx).Insert(report)

// ❌ 慎用 OmitEmpty — 会忽略所有零值字段（如 Change=0, ChangePercent=0）
result, err := dao.Reports.Ctx(ctx).OmitEmpty().Insert(report)

// ✅ 推荐 — 使用 Data(g.Map{...}) 显式指定要插入的字段，让数据库自动生成 ID
result, err := dao.Reports.Ctx(ctx).Data(g.Map{
    dao.Reports.Columns().UserId:        report.UserId,
    dao.Reports.Columns().Type:          report.Type,
    dao.Reports.Columns().TotalValue:    report.TotalValue,
    dao.Reports.Columns().Change:        report.Change,        // 允许零值
    dao.Reports.Columns().ChangePercent: report.ChangePercent,  // 允许零值
}).Insert()
```

### 安全要求

- 所有 CEX API Key/Secret/Passphrase 必须 AES-256-GCM 加密存储
- 数据库中禁止存储明文敏感数据
- 使用 GoFrame ORM 参数化查询，禁止拼接 SQL
- PIN 码使用 bcrypt 哈希存储

### 代码风格

- 代码注释必须使用中文
- `go fmt ./...` 格式化代码
- `go vet ./...` 静态检查
- 错误处理不可吞异常，必须返回或记录日志

---

## API 设计规范

前端 `webapp/src/api/client.js` 基于 `fetch` 封装，Base URL 为 `http://localhost:8080/api/v1`。所有后端路由必须严格匹配前端调用约定。

### URL 设计

**前缀**：所有路由以 `/api/v1/` 开头

**RESTful 资源命名**：

| 操作 | 方法 | URL 模式 | 示例 |
|------|------|----------|------|
| 列表 | GET | `/resource` | `GET /api/v1/exchanges/accounts` |
| 创建 | POST | `/resource` | `POST /api/v1/exchanges/accounts` |
| 详情 | GET | `/resource/{id}` | `GET /api/v1/exchanges/accounts/{id}` |
| 更新 | PUT | `/resource/{id}` | `PUT /api/v1/exchanges/accounts/{id}` |
| 删除 | DELETE | `/resource/{id}` | `DELETE /api/v1/exchanges/accounts/{id}` |

**命名规则**：
- 使用小写英文 + 连字符（kebab-case）
- 资源名用复数：`/accounts`、`/alerts`、`/strategies`

### 统一响应格式

由 `middleware/response.go` 自动包装，无需手动调用：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "timestamp": 1707580800
}
```

### JSON 字段命名

**统一使用 snake_case**，与前端预期一致：

```go
type Alert struct {
    ID          uint    `json:"id"`
    TargetPrice float64 `json:"target_price"`
    IsActive    bool    `json:"is_active"`
}
```

### 认证规范

- 前端从 `localStorage('allfi-auth')` 获取 JWT Token
- 请求头：`Authorization: Bearer <token>`
- `/auth/*` 和 `/health/*` 路由免认证，其余路由需 `middleware.Auth`
- 401 响应触发前端清除 Token 并跳转登录页
- **JWT 密钥来源**：auth logic 在首次设置 PIN 时生成随机密钥，存入 `system_config` 表（key=`auth.jwt_secret`）
- **middleware/auth.go** 从同一 DB 表读取密钥验证 Token（非配置文件 `security.jwtSecret`）
- Token 格式为 `jwt.MapClaims{sub, iat, exp}`（单用户模式，user_id 默认为 1）

### Swagger / OpenAPI 文档

使用 **GoFrame 内置 Swagger 支持**，无需第三方库或手写 YAML。

#### 访问地址

| 地址 | 说明 |
|------|------|
| `http://localhost:8080/swagger/` | Swagger UI 交互式文档 |
| `http://localhost:8080/api.json` | OpenAPI 3.0 JSON 规范 |

#### 工作原理

1. **自动生成**：GoFrame 从 `api/v1/*/` 中的 `g.Meta` 标签自动提取路由、参数、请求/响应 Schema
2. **配置驱动**：`manifest/config/config.yaml` 中的 `openapiPath` 和 `swaggerPath` 启用文档端点
3. **元数据增强**：`cmd/server/main.go` 的 `enhanceOpenApi()` 函数补充 Info、Servers、JWT SecurityScheme
4. **UI 模板**：使用自定义 SwaggerUI 模板替代 GoFrame 默认的 Redoc

#### 配置位置

```
config.yaml          → openapiPath: "/api.json"  swaggerPath: "/swagger"
cmd/server/main.go   → enhanceOpenApi()  设置 Info/Servers/SecuritySchemes/SwaggerUI 模板
api/v1/*/            → g.Meta 标签定义路由、方法、标签、摘要
```

#### 开发注意事项

- **新增 API 时**：只需在 `api/v1/{module}/{module}.go` 中定义 `g.Meta` 标签的 Req/Res 结构体，Swagger 文档自动更新
- **`dc` 标签**：字段描述（如 `dc:"计价货币"`），会自动映射为 OpenAPI 的 `description`
- **`v` 标签**：校验规则（如 `v:"required|in:daily,weekly"`），会映射为 `required` 和 `enum`
- **`tags` 标签**：接口分组（如 `tags:"资产"`），映射为 Swagger 标签分组
- **`summary` 标签**：接口摘要，显示在 Swagger UI 的接口列表中
- ❌ **禁止手写 swagger.yaml**：所有文档从代码自动生成，保持单一数据源

### 前后端路由对照表

下表列出前端所有 API 调用路径与后端已实现路由。**新增功能时必须查阅此表，确保前后端路径、参数名、请求/响应体一致**。

> **约定**：所有后端路由统一前缀 `/api/v1`，下表省略该前缀。所有响应均被 `middleware/response.go` 包装为 `{ code, message, data, timestamp }` 格式，表中"响应 data"指 `data` 字段的内容。

#### 认证模块（免认证）

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/auth/status` | GET | — | `{ pin_set }` |
| `/auth/setup` | POST | `{ pin }` | `{ token }` |
| `/auth/login` | POST | `{ pin }` | `{ token }` |
| `/auth/change` | POST | `{ current_pin, new_pin }` | `{ success }` |

#### 资产模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/assets/summary` | GET | `?currency=USD` | `{ total_value, currency, by_source, updated_at }` |
| `/assets/details` | GET | `?source_type=&currency=USD` | `{ assets: [...] }` |
| `/assets/history` | GET | `?days=30&currency=USD` | `{ snapshots: [{ date, total_value, currency }] }` |
| `/assets/refresh` | POST | — | `{ message, refreshed_count }` |

#### 交易所模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/exchanges/accounts` | GET | — | `{ accounts: [...] }` |
| `/exchanges/accounts` | POST | `{ exchange_name, api_key, api_secret, passphrase?, label?, note? }` | `{ account }` |
| `/exchanges/accounts/{id}` | GET | — | `{ account }` |
| `/exchanges/accounts/{id}` | PUT | `{ api_key?, api_secret?, passphrase?, label?, note? }` | `{ account }` |
| `/exchanges/accounts/{id}` | DELETE | — | `{}` |
| `/exchanges/accounts/{id}/test` | POST | — | `{ success, message }` |
| `/exchanges/accounts/{id}/balances` | GET | — | `{ balances, total_value }` |
| `/exchanges/accounts/{id}/sync` | POST | — | `{ message }` |

#### 钱包模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/wallets/addresses` | GET | — | `{ addresses: [...] }` |
| `/wallets/addresses` | POST | `{ blockchain, address, label? }` | `{ address }` |
| `/wallets/addresses/{id}` | GET | — | `{ address }` |
| `/wallets/addresses/{id}` | PUT | `{ blockchain?, address?, label? }` | `{ address }` |
| `/wallets/addresses/{id}` | DELETE | — | `{}` |
| `/wallets/addresses/{id}/balances` | GET | — | `{ native_balance, token_balances, total_value_usd }` |
| `/wallets/addresses/{id}/sync` | POST | — | `{ message }` |
| `/wallets/batch` | POST | `{ addresses: [{ blockchain, address, label? }] }` | `{ imported, failed }` |

#### 手动资产模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/manual/assets` | GET | — | `{ assets: [...] }` |
| `/manual/assets` | POST | `{ asset_type, asset_name, amount, currency, notes?, institution? }` | `{ asset }` |
| `/manual/assets/{id}` | PUT | `{ asset_type?, asset_name?, amount?, currency?, notes?, institution? }` | `{ asset }` |
| `/manual/assets/{id}` | DELETE | — | `{}` |

#### 汇率模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/rates/current` | GET | `?currencies=BTC,ETH` | `{ rates, base, last_updated, source, is_cached }` |
| `/rates/prices` | GET | `?symbols=BTC,ETH` | `{ prices: [...] }` |
| `/rates/refresh` | POST | — | `{ message }` |
| `/rates/history` | GET | `?base=USD&quote=BTC&days=30` | `{ history: [...], base, quote, days }` |

#### 通知模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/notifications` | GET | `?page=1&page_size=20` | `{ list: [{ id, type, title, message, is_read, created_at }], pagination }` |
| `/notifications/unread-count` | GET | — | `{ count }` |
| `/notifications/{id}/read` | POST | — | `{}` |
| `/notifications/read-all` | POST | — | `{}` |
| `/notifications/preferences` | GET | — | `{ preferences: { email_enabled, push_enabled, price_alert, portfolio_alert, system_notice } }` |
| `/notifications/preferences` | PUT | `{ email_enabled?, push_enabled?, price_alert?, portfolio_alert?, system_notice? }` | `{ preferences }` |

#### WebPush 模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/notifications/webpush/vapid` | GET | — | `{ vapid_public_key }` |
| `/notifications/webpush/subscribe` | POST | `{ endpoint, keys: { p256dh, auth } }` | `{}` |
| `/notifications/webpush/unsubscribe` | POST | `{ endpoint }` | `{}` |

#### 价格预警模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/alerts` | GET | — | `{ alerts: [...] }` |
| `/alerts` | POST | `{ symbol, condition, target_price, note? }` | `{ alert }` |
| `/alerts/{id}` | PUT | `{ symbol?, condition?, target_price?, note?, is_active? }` | `{ alert }` |
| `/alerts/{id}` | DELETE | — | `{}` |

#### 报告模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/reports` | GET | `?type=&limit=20` | `{ reports: [...] }` |
| `/reports/{id}` | GET | — | `{ report }` |
| `/reports/monthly/{month}` | GET | month 格式 `2024-01` | `{ report }` |
| `/reports/annual/{year}` | GET | year 格式 `2024` | `{ report }` |
| `/reports/generate` | POST | `{ type }` (daily/weekly/monthly/annual) | `{ report }` |
| `/reports/compare` | GET | `?report_id_1=&report_id_2=` | `{ report_1, report_2, value_diff, change_diff }` |

#### DeFi 模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/defi/positions` | GET | `?currency=USD` | `{ positions, total_value, currency }` |
| `/defi/protocols` | GET | — | `{ protocols: [...] }` |
| `/defi/stats` | GET | — | `{ total_value_locked, position_count, by_protocol, by_chain, by_type }` |

#### NFT 模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/nft/assets` | GET | — | `{ assets: [{ ..., estimated_value }], total_value }` |
| `/nfts/collections` | GET | — | `{ collections, total_count, total_value }` |

#### 交易记录模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/transactions` | GET | `?page=1&page_size=20&source=&type=&start=&end=&cursor=` | `{ transactions, total, page, page_size }` |
| `/transactions/sync` | POST | — | `{ message, synced_count }` |
| `/transactions/stats` | GET | — | `{ total_transactions, total_volume, total_fees, by_type, by_source }` |
| `/settings/tx-sync` | GET | — | `{ settings: { auto_sync, sync_interval, last_sync_at } }` |
| `/settings/tx-sync` | PUT | `{ auto_sync?, sync_interval? }` | `{ settings }` |

#### 分析模块（盈亏/归因/预测/费用）

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/analytics/pnl/daily` | GET | `?days=30` | `{ daily: [{ date, pnl, pnl_percent, total_value }] }` |
| `/analytics/pnl/summary` | GET | — | `{ total_pnl, total_pnl_percent, pnl_7d, pnl_30d, pnl_90d, best_day, worst_day, ... }` |
| `/analytics/attribution` | GET | `?days=7&currency=USD` | `{ total_return, total_percent, attributions, days, currency }` |
| `/analytics/forecast` | GET | `?days=30` | `{ forecast_points, trend, confidence, slope, days }` |
| `/analytics/fees` | GET | `?range=30d&currency=USD` | `{ total_fees, trading_fees, gas_fees, currency, breakdown, daily_trend }` |

#### 策略引擎模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/strategies` | GET | — | `{ strategies: [...] }` |
| `/strategies` | POST | `{ name, type, config }` | `{ strategy }` |
| `/strategies/{id}` | PUT | `{ name?, type?, config?, is_active? }` | `{ strategy }` |
| `/strategies/{id}` | DELETE | — | `{}` |
| `/strategies/{id}/analysis` | GET | — | `{ analysis }` |
| `/strategies/{id}/rebalance` | GET | — | `{ analysis }` |

#### 成就系统模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/achievements` | GET | — | `{ achievements, total_count, unlocked_count }` |
| `/achievements/check` | POST | — | `{ newly_unlocked: [...] }` |

#### 基准对比模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/benchmark` | GET | `?range=30d` (7d/30d/90d/1y) | `{ series, range, start_date, end_date }` |

#### 市场数据模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/market/gas` | GET | — | `{ prices: [{ chain, low, standard, fast, instant, base_fee, level }], updated_at }` |

#### 目标追踪模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/goals` | GET | — | `{ goals: [...] }` |
| `/goals` | POST | `{ title, type, target_value, currency?, deadline? }` | `{ goal }` |
| `/goals/{id}` | PUT | `{ title?, type?, target_value?, currency?, deadline? }` | `{ goal }` |
| `/goals/{id}` | DELETE | — | `{}` |

#### 资产健康评分模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/portfolio/health` | GET | `?currency=USD` | `{ overall_score, level, details, currency, updated_at }` |

#### 用户模块

| 路由 | 方法 | 请求参数/请求体 | 响应 data |
|------|------|----------------|-----------|
| `/users/settings` | GET | — | `{ settings: { key: value } }` |
| `/users/settings` | PUT | `{ settings: { key: value } }` | `{ message }` |
| `/users/clear-cache` | POST | — | `{ message }` |
| `/users/reset-settings` | POST | — | `{ message }` |
| `/users/export` | GET | — | `{ exchange_accounts, wallet_addresses, manual_assets, strategies, goals, price_alerts, settings, exported_at }` |

#### 前端已调用但后端未实现

所有前端路由均已对齐，无缺失。

#### 参数对接备忘（v0.2.0 → v0.2.1 修正）

以下是前后端 API 对接过程中修正的 16 处参数/字段不匹配，开发新功能时需注意同类问题：

| 模块 | 修正内容 |
|------|---------|
| auth | `old_pin` → `current_pin` |
| asset | 前端 `timeRange: '30D'` → 后端 `days: int`，响应 `{ snapshots }` 需转为图表 `{ labels, values }` |
| notification | 未读计数 `unread_count` → `count`；偏好设置 5 字段对齐 |
| analytics | `getDailyPnL(range)` / `getForecast(days)` / `getAttribution(range, currency)` 参数重构 |
| report | `generateReport` 改为 POST body；`compareReports` 改为 `report_id_1` + `report_id_2` |
| fee | 路由 `/fees/summary` → `/analytics/fees`；参数 `timeRange` → `range` |
| benchmark | 参数 `period` → `range` |
| nft | 路由 `/nfts` → `/nft/assets`；响应 `{ assets, total_value }` |
| market/gas | 响应 `chains` → `prices`；字段 `name/normal` → `chain/standard` |
| wallet | 批量导入请求体 `{ addresses: [{ blockchain, address, label }] }` |
| user | 更新设置请求体 `{ settings: {...} }` |

---

## 19 张数据库表（DAO 已生成）

| 序号 | 表名 | 说明 |
|------|------|------|
| 1 | exchange_accounts | 交易所账户 |
| 2 | wallet_addresses | 钱包地址 |
| 3 | manual_assets | 手动资产 |
| 4 | asset_snapshots | 资产快照 |
| 5 | asset_details | 资产明细 |
| 6 | exchange_rates | 汇率缓存 |
| 7 | system_configs | 系统配置 |
| 8 | notifications | 通知 |
| 9 | notification_preferences | 通知偏好 |
| 10 | price_alerts | 价格预警 |
| 11 | reports | 报告 |
| 12 | nft_caches | NFT 缓存 |
| 13 | unified_transactions | 统一交易记录 |
| 14 | web_push_subscriptions | WebPush 订阅 |
| 15 | strategies | 策略 |
| 16 | user_achievements | 用户成就 |
| 17 | sync_metadata | 同步元数据 |
| 18 | transaction_daily_summaries | 交易日汇总 |
| 19 | goals | 目标 |

---

## 常用命令

```bash
cd core
go run cmd/server/main.go    # 启动服务（localhost:8080）
go build ./...               # 编译检查
go test ./...                # 运行测试
go fmt ./...                 # 格式化
go vet ./...                 # 静态检查
gf gen dao                   # 生成 DAO/Entity/DO（需配置 hack/config.yaml）
```

---

## 配置

配置文件路径：`manifest/config/config.yaml`

环境变量覆盖（优先级高于配置文件）：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `ALLFI_PORT` | 服务端口 | 8080 |
| `ALLFI_MODE` | 运行模式 | development |
| `ALLFI_DB_TYPE` | 数据库类型 | sqlite |
| `ALLFI_DB_PATH` | SQLite 路径 | data/allfi.db |
| `ALLFI_MASTER_KEY` | AES 加密主密钥 | — |
| `ETHERSCAN_API_KEY` | Etherscan API Key | — |
| `BSCSCAN_API_KEY` | BscScan API Key | — |
| `POLYGONSCAN_API_KEY` | PolygonScan API Key | — |
| `ALCHEMY_API_KEY` | Alchemy NFT API Key | — |
| `COINGECKO_API_KEY` | CoinGecko API Key | — |

---

## Git 提交规则

- 提交信息使用中文，简洁明了
- **禁止**包含 `Co-Authored-By: Claude` 或 `Generated with Claude Code` 等标记
- 每完成一个任务提交一次

---

## 工具函数

### GetUserID — 从上下文获取用户 ID

位于 `internal/consts/user.go`，用于统一获取当前用户 ID：

```go
import "your-finance/allfi/internal/consts"

userID := consts.GetUserID(ctx)
```

- 优先从 GoFrame 请求上下文中获取（由 `middleware.Auth` 写入的 `user_id` 参数）
- 非 HTTP 请求场景（定时任务等）或未认证时，返回默认值 `1`（单用户模式）
- **后续迁移**：新代码应使用 `consts.GetUserID(ctx)` 替代硬编码 `userID := 1`

---

## 遗留代码（已清理）

以下旧目录已于 2026-02-12 清理完毕：

| 已删除目录 | 说明 | 替代方案 |
|-----------|------|---------|
| ~~`internal/api/`~~ | 旧 handlers + router.go | `app/{module}/controller/` |
| ~~`internal/service/`~~ | 旧业务逻辑层 | `app/{module}/logic/` |
| ~~`internal/repository/`~~ | 旧 GORM 数据访问层 | `internal/dao/` |
| ~~`internal/models/`~~ | 旧 GORM 数据模型 | `internal/model/entity/` |
| ~~`internal/database/`~~ | 旧数据库初始化 | GoFrame ORM 自带连接管理 |
| ~~`cmd/migrate/`~~ | 旧 GORM 迁移工具 | GoFrame ORM `gf gen dao` |
| ~~`migrations/`~~ | 旧 GORM 迁移定义 | GoFrame ORM 自动建表 |

**仍保留的目录**：
- `internal/cron/` — 定时任务（已迁移完成，正在使用）
- `internal/integrations/` — 第三方 API 集成（正在使用）

---

## 文档参考

**在线 API 文档**（启动服务后访问）：
- Swagger UI：`http://localhost:8080/swagger/`
- OpenAPI JSON：`http://localhost:8080/api.json`

所有项目文档统一存放在根目录 `docs/` 下：

- `docs/tech/api-reference.md` — API 接口文档
- `docs/tech/tech-baseline.md` — 技术基线
- `docs/specs/backend-spec.md` — 后端需求规格
- `docs/guides/dev-guide.md` — 开发指南
- `docs/guides/deployment-guide.md` — 部署指南
- `docs/README.md` — 文档索引

---

**最后更新**：2026-02-24
