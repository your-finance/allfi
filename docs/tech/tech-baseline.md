# AllFi 技术基线文档

> 版本：v2.0 | 更新时间：2026-02-11 | 状态：与代码实现对齐

---

## 1. 技术架构概览

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户浏览器                                 │
│              (PWA / 响应式 / 移动端适配)                           │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP(S)
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      前端层 (Frontend)                            │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Vue 3 + Vite 7 + Pinia 3                                │   │
│  │  - 9 个页面 (SPA, hash 路由)                               │   │
│  │  - 39 个组件                                               │   │
│  │  - 12 个 Pinia Store                                       │   │
│  │  - Chart.js 4 数据可视化                                    │   │
│  │  - Tailwind CSS 4 + 4 套主题                               │   │
│  │  - 3 语言 i18n (zh-CN / zh-TW / en-US)                    │   │
│  └──────────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────────┘
                         │ RESTful API (JSON)
                         │ 75 条路由
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      后端层 (Backend)                             │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Go 1.24 + net/http (标准库路由)                           │   │
│  │  ┌────────────────────────────────────────────────────┐   │   │
│  │  │  中间件链: Recovery → CORS → Logger → Auth(JWT)     │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  │  ┌────────────────────────────────────────────────────┐   │   │
│  │  │  Handler 层 (24 个处理器)                            │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  │  ┌────────────────────────────────────────────────────┐   │   │
│  │  │  Service 层 (20+ 服务，接口驱动)                     │   │   │
│  │  │  - ServiceContainer 依赖注入                        │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  │  ┌────────────────────────────────────────────────────┐   │   │
│  │  │  Repository 层 (18 个 Repo，GORM v2)               │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  │  ┌────────────────────────────────────────────────────┐   │   │
│  │  │  定时任务 (6 个 CronJob)                            │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    数据存储层 (Storage)                            │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │  SQLite3          │  │  MySQL           │                     │
│  │  (默认，单文件)    │  │  (可选，生产环境)  │                     │
│  └──────────────────┘  └──────────────────┘                     │
│  18 张业务表 · GORM AutoMigrate · 软删除                          │
└─────────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    外部服务层 (External APIs)                      │
│  ┌──────────────┐  ┌───────────────┐  ┌───────────────────┐    │
│  │  CEX          │  │  区块链        │  │  价格数据          │    │
│  │  - Binance    │  │  - Etherscan  │  │  - CoinGecko      │    │
│  │  - OKX (CCXT) │  │    (6 条链)   │  │  - Yahoo Finance  │    │
│  │  - Coinbase   │  │  - Alchemy    │  │                    │    │
│  │    (CCXT)     │  │    (NFT)      │  │                    │    │
│  └──────────────┘  └───────────────┘  └───────────────────┘    │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  DeFi 协议 (7 个)                                          │  │
│  │  Lido · Rocket Pool · Aave · Compound                     │  │
│  │  Uniswap V2 · Uniswap V3 · Curve                         │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. 技术栈

### 2.1 后端

| 技术 | 版本/说明 | 选型理由 |
|------|----------|---------|
| Go | 1.24 | 高性能、并发友好、静态编译 |
| `net/http` | 标准库 | Go 1.22+ 增强路由，无需框架 |
| GORM | v2 | 成熟 ORM，支持 SQLite/MySQL |
| SQLite3 | 默认 DB | 零配置、单文件、适合自托管 |
| MySQL | 可选 DB | 大数据量或多实例场景 |
| GoFrame `gcfg` | 配置管理 | YAML 配置 + 环境变量覆盖 |
| json-iterator | JSON 序列化 | 高性能 JSON 编解码 |
| shopspring/decimal | 精确计算 | 金融数据精度要求 |
| golang-jwt/v5 | JWT 认证 | 行业标准 Token 方案 |
| golang.org/x/crypto | bcrypt | PIN 码安全哈希 |
| Prometheus | 监控 | 标准化指标暴露 |
| webpush-go | 浏览器推送 | Web Push 通知 |
| go-binance/v2 | Binance SDK | 专用客户端 |
| ccxt/go/v4 | 多交易所 SDK | OKX/Coinbase 统一接口 |
| stretchr/testify | 测试框架 | 断言和 Mock |

### 2.2 前端

| 技术 | 版本/说明 | 选型理由 |
|------|----------|---------|
| Vue 3 | Composition API + `<script setup>` | 响应式框架 |
| Vite | 7 | 极速 HMR，原生 ESM |
| Pinia | 3 | Vue 官方状态管理 |
| Tailwind CSS | 4 | 原子化 CSS，快速定制 |
| Chart.js | 4 | 轻量可视化 |
| Phosphor Icons | — | 图标库 |
| native Fetch API | — | 无需 Axios |

### 2.3 部署

| 技术 | 用途 |
|------|------|
| Docker + Docker Compose | 容器化部署 |
| Nginx | 反向代理 + 静态文件服务（可选） |
| Systemd | Linux 服务管理 |
| Makefile | 构建自动化 |

---

## 3. 项目结构

```
allfi/
├── core/                               # 后端 (Go)
│   ├── cmd/server/main.go              # 服务入口
│   ├── config/
│   │   ├── config.yaml                 # 配置文件
│   │   └── config.example.yaml         # 配置模板
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go              # 路由定义 (75 条)
│   │   │   ├── handlers/ (24 个)      # HTTP 处理器
│   │   │   └── middleware/ (4 个)      # CORS/Logger/Recovery/Auth
│   │   ├── service/ (20+ 个)          # 业务服务（接口驱动）
│   │   │   ├── interfaces.go          # 所有服务接口定义
│   │   │   └── container.go           # 依赖注入容器
│   │   ├── repository/ (18 个)        # 数据访问层
│   │   ├── models/ (18 个表)          # GORM 数据模型
│   │   ├── integrations/ (8 个模块)   # 第三方集成
│   │   │   ├── binance/               # go-binance SDK
│   │   │   ├── okx/                   # CCXT
│   │   │   ├── coinbase/              # CCXT
│   │   │   ├── etherscan/             # 6 链 Etherscan 兼容
│   │   │   ├── coingecko/             # 加密货币价格
│   │   │   ├── yahoo/                 # 法币汇率
│   │   │   ├── alchemy/               # NFT 数据
│   │   │   └── defi/ (7 个协议)       # DeFi 仓位查询
│   │   ├── cron/ (6 个任务)           # 定时任务
│   │   ├── model/                     # DTO 传输对象
│   │   └── utils/                     # 工具库
│   ├── docs/                          # Swagger 文档
│   └── data/allfi.db                  # SQLite 数据库
├── webapp/                             # 前端 (Vue 3)
│   ├── src/
│   │   ├── pages/ (9 个)              # 页面组件
│   │   ├── components/ (39 个)        # UI 组件
│   │   ├── stores/ (12 个)            # Pinia Store
│   │   ├── api/ (16 个模块)           # API 服务
│   │   ├── composables/               # Vue Composables
│   │   ├── i18n/                      # 多语言 (800+ key)
│   │   └── data/                      # Mock 数据
│   ├── public/
│   └── package.json
├── docs/                               # 项目文档
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

---

## 4. 数据库设计

### 4.1 概述

- ORM: GORM v2，使用 `AutoMigrate` 自动建表
- 默认: SQLite3（`data/allfi.db`），可切换 MySQL
- 所有表继承 `BaseModel`（id, created_at, updated_at, deleted_at 软删除）
- 金额字段使用 `decimal.Decimal`（shopspring/decimal）

### 4.2 数据表（18 张）

| # | 表名 | 说明 | 核心字段 |
|---|------|------|---------|
| 1 | `exchange_accounts` | 交易所账户 | exchange, api_key_encrypted, api_secret_encrypted, status |
| 2 | `wallet_addresses` | 钱包地址 | blockchain(eth/bsc/polygon), address, label |
| 3 | `manual_assets` | 传统资产 | asset_type(bank/cash/stock/fund), amount(Decimal), currency |
| 4 | `asset_details` | 资产明细快照 | source_type, symbol, balance, price_usd, value_usd |
| 5 | `asset_snapshots` | 资产总值快照 | total_value_usd, exchange_value, blockchain_value, manual_value |
| 6 | `exchange_rates` | 汇率缓存 | from_currency, to_currency, rate, source, expires_at |
| 7 | `system_config` | 系统配置 KV | config_key, config_value |
| 8 | `notifications` | 通知消息 | type, title, content, is_read |
| 9 | `price_alerts` | 价格预警 | symbol, condition, target_price, is_active |
| 10 | `reports` | 资产报告 | report_type(daily/weekly/monthly/annual), content |
| 11 | `unified_transactions` | 统一交易记录 | tx_type(buy/sell/swap/transfer), source, from_asset, to_asset |
| 12 | `transaction_daily_summaries` | 交易日汇总 | date, buy_count, sell_count, total_fee_usd |
| 13 | `strategies` | 投资策略 | name, target_allocations, rebalance_threshold |
| 14 | `achievements` | 成就记录 | achievement_id, unlocked_at |
| 15 | `nfts` | NFT 资产 | chain, contract_address, token_id, collection |
| 16 | `goals` | 目标追踪 | name, target_amount, deadline |
| 17 | `sync_metadata` | 同步元数据 | source, last_sync_at, sync_status |

### 4.3 支持的常量

| 类型 | 值 |
|------|-----|
| 交易所 | binance, okx, coinbase |
| 区块链 | ethereum, bsc, polygon |
| 资产来源 | cex, blockchain, manual |
| 传统资产类型 | bank, cash, stock, fund |
| 计价货币 | USDC, USD, BTC, ETH, CNY |
| 交易类型 | buy, sell, swap, transfer, deposit, withdraw |
| 汇率数据源 | yahoo_finance, coingecko |

---

## 5. API 设计

### 5.1 规范

- **Base URL**: `http://localhost:8080/api/v1`
- **路由风格**: Go 1.22+ 增强路由（`METHOD /path/{param}`）
- **响应格式**: 统一 JSON（code + message + data + timestamp）
- **认证**: PIN 码 + JWT Bearer Token
- **文档**: Swagger UI 内置（`/api/v1/docs`）

### 5.2 错误码体系

| 范围 | 类别 | 码值 |
|------|------|------|
| 0 | 成功 | 0 |
| 1001-1999 | 客户端错误 | 1001 参数错误, 1002 验证失败, 1003 资源不存在, 1004 重复条目, 1005 未授权 |
| 2001-2999 | 服务器错误 | 2001 内部错误, 2002 数据库错误, 2003 外部 API 错误, 2004 加密错误 |
| 3001-3999 | 业务错误 | 3001 交易所 API 错误, 3004 地址无效, 3006 汇率获取失败, 3007 快照失败 |

### 5.3 接口统计（75 条路由）

| 模块 | 数量 | 说明 |
|------|------|------|
| 认证 | 4 | PIN 设置/登录/修改/状态 |
| 资产 | 8 | 总览/详情/历史/刷新 + 传统资产 CRUD |
| 交易所 | 7 | CRUD + 测试连接 + 余额 |
| 钱包 | 8 | CRUD + 余额 + 批量导入 |
| 汇率 | 3 | 当前/价格/刷新 |
| 通知 | 9 | 列表/已读/偏好 + WebPush |
| 价格预警 | 4 | CRUD |
| 报告 | 3 | 列表/详情/生成 |
| 交易记录 | 5 | 列表/同步/统计/设置 |
| 策略引擎 | 5 | CRUD + 分析 |
| DeFi | 2 | 仓位 + 协议列表 |
| NFT | 1 | 资产列表 |
| 成就 | 2 | 列表 + 检查 |
| 目标 | 4 | CRUD |
| 分析 | 5 | 费用/盈亏/归因/预测 |
| 其他 | 5 | 健康检查/基准/Gas/健康评分/文档 |

详细路由表见 [后端需求规格文档](../specs/backend-spec.md#43-完整路由表75-条)。

---

## 6. 认证与安全

### 6.1 认证方式

单用户 PIN 码（4-8 位数字）+ JWT Token：

```
首次使用 → 所有接口开放 → 用户设置 PIN
       → PIN bcrypt 哈希存储 → 返回 JWT Token
       → 后续请求需 Bearer Token → 中间件校验
```

白名单路由（无需认证）：`/api/v1/health`、`/api/v1/auth/*`

### 6.2 API Key 加密

- 算法: AES-256-GCM
- 主密钥: 32 字节，环境变量 `ALLFI_MASTER_KEY` 或配置文件
- 存储: Base64(nonce + ciphertext)
- 工具: `utils.EncryptAPIKey()` / `utils.DecryptAPIKey()`

### 6.3 中间件链

```
Recovery → CORS → Logger → Auth(JWT) → Handler
```

---

## 7. 第三方集成

### 7.1 交易所

| 交易所 | SDK | 说明 |
|--------|-----|------|
| Binance | `go-binance/v2` | 专用 SDK |
| OKX | `ccxt/go/v4` | CCXT 统一接口 |
| Coinbase | `ccxt/go/v4` | CCXT 统一接口 |

### 7.2 区块链

| 服务 | 说明 |
|------|------|
| Etherscan | 统一客户端，6 条 EVM 链（Ethereum/BSC/Polygon/Arbitrum/Optimism/Base） |
| Alchemy | NFT 数据（列表/收藏集/估值） |

### 7.3 DeFi 协议（7 个）

| 协议 | 类型 |
|------|------|
| Lido | ETH 质押 |
| Rocket Pool | ETH 质押 |
| Aave | 借贷 |
| Compound | 借贷 |
| Uniswap V2 | DEX 流动性 |
| Uniswap V3 | DEX 流动性 |
| Curve | 稳定币 DEX |

通过 `defi/registry.go` 注册表模式统一管理，所有协议实现同一接口。

### 7.4 价格数据

| 数据源 | 用途 |
|--------|------|
| CoinGecko | 加密货币实时价格 |
| Yahoo Finance | 法币汇率、传统资产价格 |

---

## 8. 定时任务

由 `CronManager` 统一管理 6 个定时任务：

| 任务 | 默认间隔 | 说明 |
|------|---------|------|
| 资产快照 | 1 小时 | 刷新价格 → 创建快照 → 清理旧数据（90 天） |
| 通知摘要 | 每日 | 为启用摘要的用户生成每日通知 |
| 价格预警 | 周期性 | 检查活跃预警条件 |
| 报告生成 | 每日 | 自动生成日报/周报/月报 |
| 策略监控 | 周期性 | 检查策略偏移 |
| 风险预警 | 周期性 | 检查资产集中度、波动率 |

---

## 9. 配置管理

使用 GoFrame `gcfg` 封装，支持配置文件 + 环境变量覆盖。

### 9.1 核心配置项

| 配置项 | 环境变量 | 默认值 |
|--------|---------|--------|
| `server.port` | `ALLFI_PORT` | 8080 |
| `server.mode` | `ALLFI_MODE` | development |
| `database.type` | `ALLFI_DB_TYPE` | sqlite |
| `database.sqlite.path` | `ALLFI_DB_PATH` | ../data/allfi.db |
| `security.master_key` | `ALLFI_MASTER_KEY` | —（必填） |
| `external_apis.etherscan.api_key` | `ETHERSCAN_API_KEY` | — |
| `external_apis.coingecko.api_key` | `COINGECKO_API_KEY` | — |
| `cron.snapshot_interval` | — | 3600（秒） |
| `cron.price_cache_ttl` | — | 300（秒） |
| `defaults.currency` | — | USDC |
| `defaults.history_retention_days` | — | 90 |

---

## 10. 开发规范

### 10.1 Go 代码规范

- 包名: 小写单词（`service`, `models`）
- 文件名: 蛇形命名（`asset_service.go`）
- 结构体: 大驼峰（`AssetService`）
- 函数: 大驼峰（导出）/ 小驼峰（内部）
- 所有注释使用中文

### 10.2 数据库规范

- 表名: 蛇形、复数（`exchange_accounts`）
- 字段: 蛇形（`api_key_encrypted`）
- 索引: `idx_表名_字段名`

### 10.3 API 规范

- 路径: 小写、复数名词（`/assets`, `/exchanges/accounts`）
- 参数: 蛇形（`?currency=USDC&page_size=20`）

### 10.4 测试

- 测试框架: `stretchr/testify`
- 接口测试: Mock 服务实现接口
- 文件命名: `*_test.go`

---

## 11. 部署方案

### 11.1 Docker 部署

```bash
docker-compose up -d
# 访问 http://localhost:8080
```

### 11.2 手动部署

```bash
# 后端
cd core && go build -o allfi ./cmd/server && ./allfi

# 前端
cd webapp && npm install && npm run build
# 构建产物在 dist/，由后端静态服务或 Nginx 提供
```

### 11.3 数据备份

```bash
# SQLite 备份
sqlite3 data/allfi.db ".backup data/backup/allfi_$(date +%Y%m%d).db"
```

---

## 12. 技术决策记录 (ADR)

### ADR-001: SQLite 作为默认数据库

- 零配置、单文件、便于备份
- 单用户场景性能充足
- 可切换 MySQL（修改配置即可）

### ADR-002: 标准库 net/http 替代 Web 框架

- Go 1.22+ 增强路由功能足够（支持 `METHOD /path/{param}`）
- 无外部框架依赖，减少升级风险
- GoFrame 仅用于配置管理（`gcfg`）

### ADR-003: AES-256-GCM 加密 API Key

- 行业标准加密算法
- 自托管环境用户完全控制主密钥
- 数据库泄露不影响 API Key 安全

### ADR-004: CCXT 统一交易所集成

- OKX、Coinbase 通过 CCXT Go SDK 接入
- 统一接口降低维护成本
- Binance 保留专用 SDK（功能更完整）

### ADR-005: 接口驱动的服务层

- 所有服务定义接口（`*Interface`）
- `ServiceContainer` 管理依赖注入
- 便于 Mock 测试和替换实现

---

## 13. 数据统计总览

| 维度 | 数量 |
|------|------|
| 后端 API 路由 | 75 条 |
| 后端 Handler | 24 个 |
| 后端服务接口 | 20+ 个 |
| 后端 Repository | 18 个 |
| 数据库表 | 18 张 |
| 第三方集成 | 8 个模块 |
| DeFi 协议 | 7 个 |
| 定时任务 | 6 个 |
| 前端页面 | 9 个 |
| 前端组件 | 39 个 |
| Pinia Store | 12 个 |
| API 服务模块 | 16 个 |
| i18n 翻译 key | 800+ |
| 主题 | 4 套 |
| 支持语言 | 3 种 |

---

## 相关文档

- [后端需求规格](../specs/backend-spec.md) — API 路由表、数据模型、服务接口详情
- [前端需求规格](../specs/frontend-spec.md) — 页面、组件、Store、路由详情
- [API 接口文档](./api-reference.md) — 完整接口请求/响应示例
- [UI/UX 设计规范](../design/ui-ux-standards.md) — 设计系统
- [部署指南](../guides/deployment-guide.md) — Docker/手动部署

---

**文档维护者**: @allfi
**最后更新**: 2026-02-11
