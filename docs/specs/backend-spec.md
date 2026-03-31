# AllFi 后端规格

> 审计日期：2026-03-31
> 范围：当前仓库实现，不再描述历史 GORM / 旧路由方案。

---

## 1. 技术栈

| 类别 | 当前实现 |
|------|---------|
| 语言 | Go 1.24.11 |
| Web 框架 | GoFrame `ghttp` |
| 数据访问 | GoFrame `gdb` + DAO/DO/Entity |
| 配置 | GoFrame `gcfg` |
| OpenAPI | GoFrame `goai` |
| JSON | `json-iterator/go` |
| 认证 | bcrypt + JWT + 2FA/TOTP |
| 加密 | AES-256-GCM |
| 推送 | `webpush-go` |
| 默认数据库 | SQLite |

### 重要更正

- 当前后端不是 `net/http` 自定义路由主架构
- 当前后端不是 GORM 主数据访问层
- 默认文档化数据库路径是 SQLite，MySQL 仅保留零散辅助代码，不是当前主路径

---

## 2. 目录结构

```text
core/
├── api/v1/                  API 定义与 OpenAPI 元数据
├── cmd/server/main.go       服务入口
├── internal/
│   ├── app/                 29 个业务模块
│   ├── cron/                10 个定时任务
│   ├── database/            启动建表与补充迁移
│   ├── integrations/        第三方集成
│   ├── middleware/          CORS / Auth / Logger / ErrorHandler / Response 等
│   ├── statics/             嵌入式前端静态资源
│   ├── utils/               配置、缓存、加密、响应辅助
│   └── version/             构建时版本信息
└── manifest/config/         配置模板与默认配置
```

---

## 3. 业务模块

### 当前模块数

- 29 个 `core/internal/app/*`

### 模块列表

- `achievement`
- `asset`
- `attribution`
- `auth`
- `benchmark`
- `cross_chain`
- `defi`
- `exchange`
- `exchange_rate`
- `fee`
- `forecast`
- `gas`
- `goal`
- `health`
- `health_score`
- `manual_asset`
- `market`
- `nft`
- `notification`
- `pnl`
- `price_alert`
- `report`
- `risk`
- `strategy`
- `system`
- `transaction`
- `user`
- `wallet`
- `webpush`

---

## 4. 路由与接口

### 路由组织

- 所有 API 统一挂载在 `/api/v1`
- `health`、`auth`、`system` 部分接口为免认证
- 其余业务模块通过 `middleware.Auth` 保护

### 当前接口暴露

- `core/api/v1/**/*.go` 中共有 78 个 `g.Meta path`
- Swagger UI：`/swagger/`
- OpenAPI JSON：`/api.json`

### 典型模块分组

| 模块 | 说明 |
|------|------|
| `auth` | 状态、设置密码、登录、改密、2FA |
| `asset` | 汇总、明细、快照 |
| `exchange` / `wallet` / `manual_asset` | 资产来源管理 |
| `defi` / `nft` / `cross_chain` | Web3 扩展数据 |
| `report` / `pnl` / `benchmark` / `risk` | 分析与报告 |
| `system` | 版本、更新检查、更新/回滚状态 |

---

## 5. 中间件

当前中间件集中在 `core/internal/middleware/`，主要包括：

- CORS
- Context 注入
- Logger
- ErrorHandler
- Response 包装
- Auth

说明：

- 中间件链由 `core/cmd/server/main.go` 统一注册
- Swagger 与 `/api.json` 明确绕过 SPA fallback

---

## 6. 数据存储

### 默认路径

- 配置文件：`core/manifest/config/config.yaml`
- 默认连接：`sqlite::@file(./data/allfi.db)`

### 初始化方式

- `core/internal/database/migrate.go` 在启动时自动执行
- 采用手写 DDL + 补充 SQL 迁移文件
- 不再沿用“GORM AutoMigrate”表述

### 数据域示例

- 系统配置
- 用户与认证状态
- 交易所账户
- 钱包地址
- 手动资产
- 资产快照 / 资产明细
- 汇率
- 统一交易记录
- 通知与通知偏好
- 价格预警
- 报告
- 策略与目标
- DeFi 借贷历史
- NFT 缓存
- Cross-chain 交易

---

## 7. 外部集成

### CEX

- Binance
- OKX
- Coinbase

### 链上 / NFT

- Etherscan 系列
- Alchemy
- Bridge/Stargate

### 价格与汇率

- CoinGecko
- Yahoo
- Frankfurter
- GateIO
- Binance
- Local provider

### DeFi 协议

- Lido
- Rocket Pool
- Aave
- Compound
- Uniswap V2
- Uniswap V3
- Curve

---

## 8. 认证与更新

### 认证

- 密码模式：
  - PIN：4-20 位数字
  - 复杂密码：8-20 位，含大小写字母和数字
- 支持 2FA / TOTP
- JWT Bearer Token 用于会话

### 更新

| 场景 | 当前实现 |
|------|---------|
| 宿主机运行 | 下载 release tarball，`go-update` 原位替换 |
| Docker 运行 | 调用 `updater` sidecar，执行 git checkout + docker build + restart |

---

## 9. 测试现状

执行命令：

```bash
cd core && go test ./...
```

当前结果：

- 默认测试当前可通过
- 建议命令：

```bash
cd core && go test ./... -timeout 60s
```

- `exchange_rate/provider` 中依赖公网的集成测试默认跳过，需设置 `ALLFI_RUN_ONLINE_TESTS=1` 才执行
