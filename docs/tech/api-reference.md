# AllFi API 接口文档

> **版本**：v3.0
> **更新时间**：2026-02-12
> **API 版本**：v1
> **Base URL**: `http://localhost:8080/api/v1`
> **Swagger UI**: `http://localhost:8080/api/v1/docs`

---

## 目录

1. [接口概述](#1-接口概述)
2. [认证与安全](#2-认证与安全)
3. [响应格式](#3-响应格式)
4. [错误码](#4-错误码)
5. [认证接口](#5-认证接口)
6. [资产接口](#6-资产接口)
7. [手动资产接口](#7-手动资产接口)
8. [交易所账户接口](#8-交易所账户接口)
9. [钱包地址接口](#9-钱包地址接口)
10. [汇率接口](#10-汇率接口)
11. [通知接口](#11-通知接口)
12. [WebPush 推送接口](#12-webpush-推送接口)
13. [价格预警接口](#13-价格预警接口)
14. [报告接口](#14-报告接口)
15. [DeFi 接口](#15-defi-接口)
16. [NFT 接口](#16-nft-接口)
17. [交易记录接口](#17-交易记录接口)
18. [费用分析接口](#18-费用分析接口)
19. [分析接口（盈亏/归因/预测）](#19-分析接口盈亏归因预测)
20. [策略引擎接口](#20-策略引擎接口)
21. [成就系统接口](#21-成就系统接口)
22. [基准对比接口](#22-基准对比接口)
23. [目标追踪接口](#23-目标追踪接口)
24. [资产健康评分接口](#24-资产健康评分接口)
25. [市场数据接口](#25-市场数据接口)
26. [用户设置接口](#26-用户设置接口)
27. [系统接口](#27-系统接口)

---

## 1. 接口概述

### 1.1 基础信息

| 项目 | 说明 |
|------|------|
| **协议** | HTTP/HTTPS |
| **请求方式** | GET, POST, PUT, DELETE |
| **Base URL** | `http://localhost:8080/api/v1` |
| **Content-Type** | `application/json` |
| **字符编码** | UTF-8 |
| **后端框架** | GoFrame v2 |

### 1.2 版本说明

当前 API 版本：**v1**

所有接口路径以 `/api/v1` 开头。后端同时提供 Swagger UI，访问 `/api/v1/docs` 查看交互式文档。

### 1.3 接口统计

| 模块 | 端点数 | 说明 |
|------|--------|------|
| 认证 | 4 | PIN 模式认证 |
| 资产 | 4 | 汇总、详情、历史、刷新 |
| 手动资产 | 4 | CRUD |
| 交易所账户 | 8 | CRUD + 测试 + 余额 + 同步 |
| 钱包地址 | 8 | CRUD + 余额 + 批量导入 + 同步 |
| 汇率 | 4 | 当前 + 价格 + 历史 + 刷新 |
| 通知 | 6 | 列表 + 已读 + 偏好 |
| WebPush | 3 | VAPID + 订阅/退订 |
| 价格预警 | 4 | CRUD |
| 报告 | 6 | 列表/详情/生成/月度/年度/对比 |
| DeFi | 3 | 仓位 + 协议 + 统计 |
| NFT | 2 | 资产列表 + 收藏集 |
| 交易记录 | 5 | 列表 + 同步 + 统计 + 设置 |
| 费用分析 | 1 | 费用统计 |
| 分析 | 4 | 盈亏/摘要/归因/预测 |
| 策略引擎 | 6 | CRUD + 分析 + 再平衡 |
| 成就系统 | 2 | 列表 + 检查 |
| 基准对比 | 1 | 收益率对比 |
| 目标追踪 | 4 | CRUD |
| 资产健康 | 1 | 健康评分 |
| 市场数据 | 1 | Gas 费 |
| 用户设置 | 5 | 获取/更新/导出/清缓存/重置 |
| 系统 | 1 | 健康检查 |
| **合计** | **~87** | |

---

## 2. 认证与安全

### 2.1 认证方式

AllFi 采用 **PIN 码 + JWT** 的轻量认证模式，适合自托管单用户场景。

**流程**：
1. 客户端调用 `GET /auth/status` 检查是否已设置 PIN
2. 首次使用调用 `POST /auth/setup` 设置 PIN，返回 JWT token
3. 后续登录调用 `POST /auth/login`，返回 JWT token
4. 后续请求在 Header 中携带 token：

```http
Authorization: Bearer <jwt_token>
```

### 2.2 免认证路由

以下路由不需要 JWT 认证：

- `GET /api/v1/health` — 健康检查
- `GET /api/v1/auth/status` — 认证状态
- `POST /api/v1/auth/setup` — 首次设置 PIN
- `POST /api/v1/auth/login` — PIN 登录

### 2.3 安全建议

- 仅在受信任的网络环境中使用
- 建议使用 HTTPS（生产环境）
- 定期备份数据库文件
- 不要将 API 暴露到公网

---

## 3. 响应格式

### 3.1 标准响应格式

所有接口由 GoFrame 中间件自动包装为统一的 JSON 格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1707382800
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码：0=成功，非0=错误 |
| message | string | 消息说明 |
| data | object/null | 响应数据 |
| timestamp | int64 | Unix 时间戳 |

---

## 4. 错误码

### 客户端错误（1001-1999）

| 错误码 | 常量名 | 说明 | HTTP 状态码 |
|--------|--------|------|-------------|
| **0** | `CodeSuccess` | 成功 | 200 |
| **1001** | `CodeInvalidParams` | 请求参数错误 | 400 |
| **1002** | `CodeValidationFailed` | 数据验证失败 | 422 |
| **1003** | `CodeResourceNotFound` | 资源不存在 | 404 |
| **1004** | `CodeDuplicateEntry` | 数据已存在 | 409 |
| **1005** | `CodeUnauthorized` | 未授权访问 | 401 |

### 服务器错误（2001-2999）

| 错误码 | 常量名 | 说明 | HTTP 状态码 |
|--------|--------|------|-------------|
| **2001** | `CodeInternalError` | 服务器内部错误 | 500 |
| **2002** | `CodeDatabaseError` | 数据库操作失败 | 500 |
| **2003** | `CodeExternalAPIError` | 外部服务调用失败 | 502 |
| **2004** | `CodeEncryptionError` | 数据加密/解密失败 | 500 |
| **2005** | `CodeCacheError` | 缓存操作失败 | 500 |
| **2006** | `CodeConfigError` | 配置加载失败 | 500 |

### 业务错误（3001-3999）

| 错误码 | 常量名 | 说明 | HTTP 状态码 |
|--------|--------|------|-------------|
| **3001** | `CodeExchangeAPIError` | 交易所 API 调用失败 | 502 |
| **3002** | `CodeBlockchainAPIError` | 区块链 API 调用失败 | 502 |
| **3003** | `CodeInvalidAPIKey` | API Key 无效 | 400 |
| **3004** | `CodeInvalidAddress` | 钱包地址格式无效 | 400 |
| **3005** | `CodeBalanceFetchFailed` | 获取余额失败 | 502 |
| **3006** | `CodeRateFetchFailed` | 获取汇率失败 | 502 |
| **3007** | `CodeSnapshotFailed` | 创建资产快照失败 | 500 |

---

## 5. 认证接口

### 5.1 获取认证状态

```
GET /api/v1/auth/status
```

**功能**：检查是否已设置 PIN 码

**响应**：
```json
{
  "code": 0,
  "data": {
    "pin_set": false
  }
}
```

### 5.2 首次设置 PIN

```
POST /api/v1/auth/setup
```

**请求体**：
```json
{
  "pin": "123456"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pin | string | 是 | 4-20 位 PIN 码 |

**响应**：
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 5.3 PIN 登录

```
POST /api/v1/auth/login
```

**请求体**：
```json
{
  "pin": "123456"
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 5.4 修改 PIN

```
POST /api/v1/auth/change
```

**请求体**：
```json
{
  "current_pin": "123456",
  "new_pin": "654321"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| current_pin | string | 是 | 当前 PIN 码 |
| new_pin | string | 是 | 新 PIN 码（4-20 位） |

**响应**：
```json
{
  "code": 0,
  "data": {
    "success": true
  }
}
```

---

## 6. 资产接口

### 6.1 资产总览

```
GET /api/v1/assets/summary?currency=USD
```

**功能**：获取所有资产的汇总信息，按来源分类

**请求参数**：

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| currency | string | 否 | 计价币种（USD/BTC/ETH/CNY） | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "total_value": 125430.78,
    "currency": "USD",
    "by_source": {
      "cex": 80200,
      "blockchain": 35230,
      "manual": 10000
    },
    "updated_at": "2026-02-10T10:30:00Z"
  }
}
```

### 6.2 资产详情列表

```
GET /api/v1/assets/details?source_type=cex&currency=USD
```

**功能**：获取扁平化的资产详情列表（按持仓逐条列出）

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| source_type | string | 否 | 来源类型筛选（cex/blockchain/manual） | 全部 |
| currency | string | 否 | 计价币种 | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "assets": [
      {
        "id": 1,
        "symbol": "BTC",
        "amount": 0.5,
        "value": 32000,
        "price": 64000,
        "source": "binance",
        "source_type": "cex",
        "updated_at": "2026-02-10T10:30:00Z"
      }
    ]
  }
}
```

### 6.3 资产历史

```
GET /api/v1/assets/history?days=30&currency=USD
```

**功能**：获取历史资产总值曲线数据

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| days | int | 否 | 查询天数 | 30 |
| currency | string | 否 | 计价币种 | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "snapshots": [
      {
        "date": "2026-01-12",
        "total_value": 120000,
        "currency": "USD"
      }
    ]
  }
}
```

### 6.4 刷新资产

```
POST /api/v1/assets/refresh
```

**功能**：手动触发全量资产数据刷新（CEX + 链上 + 手动）

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "资产刷新成功",
    "refreshed_count": 15
  }
}
```

---

## 7. 手动资产接口

> **路径前缀**：`/api/v1/manual/assets`

### 7.1 获取手动资产列表

```
GET /api/v1/manual/assets
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "assets": [
      {
        "id": 1,
        "asset_type": "bank",
        "asset_name": "银行存款",
        "amount": 70000,
        "currency": "CNY",
        "notes": "工商银行定期",
        "institution": "工商银行",
        "created_at": "2026-01-15T00:00:00Z",
        "updated_at": "2026-02-10T00:00:00Z"
      }
    ]
  }
}
```

### 7.2 添加手动资产

```
POST /api/v1/manual/assets
```

**请求体**：
```json
{
  "asset_type": "bank",
  "asset_name": "银行存款",
  "amount": 70000,
  "currency": "CNY",
  "notes": "工商银行定期",
  "institution": "工商银行"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| asset_type | string | 是 | 类型（cash/bank/stock/fund） |
| asset_name | string | 是 | 资产名称（最大 100 字符） |
| amount | float64 | 是 | 数量（最小 0） |
| currency | string | 是 | 币种（CNY/USD/EUR/HKD 等） |
| notes | string | 否 | 备注 |
| institution | string | 否 | 机构名称 |

### 7.3 更新手动资产

```
PUT /api/v1/manual/assets/{id}
```

**请求体**：所有字段可选，仅传需要更新的字段。

### 7.4 删除手动资产

```
DELETE /api/v1/manual/assets/{id}
```

---

## 8. 交易所账户接口

### 8.1 获取账户列表

```
GET /api/v1/exchanges/accounts
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "accounts": [
      {
        "id": 1,
        "exchange_name": "binance",
        "label": "主账户",
        "note": "",
        "status": "active",
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-02-10T10:30:00Z"
      }
    ]
  }
}
```

### 8.2 获取单个账户

```
GET /api/v1/exchanges/accounts/{id}
```

### 8.3 添加账户

```
POST /api/v1/exchanges/accounts
```

**请求体**：
```json
{
  "exchange_name": "binance",
  "api_key": "YOUR_API_KEY",
  "api_secret": "YOUR_API_SECRET",
  "passphrase": "",
  "label": "Trading Account",
  "note": ""
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| exchange_name | string | 是 | 交易所（binance/okx/coinbase） |
| api_key | string | 是 | API Key |
| api_secret | string | 是 | API Secret |
| passphrase | string | 否 | API Passphrase（OKX 必填） |
| label | string | 否 | 账户标签 |
| note | string | 否 | 备注 |

### 8.4 更新账户

```
PUT /api/v1/exchanges/accounts/{id}
```

### 8.5 删除账户

```
DELETE /api/v1/exchanges/accounts/{id}
```

### 8.6 测试连接

```
POST /api/v1/exchanges/accounts/{id}/test
```

**功能**：测试 API Key 是否有效

**响应**：
```json
{
  "code": 0,
  "data": {
    "success": true,
    "message": "连接成功"
  }
}
```

### 8.7 获取账户余额

```
GET /api/v1/exchanges/accounts/{id}/balances
```

**功能**：获取指定账户的所有持仓余额

**响应**：
```json
{
  "code": 0,
  "data": {
    "balances": [
      { "symbol": "BTC", "free": 0.5, "locked": 0.0, "total": 0.5, "value_usd": 32000 },
      { "symbol": "ETH", "free": 8.0, "locked": 0.0, "total": 8.0, "value_usd": 18234 }
    ],
    "total_value": 50234
  }
}
```

### 8.8 同步账户资产

```
POST /api/v1/exchanges/accounts/{id}/sync
```

**功能**：触发交易所账户资产同步刷新

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "同步成功"
  }
}
```

---

## 9. 钱包地址接口

### 9.1 获取钱包列表

```
GET /api/v1/wallets/addresses
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "addresses": [
      {
        "id": 1,
        "blockchain": "ethereum",
        "address": "0x1234...5678",
        "label": "主钱包",
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-02-10T00:00:00Z"
      }
    ]
  }
}
```

### 9.2 获取单个钱包

```
GET /api/v1/wallets/addresses/{id}
```

### 9.3 添加钱包

```
POST /api/v1/wallets/addresses
```

**请求体**：
```json
{
  "blockchain": "ethereum",
  "address": "0x1234567890abcdef1234567890abcdef12345678",
  "label": "主钱包"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| blockchain | string | 是 | 区块链（ethereum/bsc/polygon/arbitrum/optimism/base） |
| address | string | 是 | 钱包地址（0x 开头） |
| label | string | 否 | 地址标签 |

### 9.4 更新钱包

```
PUT /api/v1/wallets/addresses/{id}
```

### 9.5 删除钱包

```
DELETE /api/v1/wallets/addresses/{id}
```

### 9.6 获取钱包余额

```
GET /api/v1/wallets/addresses/{id}/balances
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "native_balance": 5.0,
    "token_balances": {
      "USDC": 5000,
      "UNI": 100
    },
    "total_value_usd": 16400
  }
}
```

### 9.7 批量导入钱包

```
POST /api/v1/wallets/batch
```

**请求体**：
```json
{
  "addresses": [
    { "blockchain": "ethereum", "address": "0xabc123...", "label": "钱包1" },
    { "blockchain": "ethereum", "address": "0xdef456...", "label": "钱包2" }
  ]
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "imported": 2,
    "failed": 0
  }
}
```

### 9.8 同步钱包资产

```
POST /api/v1/wallets/addresses/{id}/sync
```

**功能**：触发钱包地址资产同步刷新

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "同步成功"
  }
}
```

---

## 10. 汇率接口

### 10.1 获取当前汇率

```
GET /api/v1/rates/current?currencies=BTC,ETH,USD
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| currencies | string | 否 | 币种列表（逗号分隔） |

**响应**：
```json
{
  "code": 0,
  "data": {
    "rates": {
      "BTC": 64000,
      "ETH": 2280,
      "CNY": 0.14
    },
    "base": "USD",
    "last_updated": 1707580800000,
    "source": "coingecko",
    "is_cached": true
  }
}
```

### 10.2 获取币种价格

```
GET /api/v1/rates/prices?symbols=BTC,ETH
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| symbols | string | 否 | 币种列表（逗号分隔） |

**响应**：
```json
{
  "code": 0,
  "data": {
    "prices": [
      {
        "symbol": "BTC",
        "price_usd": 64000,
        "change_24h": 3.2,
        "last_updated": 1707580800
      }
    ]
  }
}
```

### 10.3 获取历史汇率

```
GET /api/v1/rates/history?base=USD&quote=CNY&days=30
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| base | string | 否 | 基准货币 | USD |
| quote | string | 否 | 目标货币 | — |
| days | int | 否 | 天数 | 30 |

**响应**：
```json
{
  "code": 0,
  "data": {
    "history": [
      { "date": "2026-01-13", "rate": 7.25, "base": "USD", "quote": "CNY" }
    ],
    "base": "USD",
    "quote": "CNY",
    "days": 30
  }
}
```

### 10.4 刷新汇率

```
POST /api/v1/rates/refresh
```

**功能**：手动触发汇率数据刷新

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "汇率刷新成功"
  }
}
```

---

## 11. 通知接口

### 11.1 获取通知列表

```
GET /api/v1/notifications?page=1&page_size=20
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页数量 | 20 |

**响应**：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "type": "daily_digest",
        "title": "每日资产摘要",
        "message": "总资产: $125,430.00 | CEX: $80,200 | 链上: $35,230",
        "is_read": false,
        "created_at": "2026-02-10 09:00:00"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 25,
      "total_pages": 2,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### 11.2 获取未读数量

```
GET /api/v1/notifications/unread-count
```

**响应**：
```json
{
  "code": 0,
  "data": { "count": 3 }
}
```

### 11.3 标记为已读

```
POST /api/v1/notifications/{id}/read
```

### 11.4 全部已读

```
POST /api/v1/notifications/read-all
```

### 11.5 获取通知偏好

```
GET /api/v1/notifications/preferences
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "preferences": {
      "email_enabled": false,
      "push_enabled": true,
      "price_alert": true,
      "portfolio_alert": true,
      "system_notice": true
    }
  }
}
```

### 11.6 更新通知偏好

```
PUT /api/v1/notifications/preferences
```

**请求体**：
```json
{
  "push_enabled": true,
  "price_alert": true,
  "portfolio_alert": true,
  "system_notice": true
}
```

---

## 12. WebPush 推送接口

### 12.1 获取 VAPID 公钥

```
GET /api/v1/notifications/webpush/vapid
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "vapid_public_key": "BEl62iUYgU..."
  }
}
```

### 12.2 订阅推送

```
POST /api/v1/notifications/webpush/subscribe
```

**请求体**：
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/...",
  "keys": {
    "p256dh": "...",
    "auth": "..."
  }
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| endpoint | string | 是 | 推送服务端点 URL |
| keys.p256dh | string | 是 | P-256 Diffie-Hellman 公钥 |
| keys.auth | string | 是 | 认证密钥 |

### 12.3 退订推送

```
POST /api/v1/notifications/webpush/unsubscribe
```

**请求体**：
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/..."
}
```

---

## 13. 价格预警接口

### 13.1 获取预警列表

```
GET /api/v1/alerts
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "alerts": [
      {
        "id": 1,
        "symbol": "BTC",
        "condition": "above",
        "target_price": 100000,
        "current_price": 64000,
        "is_active": true,
        "triggered_at": null,
        "note": "突破十万提醒",
        "created_at": "2026-02-10T00:00:00Z"
      }
    ]
  }
}
```

### 13.2 创建预警

```
POST /api/v1/alerts
```

**请求体**：
```json
{
  "symbol": "BTC",
  "condition": "above",
  "target_price": 100000,
  "note": "突破十万提醒"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| symbol | string | 是 | 币种代码 |
| condition | string | 是 | 条件（above/below） |
| target_price | float64 | 是 | 目标价格（最小 0） |
| note | string | 否 | 备注 |

### 13.3 更新预警

```
PUT /api/v1/alerts/{id}
```

**请求体**：所有字段可选，额外支持 `is_active` (bool) 控制启用/禁用。

### 13.4 删除预警

```
DELETE /api/v1/alerts/{id}
```

---

## 14. 报告接口

### 14.1 获取报告列表

```
GET /api/v1/reports?type=daily&limit=20
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| type | string | 否 | 报告类型（daily/weekly/monthly/annual） | 全部 |
| limit | int | 否 | 返回数量 | 20 |

**响应**：
```json
{
  "code": 0,
  "data": {
    "reports": [
      {
        "id": 1,
        "type": "daily",
        "title": "2026-02-10 日报",
        "total_value": 125430,
        "pnl": 1230,
        "pnl_percent": 0.99,
        "created_at": "2026-02-10T21:00:00Z"
      }
    ]
  }
}
```

### 14.2 获取报告详情

```
GET /api/v1/reports/{id}
```

### 14.3 生成报告

```
POST /api/v1/reports/generate
```

**请求体**：
```json
{
  "type": "daily"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 报告类型（daily/weekly/monthly/annual） |

### 14.4 月度报告

```
GET /api/v1/reports/monthly/{month}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| month | string | 是 | 月份（格式: 2024-01，路径参数） |

### 14.5 年度报告

```
GET /api/v1/reports/annual/{year}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| year | string | 是 | 年份（格式: 2024，路径参数） |

### 14.6 报告对比

```
GET /api/v1/reports/compare?report_id_1=1&report_id_2=2
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| report_id_1 | int | 是 | 报告 1 ID |
| report_id_2 | int | 是 | 报告 2 ID |

**响应**：
```json
{
  "code": 0,
  "data": {
    "report_1": {
      "id": 1,
      "type": "daily",
      "period": "2026-02-10",
      "total_value": 125430,
      "change": 1230,
      "change_percent": 0.99,
      "generated_at": "2026-02-10T21:00:00Z"
    },
    "report_2": {
      "id": 2,
      "type": "daily",
      "period": "2026-02-09",
      "total_value": 124200,
      "change": 800,
      "change_percent": 0.65,
      "generated_at": "2026-02-09T21:00:00Z"
    },
    "value_diff": 1230,
    "change_diff": 430
  }
}
```

---

## 15. DeFi 接口

### 15.1 获取 DeFi 仓位

```
GET /api/v1/defi/positions?currency=USD
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| currency | string | 否 | 计价货币 | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "positions": [
      {
        "protocol": "Uniswap V3",
        "type": "liquidity",
        "token": "ETH-USDC",
        "amount": 2.5,
        "value": 15000,
        "apy": 12.5,
        "chain": "ethereum",
        "wallet_addr": "0x1234..."
      }
    ],
    "total_value": 45000,
    "currency": "USD"
  }
}
```

### 15.2 获取支持的协议

```
GET /api/v1/defi/protocols
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "protocols": [
      {
        "name": "Uniswap V3",
        "chains": ["ethereum", "polygon", "arbitrum"],
        "types": ["liquidity"],
        "is_active": true
      }
    ]
  }
}
```

### 15.3 DeFi 统计

```
GET /api/v1/defi/stats
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "total_value_locked": 45000,
    "position_count": 5,
    "by_protocol": { "Uniswap V3": 15000, "Aave": 20000, "Lido": 10000 },
    "by_chain": { "ethereum": 35000, "polygon": 10000 },
    "by_type": { "lending": 20000, "liquidity": 15000, "staking": 10000 }
  }
}
```

---

## 16. NFT 接口

### 16.1 获取 NFT 资产

```
GET /api/v1/nft/assets
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "assets": [
      {
        "id": 1,
        "contract_address": "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
        "token_id": "1234",
        "name": "Bored Ape #1234",
        "description": "...",
        "image_url": "ipfs://...",
        "collection": "Bored Ape Yacht Club",
        "chain": "ethereum",
        "floor_price": 12.5,
        "estimated_value": 28500,
        "wallet_addr": "0x1234..."
      }
    ],
    "total_value": 28500
  }
}
```

### 16.2 获取收藏集统计

```
GET /api/v1/nfts/collections
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "collections": [
      {
        "collection": "Bored Ape Yacht Club",
        "count": 2,
        "total_floor_price": 57000,
        "chain": "ethereum"
      }
    ],
    "total_count": 5,
    "total_value": 85000
  }
}
```

---

## 17. 交易记录接口

### 17.1 获取交易记录

```
GET /api/v1/transactions?page=1&page_size=20&type=buy&source=binance&start=2026-01-01&end=2026-02-10
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页数量 | 20 |
| type | string | 否 | 交易类型（buy/sell/transfer/deposit/withdraw） | 全部 |
| source | string | 否 | 来源筛选（binance/okx/coinbase/ethereum 等） | 全部 |
| start | string | 否 | 开始日期（YYYY-MM-DD） | |
| end | string | 否 | 结束日期（YYYY-MM-DD） | |
| cursor | string | 否 | 游标分页（RFC3339 格式时间戳） | |

**响应**：
```json
{
  "code": 0,
  "data": {
    "transactions": [
      {
        "id": 1,
        "source": "binance",
        "tx_type": "buy",
        "symbol": "BTC",
        "amount": 0.1,
        "price": 64000,
        "total": 6400,
        "fee": 6.4,
        "fee_coin": "USDT",
        "tx_hash": "",
        "timestamp": "2026-02-10T08:30:00Z"
      }
    ],
    "total": 200,
    "page": 1,
    "page_size": 20
  }
}
```

### 17.2 同步交易记录

```
POST /api/v1/transactions/sync
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "同步完成",
    "synced_count": 12
  }
}
```

### 17.3 交易统计

```
GET /api/v1/transactions/stats
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "total_transactions": 200,
    "total_volume": 150000,
    "total_fees": 320,
    "by_type": { "buy": 120, "sell": 50, "transfer": 20, "deposit": 10 },
    "by_source": { "binance": 150, "ethereum": 50 }
  }
}
```

### 17.4 获取同步设置

```
GET /api/v1/settings/tx-sync
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "settings": {
      "auto_sync": false,
      "sync_interval": 360,
      "last_sync_at": "2026-02-10T08:00:00Z"
    }
  }
}
```

### 17.5 更新同步设置

```
PUT /api/v1/settings/tx-sync
```

**请求体**：
```json
{
  "auto_sync": true,
  "sync_interval": 180
}
```

---

## 18. 费用分析接口

### 18.1 获取费用分析

```
GET /api/v1/analytics/fees?range=30d&currency=USD
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| range | string | 否 | 时间范围（7d/30d/90d/1y） | 30d |
| currency | string | 否 | 计价货币 | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "total_fees": 456.78,
    "trading_fees": 280.50,
    "gas_fees": 126.28,
    "currency": "USD",
    "breakdown": [
      { "source": "binance", "type": "trading_fee", "amount": 200, "currency": "USD", "count": 50 },
      { "source": "ethereum", "type": "gas_fee", "amount": 126.28, "currency": "USD", "count": 30 }
    ],
    "daily_trend": [
      { "date": "2026-02-10", "amount": 12.5 }
    ]
  }
}
```

---

## 19. 分析接口（盈亏/归因/预测）

### 19.1 每日盈亏

```
GET /api/v1/analytics/pnl/daily?days=30
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| days | int | 否 | 查询天数 | 30 |

**响应**：
```json
{
  "code": 0,
  "data": {
    "daily": [
      {
        "date": "2026-02-10",
        "pnl": 320,
        "pnl_percent": 0.58,
        "total_value": 55320
      }
    ]
  }
}
```

### 19.2 盈亏摘要

```
GET /api/v1/analytics/pnl/summary
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "total_pnl": 8500,
    "total_pnl_percent": 7.2,
    "pnl_7d": 1200,
    "pnl_30d": 3200,
    "pnl_90d": 8500,
    "best_day": "2026-02-05",
    "worst_day": "2026-01-20",
    "best_day_pnl": 1500,
    "worst_day_pnl": -800
  }
}
```

### 19.3 资产归因分析

```
GET /api/v1/analytics/attribution?days=7&currency=USD
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| days | int | 否 | 分析天数 | 7 |
| currency | string | 否 | 计价货币 | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "total_return": 1500,
    "total_percent": 2.8,
    "attributions": [
      { "symbol": "BTC", "source": "binance", "contribution": 800, "weight": 42, "return": 3.2 },
      { "symbol": "ETH", "source": "ethereum", "contribution": 400, "weight": 28, "return": 2.1 }
    ],
    "days": 7,
    "currency": "USD"
  }
}
```

### 19.4 趋势预测

```
GET /api/v1/analytics/forecast?days=30
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| days | int | 否 | 预测天数 | 30 |

**响应**：
```json
{
  "code": 0,
  "data": {
    "forecast_points": [
      {
        "date": "2026-02-13",
        "value": 55500,
        "lower": 54000,
        "upper": 57000
      }
    ],
    "trend": "up",
    "confidence": 0.72,
    "slope": 82.5,
    "days": 30
  }
}
```

---

## 20. 策略引擎接口

### 20.1 获取策略列表

```
GET /api/v1/strategies
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "strategies": [
      {
        "id": 1,
        "name": "BTC/ETH 7:3 策略",
        "type": "rebalance",
        "config": {
          "targets": [
            { "symbol": "BTC", "percentage": 70 },
            { "symbol": "ETH", "percentage": 30 }
          ],
          "rebalance_threshold": 5.0
        },
        "is_active": true,
        "created_at": "2026-01-15T00:00:00Z",
        "updated_at": "2026-02-10T10:00:00Z"
      }
    ]
  }
}
```

### 20.2 创建策略

```
POST /api/v1/strategies
```

**请求体**：
```json
{
  "name": "BTC/ETH 7:3 策略",
  "type": "rebalance",
  "config": {
    "targets": [
      { "symbol": "BTC", "percentage": 70 },
      { "symbol": "ETH", "percentage": 30 }
    ],
    "rebalance_threshold": 5.0
  }
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 策略名称（最大 100 字符） |
| type | string | 是 | 策略类型（rebalance/dca/stop_limit） |
| config | object | 是 | 策略配置（JSON 对象） |

### 20.3 更新策略

```
PUT /api/v1/strategies/{id}
```

**请求体**：所有字段可选，额外支持 `is_active` (bool) 控制启用/禁用。

### 20.4 删除策略

```
DELETE /api/v1/strategies/{id}
```

### 20.5 策略分析

```
GET /api/v1/strategies/{id}/analysis
```

**功能**：获取策略的偏差分析和再平衡建议

**响应**：
```json
{
  "code": 0,
  "data": {
    "analysis": {
      "strategy_id": 1,
      "current_alloc": { "BTC": 75, "ETH": 25 },
      "target_alloc": { "BTC": 70, "ETH": 30 },
      "deviation": { "BTC": 5, "ETH": -5 },
      "recommendations": [
        {
          "symbol": "BTC",
          "action": "sell",
          "amount": 0.04,
          "value_usd": 2750,
          "reason": "BTC 超配 5%，建议卖出"
        },
        {
          "symbol": "ETH",
          "action": "buy",
          "amount": 1.2,
          "value_usd": 2750,
          "reason": "ETH 欠配 5%，建议买入"
        }
      ]
    }
  }
}
```

### 20.6 再平衡建议

```
GET /api/v1/strategies/{id}/rebalance
```

**功能**：获取再平衡建议（等同于策略分析接口，返回格式相同）

---

## 21. 成就系统接口

### 21.1 获取成就列表

```
GET /api/v1/achievements
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "achievements": [
      {
        "id": "first_deposit",
        "name": "初入币圈",
        "description": "完成首次资产录入",
        "icon": "rocket",
        "category": "beginner",
        "is_unlocked": true,
        "unlocked_at": "2026-01-15T00:00:00Z",
        "progress": 100
      },
      {
        "id": "diamond_hands",
        "name": "钻石手",
        "description": "持仓超过 365 天未卖出",
        "icon": "diamond",
        "category": "advanced",
        "is_unlocked": false,
        "progress": 33
      }
    ],
    "total_count": 11,
    "unlocked_count": 3
  }
}
```

### 21.2 检查成就进度

```
POST /api/v1/achievements/check
```

**功能**：触发成就进度检查，返回新解锁的成就

**响应**：
```json
{
  "code": 0,
  "data": {
    "newly_unlocked": [
      {
        "id": "whale",
        "name": "鲸鱼",
        "description": "总资产超过 10 万美元"
      }
    ]
  }
}
```

---

## 22. 基准对比接口

### 22.1 获取基准对比

```
GET /api/v1/benchmark?range=30d
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| range | string | 否 | 周期（7d/30d/90d/1y） | 30d |

**响应**：
```json
{
  "code": 0,
  "data": {
    "series": [
      {
        "name": "Portfolio",
        "points": [
          { "date": "2026-01-12", "value": 100 },
          { "date": "2026-02-10", "value": 106.1 }
        ],
        "return": 6.1
      },
      {
        "name": "BTC",
        "points": [
          { "date": "2026-01-12", "value": 100 },
          { "date": "2026-02-10", "value": 108.5 }
        ],
        "return": 8.5
      },
      {
        "name": "ETH",
        "points": [ "..." ],
        "return": 12.3
      },
      {
        "name": "S&P500",
        "points": [ "..." ],
        "return": 2.1
      }
    ],
    "range": "30d",
    "start_date": "2026-01-12",
    "end_date": "2026-02-10"
  }
}
```

---

## 23. 目标追踪接口

### 23.1 获取目标列表

```
GET /api/v1/goals
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "goals": [
      {
        "id": 1,
        "title": "百万目标",
        "type": "asset_value",
        "target_value": 1000000,
        "current_value": 125430,
        "currency": "USD",
        "progress": 12.5,
        "deadline": "2028-12-31",
        "is_completed": false,
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-02-10T00:00:00Z"
      }
    ]
  }
}
```

### 23.2 创建目标

```
POST /api/v1/goals
```

**请求体**：
```json
{
  "title": "百万目标",
  "type": "asset_value",
  "target_value": 1000000,
  "currency": "USD",
  "deadline": "2028-12-31"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 目标标题（最大 100 字符） |
| type | string | 是 | 目标类型（asset_value/holding_amount/return_rate） |
| target_value | float64 | 是 | 目标值（最小 0） |
| currency | string | 否 | 货币（默认 USD） |
| deadline | string | 否 | 截止日期（ISO 8601 格式） |

### 23.3 更新目标

```
PUT /api/v1/goals/{id}
```

### 23.4 删除目标

```
DELETE /api/v1/goals/{id}
```

---

## 24. 资产健康评分接口

### 24.1 获取健康评分

```
GET /api/v1/portfolio/health?currency=USD
```

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| currency | string | 否 | 计价货币 | USD |

**响应**：
```json
{
  "code": 0,
  "data": {
    "overall_score": 72,
    "level": "good",
    "details": [
      {
        "category": "diversification",
        "score": 65,
        "weight": 0.3,
        "description": "资产集中度偏高",
        "suggestion": "建议增加稳定币比例以降低波动"
      },
      {
        "category": "stability",
        "score": 80,
        "weight": 0.25,
        "description": "收益稳定性良好",
        "suggestion": ""
      },
      {
        "category": "risk",
        "score": 70,
        "weight": 0.25,
        "description": "风险水平适中",
        "suggestion": "考虑分散到更多资产类别"
      },
      {
        "category": "growth",
        "score": 75,
        "weight": 0.2,
        "description": "增长趋势向好",
        "suggestion": ""
      }
    ],
    "currency": "USD",
    "updated_at": "2026-02-10T10:30:00Z"
  }
}
```

---

## 25. 市场数据接口

### 25.1 Gas 费查询

```
GET /api/v1/market/gas
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "prices": [
      {
        "chain": "Ethereum",
        "low": 10.0,
        "standard": 15.0,
        "fast": 20.0,
        "instant": 25.0,
        "base_fee": 12.5,
        "level": "medium"
      },
      {
        "chain": "BSC",
        "low": 1.0,
        "standard": 3.0,
        "fast": 5.0,
        "instant": 7.0,
        "base_fee": 1.0,
        "level": "low"
      },
      {
        "chain": "Polygon",
        "low": 30.0,
        "standard": 50.0,
        "fast": 80.0,
        "instant": 100.0,
        "base_fee": 30.0,
        "level": "low"
      }
    ],
    "updated_at": 1707580800
  }
}
```

> **单位**：所有 Gas 价格单位为 Gwei

> **拥堵等级**：`level` 取值为 `low`（低）/ `medium`（中）/ `high`（高）

---

## 26. 用户设置接口

### 26.1 获取设置

```
GET /api/v1/users/settings
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "settings": {
      "language": "zh-CN",
      "currency": "USDC",
      "theme": "nexus-pro",
      "timezone": "Asia/Shanghai",
      "refresh_interval": "3600",
      "snapshot_time": "00:00"
    }
  }
}
```

### 26.2 更新设置

```
PUT /api/v1/users/settings
```

**请求体**：
```json
{
  "settings": {
    "language": "zh-CN",
    "currency": "USDC",
    "theme": "nexus-pro"
  }
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| settings | map[string]string | 是 | 设置键值对 |

### 26.3 导出数据

```
GET /api/v1/users/export
```

**功能**：导出所有用户数据（JSON 格式，不含加密凭证）

**响应**：
```json
{
  "code": 0,
  "data": {
    "exchange_accounts": [
      { "id": 1, "exchange_name": "binance", "label": "主账户", "note": "" }
    ],
    "wallet_addresses": [
      { "id": 1, "blockchain": "ethereum", "address": "0x1234...", "label": "主钱包" }
    ],
    "manual_assets": [
      { "id": 1, "asset_type": "bank", "asset_name": "银行存款", "amount": 70000, "amount_usd": 9660, "currency": "CNY", "notes": "" }
    ],
    "strategies": [
      { "id": 1, "name": "BTC/ETH 7:3", "type": "rebalance", "config": "{...}", "is_active": true }
    ],
    "goals": [
      { "id": 1, "title": "百万目标", "type": "asset_value", "target_value": 1000000, "currency": "USD", "deadline": "2028-12-31" }
    ],
    "price_alerts": [
      { "id": 1, "symbol": "BTC", "condition": "above", "target_price": 100000, "is_active": true, "note": "" }
    ],
    "settings": { "language": "zh-CN", "currency": "USDC" },
    "exported_at": "2026-02-12T10:00:00Z"
  }
}
```

### 26.4 清除缓存

```
POST /api/v1/users/clear-cache
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "缓存已清除"
  }
}
```

### 26.5 重置设置

```
POST /api/v1/users/reset-settings
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "message": "设置已重置为默认值"
  }
}
```

---

## 27. 系统接口

### 27.1 健康检查

```
GET /api/v1/health
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "status": "ok",
    "version": "0.2.0",
    "timestamp": 1707580800
  }
}
```

---

## 相关文档

- [Swagger UI（交互式文档）](http://localhost:8080/api/v1/docs)
- [技术基线](./tech-baseline.md)
- [部署指南](../guides/deployment-guide.md)

---

**文档维护者**: @allfi
**最后更新**: 2026-02-12
**版本**: v3.0
