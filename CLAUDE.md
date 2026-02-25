# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## ⚠️ 重要规范：中文优先原则

**本项目强制要求使用中文进行所有文档化工作：**

1. **代码注释必须使用中文**
   - 所有函数、方法、结构体、类型的注释必须用中文编写
   - 复杂逻辑的行内注释必须用中文
   - TODO、FIXME 等标记后的说明必须用中文

2. **文档内容必须使用中文**
   - 所有 Markdown 文档（`.md` 文件）的正文必须用中文
   - API 文档、技术设计文档、README 等均使用中文
   - 代码示例中的说明性文字使用中文

3. **文档文件名使用中文拼音或中文**
   - 推荐使用有意义的中文拼音命名（如 `biz-overview.md`）
   - 核心文档可使用中文文件名（如 `技术基线.md`）
   - 避免使用纯英文文件名（除非是约定俗成的如 `README.md`）

4. **例外情况**
   - 代码标识符（变量名、函数名、包名等）仍使用英文（遵循语言规范）
   - Git commit message 可以使用中文或英文
   - 第三方库、框架名称保持原样
   - 技术术语可保留英文（如 API、HTTP、JSON）

5. **文档维护**:
  - 每次功能修改、新增或删除后，必须同步更新 README.md 和 CLAUDE.md
  - README.md 侧重于项目概述、功能列表、使用说明
  - CLAUDE.md 侧重于开发规范、技术细节、实现进度
  - 前端 `webapp/`、后端 `core/` 分别维护各自的 README.md 和 CLAUDE.md
  - 所有项目文档统一存放在根目录 `docs/` 下，子目录中不放文档
  - 完成任务后提交代码

**示例对比：**

✅ **正确示例**：
```go
// 获取用户的所有资产总值（按指定货币计价）
// currency: 计价货币（USDC/BTC/ETH/CNY）
// 返回：总资产价值和错误信息
func GetTotalAssetValue(userID int64, currency string) (float64, error) {
    // 并发获取 CEX、区块链、手动资产的余额
    var wg sync.WaitGroup
    // ... 具体实现
}
```

❌ **错误示例**：
```go
// Get total asset value for user in specified currency
// currency: pricing currency (USDC/BTC/ETH/CNY)
// Returns: total asset value and error
func GetTotalAssetValue(userID int64, currency string) (float64, error) {
    // Fetch balances from CEX, blockchain and manual assets concurrently
    var wg sync.WaitGroup
    // ... implementation
}
```

## 工作准则

- **语言要求**: 只用中文解释，文档名、文档内容、代码注释均用中文
- **任务管理**: todo list 要详细，任务拆解要细致，实现要全面
- **编码规范**: 正确使用 UTF-8 编码，确保中文显示正常
- **提交代码**: 每次完成一个任务，提交一次代码
- **Go 版本要求**: Go 1.24
- **前端包管理**: 使用 pnpm，禁止使用 npm/yarn

## 版本管理规则

1. **版本号通过 Git Tag 管理**
   - 版本号直接使用 Git Tag（如 `v0.1.8`），遵循语义化版本（SemVer）
   - 后端版本号通过 `ldflags` 在构建时注入到 `core/internal/version/version.go`
   - 前端从后端 API 获取版本号
   - 发版流程：创建并推送 Git Tag（如 `git tag v0.2.0 && git push origin v0.2.0`）

2. **README.md 不写历史记录**
   - `README.md` 只包含项目概述、功能列表、快速开始、技术栈等内容
   - **禁止**在 `README.md` 中写版本历史、更新日志、变更记录
   - 所有版本历史记录统一在 GitHub Releases 中查看
   - `CLAUDE.md` 中的实现进度也引用 GitHub Releases，不重复记录

3. **CHANGEHISTORY.md 自动生成规则**
   - `CHANGEHISTORY.md` 由 GitHub Action（`.github/workflows/release.yml`）**自动生成**
   - 当推送 `v*` 格式的 Git Tag 时，自动触发：
     1. 提取上一个 Tag 到当前 Tag 之间的所有 commit 记录
     2. 按提交前缀自动分类（功能新增 / 问题修复 / 重构优化 / 文档更新 / 构建与部署 / 其他变更）
     3. 生成完整的 `CHANGEHISTORY_full.md` 文件（包含所有历史版本）
     4. 将 `CHANGEHISTORY_full.md` 作为 Release Asset 上传到 GitHub Releases
     5. **不再提交任何文件到 master 分支**（避免产生额外的 bot commit）
   - **查看变更历史**：访问 [GitHub Releases](https://github.com/your-finance/allfi/releases) 页面，下载 `CHANGEHISTORY_full.md` 文件
   - **提交分类规则**：为了让自动分类更准确，建议 commit message 使用以下前缀：
     - `feat:` / `添加` / `新增` / `实现` / `支持` → 功能新增
     - `fix:` / `修复` / `修正` → 问题修复
     - `refactor:` / `重构` → 重构优化
     - `docs:` / `文档` → 文档更新
     - `build:` / `ci:` / `构建` / `部署` → 构建与部署

## Git 提交规则

### 提交信息格式

提交信息应该简洁明了，使用中文描述变更内容。

**不要**在 commit message 中包含以下内容：
- ❌ `🤖 Generated with [Claude Code](https://claude.com/claude-code)`
- ❌ `Co-Authored-By: Claude <noreply@anthropic.com>`
- ❌ 任何 Claude Code 相关的标记或签名

### 提交信息示例

```
添加ASan支持到build.sh

- 在build.sh中添加asan参数检查
- 支持 release/debug 模式下的ASan编译
- 修改CMakeLists.txt配置对应选项
```

### 提交前检查

在提交代码前，确保：
1. 提交信息使用中文
2. 提交信息简洁明了，说明改动的目的
3. 不包含任何工具生成的标记
4. 必要时包含详细的变更说明

---

## 项目技术栈

### 后端（`core/`）
- **Go 1.24** + **GoFrame v2**（HTTP 路由 + ORM + 配置）
- **GoFrame ORM** — SQLite3（纯 Go 实现，无需 CGO）/ MySQL 可选
- **SQLite 驱动**：`gogf/gf/contrib/drivers/sqlite/v2` → `glebarez/go-sqlite` → `modernc.org/sqlite`
- **架构**：api/v1/ → controller/ → service/ → logic/ → dao/ → Database
- **服务注册**：`init()` 自动注册 + `cmd/server/main.go` blank import 触发

### 前端（`webapp/`）
- **Vue 3** Composition API + `<script setup>`
- **Vite 7** + **pnpm** 构建
- **Pinia** 状态管理、**Vue Router 4**、**Chart.js 4**
- **Phosphor Icons**（`@phosphor-icons/vue`）
- **Mock/Real API 切换**：`VITE_USE_MOCK_API` 环境变量
- **i18n**：简体中文 / 繁體中文 / English
- **4 套主题**：Nexus Pro / Vestia / XChange / Aurora

---

## 已实现功能

详细版本历史记录见 [GitHub Releases](https://github.com/your-finance/allfi/releases)。

---

## RPC 管理规则

### 免费 RPC（用户未配置 API Key / RPC URL 时的默认行为）

1. **动态获取**：通过 Chainlist API（`https://chainlist.org/rpcs.json`）获取各链的免费公共 RPC 端点列表
2. **过滤商业节点**：自动排除需要 API Key 的商业 RPC 提供商（如 Alchemy、Infura、Ankr、BlastAPI 等）
3. **隐私优先**：优先选择 `tracking: "none"` 的 RPC 节点
4. **可用性验证**：获取 RPC URL 后，使用 `eth_gasPrice` 请求验证可用性，只保存验证通过的 URL
5. **刷新频率**：1 小时更新一次，避免频繁请求 Chainlist
6. **Fallback**：动态获取失败时，使用 `chains.go` 中的硬编码 PublicRPC 作为最终兜底

### 用户自定义 RPC（优先级最高）

1. **当用户配置了对应链的 RPC URL 后，仅使用用户输入的 RPC**，不再使用免费 RPC
2. **配置方式**：通过前端 Settings → API 配置页面输入，或通过环境变量设置
3. **支持的链**：`ethereum_rpc`、`bsc_rpc`、`polygon_rpc`、`arbitrum_rpc`、`optimism_rpc`、`base_rpc`
4. **存储方式**：数据库 AES 加密存储（优先级最高）> 环境变量 > 空字符串

### RPC URL 解析优先级（`GetRPCURL()` 函数）

```
① 用户自定义 RPC（数据库/环境变量） → ② Chainlist 动态免费 RPC → ③ 硬编码 Fallback RPC
```

### RPC 的能力范围

RPC（无论免费或用户自定义）可用于以下**链上只读操作**：
- ✅ 获取 Gas 价格（`eth_gasPrice`）
- ✅ 查询原生代币余额（`eth_getBalance`）— ETH/BNB/MATIC 等
- ✅ 未来可扩展：ERC20 代币余额查询（`eth_call` → `balanceOf()`）
- ✅ 未来可扩展：授权操作（需签名的写操作）

**已知限制**：
- ❌ 无 Etherscan API Key 时，无法自动扫描钱包内的所有 ERC20 代币（需已知合约地址）
- ❌ RPC 不支持历史交易记录查询（需 Etherscan API）

### 相关代码文件

| 文件 | 作用 |
|------|------|
| `internal/integrations/etherscan/rpc_manager.go` | RPC URL 管理：Chainlist 动态获取 + 用户自定义解析 |
| `internal/integrations/etherscan/rpc.go` | 通过 RPC 获取 Gas 价格和原生代币余额 |
| `internal/integrations/etherscan/chains.go` | 链配置：ChainID、PublicRPC fallback 等 |
| `internal/utils/apikey.go` | API Key / RPC URL 解析（数据库 > 环境变量） |
| `internal/app/system/logic/apikey.go` | API Key 管理（加密存储、脱敏查询） |
| `internal/app/wallet/logic/wallet.go` | 钱包余额查询（自动降级：Etherscan → RPC） |
| `internal/app/market/logic/market.go` | Gas 价格查询（自动降级：Etherscan → RPC） |

---

## 行情与汇率获取规则

1. **加密货币汇率**：汇率获取需要使用 Binance Spot，或 Gate.io 等 CEX 上的行情。
2. **法币汇率**：CNY 的行情通过 Yahoo 获取。

---

## 后端核心服务一览

| 服务 | 文件 | 接口 |
|------|------|------|
| 资产聚合 | `service/asset_service.go` | `AssetServiceInterface` |
| 交易所 | `service/exchange_service.go` | `ExchangeServiceInterface` |
| 区块链 | `service/blockchain_service.go` | `BlockchainServiceInterface` |
| 价格 | `service/price_service.go` | `PriceServiceInterface` |
| 快照 | `service/snapshot_service.go` | `SnapshotServiceInterface` |
| DeFi | `service/defi_service.go` | `DeFiServiceInterface` |
| NFT | `service/nft_service.go` | `NFTServiceInterface` |
| 交易记录 | `service/transaction_service.go` | `TransactionServiceInterface` |
| 费用分析 | `service/fee_service.go` | `FeeServiceInterface` |
| 策略引擎 | `service/strategy_service.go` | `StrategyServiceInterface` |
| 成就系统 | `service/achievement_service.go` | `AchievementServiceInterface` |
| 基准对比 | `service/benchmark_service.go` | `BenchmarkServiceInterface` |
| 通知 | `service/notification_service.go` | `NotificationServiceInterface` |
| 价格预警 | `service/price_alert_service.go` | `PriceAlertServiceInterface` |
| 报告 | `service/report_service.go` | `ReportServiceInterface` |
| WebPush | `service/webpush_service.go` | — |
| 系统管理 | `internal/app/system/logic/system.go` | `ISystem` |

所有服务接口定义在 `service/interfaces.go`，入口文件 `cmd/server/main.go` 中完成依赖注入。
