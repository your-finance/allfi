# AllFi 开发指南

> **版本**：v2.1 | **更新时间**：2026-02-13

---

## 目录

1. [开发环境搭建](#1-开发环境搭建)
2. [代码规范](#2-代码规范)
3. [项目结构](#3-项目结构)
4. [开发流程](#4-开发流程)
5. [测试指南](#5-测试指南)
6. [贡献指南](#6-贡献指南)

---

## 1. 开发环境搭建

### 1.0 一键启动（推荐）

项目提供了 Makefile 和快速启动脚本，可以快速完成环境搭建：

```bash
# 方式一：Makefile
make setup    # 自动生成 .env + 安装前后端依赖
make dev      # 同时启动前后端开发服务器

# 方式二：快速启动脚本
bash scripts/quickstart.sh          # 完整启动
bash scripts/quickstart.sh --mock   # 仅前端 Mock 模式
bash scripts/quickstart.sh --check  # 检测依赖环境
```

#### Makefile 命令速查

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有命令 |
| `make setup` | 一键初始化（生成 .env + 安装依赖） |
| `make dev` | 同时启动前后端 |
| `make dev-mock` | 纯前端 Mock 模式（无需后端） |
| `make dev-backend` | 仅启动后端 |
| `make dev-frontend` | 仅启动前端 |
| `make build` | 构建前后端 |
| `make docker` | Docker Compose 启动 |
| `make docker-build` | Docker 重新构建并启动 |
| `make health` | 健康检查 |
| `make swagger` | 打开 Swagger UI |
| `make clean` | 清理构建产物 |

> 如果不想使用 Makefile，也可以按下面的步骤手动搭建。

### 1.1 安装开发工具

#### 必需工具

```bash
# Go 1.24（项目固定版本）
go version   # 输出应为 go1.24.x

# Node.js 18+
node --version

# pnpm（包管理器）
pnpm --version

# Git
git --version
```

#### 推荐工具

```bash
# Go 热重载
go install github.com/cosmtrek/air@latest

# Go 代码检查
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 1.2 克隆代码

```bash
git clone https://github.com/your-finance/allfi.git
cd allfi
```

### 1.3 后端开发环境

```bash
cd core

# 安装 Go 依赖
go mod download && go mod verify

# 配置文件（开发环境默认配置即可使用）
# 如需自定义，编辑 manifest/config/config.yaml

# 启动后端（热重载）
air

# 或手动运行
go run cmd/server/main.go
```

验证：

```bash
curl http://localhost:8080/api/v1/health
```

> **数据库**：首次启动时 GORM AutoMigrate 自动创建所有表，无需手动迁移。

### 1.4 前端开发环境

```bash
cd webapp

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

访问 http://localhost:3174

### 1.5 Mock 模式开发

前端支持 Mock 数据模式，无需启动后端即可开发和体验 UI：

```bash
cd webapp

# 方式一：pnpm 脚本（推荐）
pnpm dev:mock

# 方式二：手动设置环境变量
VITE_USE_MOCK_API=true pnpm dev
```

Mock 模式使用 `webapp/src/data/` 下的模拟数据，所有 API 调用返回预设数据。

> 切换回真实 API 模式：使用 `pnpm dev` 即可（需要后端运行）。

### 1.6 代理配置（国内用户）

如果下载依赖或访问外部 API 较慢，请参考 [代理配置指南](./proxy-guide.md) 进行配置：

```bash
# Go 模块代理
export GOPROXY=https://goproxy.cn,direct

# pnpm 镜像源
pnpm config set registry https://registry.npmmirror.com
```

### 1.7 IDE 配置（VSCode）

推荐扩展：
- **Go** — Go 语言支持
- **Vue - Official (Volar)** — Vue 3 支持
- **Tailwind CSS IntelliSense** — Tailwind 智能提示

工作区配置 `.vscode/settings.json`：

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[vue]": {
    "editor.defaultFormatter": "Vue.volar"
  }
}
```

---

## 2. 代码规范

### 2.1 后端规范（Go）

#### 命名规范

```go
// 类型：PascalCase
type ExchangeService struct {}

// 函数：PascalCase（导出）/ camelCase（未导出）
func (s *ExchangeService) GetBalance() error {}

// 常量：PascalCase
const DefaultCurrency = "USDC"

// 错误：Err 前缀
var ErrInvalidAPIKey = errors.New("invalid API key")
```

#### 注释规范（必须中文）

```go
// GetAssetsSummary 获取资产总览
// 参数：
//   - currency: 计价币种（USDC/BTC/ETH/CNY）
// 返回：
//   - *AssetSummary: 资产总览数据
//   - error: 错误信息
func (s *AssetService) GetAssetsSummary(currency string) (*AssetSummary, error) {
    // 验证计价币种
    if !isValidCurrency(currency) {
        return nil, ErrInvalidCurrency
    }
    // 并发获取三类资产
    var wg sync.WaitGroup
    // ...
}
```

#### 错误处理

```go
// 使用有意义的错误信息
if err != nil {
    return nil, fmt.Errorf("获取交易所余额失败: %w", err)
}

// 使用预定义错误码（见 utils/response.go）
utils.Error(r, utils.CodeExternalAPIError, "交易所 API 调用失败")
```

### 2.2 前端规范（Vue 3）

#### 必须使用 Composition API

```vue
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAssetStore } from '../stores/assetStore'

const assetStore = useAssetStore()
const totalValue = computed(() => assetStore.totalValue)

onMounted(() => {
  assetStore.loadSummary()
})
</script>
```

#### CSS 规范

- 必须使用 `<style scoped>`
- 必须使用 CSS 变量（`--gap-*`、`--radius-*`、`--color-*`），禁止硬编码
- 数字必须使用 `.font-mono` 类
- 禁止 `backdrop-filter: blur()`、发光边框等装饰效果

#### 多语言

所有用户可见文本必须通过 `useI18n()` 的 `t()` 函数翻译：

```vue
<script setup>
import { useI18n } from '../composables/useI18n'
const { t } = useI18n()
</script>

<template>
  <h1>{{ t('dashboard.title') }}</h1>
</template>
```

### 2.3 提交规范

```bash
# 格式：中文动词 + 简要描述
git commit -m "添加资产总览 API"
git commit -m "修复交易所连接失败的问题"
git commit -m "优化资产刷新性能"

# 禁止包含 Claude Code 标记
```

| 前缀 | 说明 | 示例 |
|-----|------|------|
| 添加 | 新功能 | 添加 Gas 费查询功能 |
| 修复 | Bug 修复 | 修复资产计算错误 |
| 优化 | 性能优化 | 优化并发查询性能 |
| 重构 | 代码重构 | 重构资产服务层 |
| 文档 | 文档更新 | 更新 API 文档 |

---

## 3. 项目结构

### 3.1 后端结构

```
core/
├── cmd/
│   ├── server/main.go              # 服务器入口
│   └── migrate/main.go             # 数据库迁移工具
├── internal/
│   ├── api/
│   │   ├── router.go               # 路由注册（75 条路由）
│   │   ├── handlers/               # 24 个 Handler 文件
│   │   │   ├── auth.go             # 认证（PIN 设置/登录/修改）
│   │   │   ├── asset.go            # 资产总览/详情/历史/刷新
│   │   │   ├── exchange.go         # 交易所账户 CRUD
│   │   │   ├── wallet.go           # 钱包地址 CRUD
│   │   │   ├── defi.go             # DeFi 仓位
│   │   │   ├── nft.go              # NFT 资产
│   │   │   ├── transaction.go      # 交易记录
│   │   │   ├── pnl.go              # 每日盈亏
│   │   │   ├── strategy.go         # 投资策略
│   │   │   ├── achievement.go      # 成就系统
│   │   │   └── ...                 # 其余 handler
│   │   └── middleware/             # 中间件
│   │       └── auth.go             # JWT 认证中间件
│   ├── service/
│   │   ├── interfaces.go           # 20+ 服务接口定义
│   │   ├── container.go            # ServiceContainer 依赖注入
│   │   ├── asset_service.go
│   │   ├── exchange_service.go
│   │   └── ...
│   ├── repository/                 # 数据访问层（18 个 Repo）
│   ├── models/                     # GORM 数据模型（18 张表）
│   │   ├── base.go                 # BaseModel + 常量定义
│   │   ├── exchange_account.go
│   │   ├── wallet_address.go
│   │   ├── manual_asset.go
│   │   └── ...
│   ├── integrations/               # 第三方集成（8 个模块）
│   │   ├── binance/                # Binance（go-binance/v2）
│   │   ├── okx/                    # OKX（CCXT）
│   │   ├── coinbase/               # Coinbase（CCXT）
│   │   ├── etherscan/              # Etherscan（6 条 EVM 链）
│   │   ├── coingecko/              # CoinGecko 价格
│   │   ├── yahoo/                  # Yahoo Finance 汇率
│   │   ├── alchemy/                # Alchemy NFT API
│   │   ├── defi/                   # DeFi 协议（7 个）
│   │   └── base_client.go          # HTTP 客户端基类
│   ├── cron/                       # 6 个定时任务
│   │   └── snapshot_job.go         # CronManager
│   └── utils/
│       ├── config.go               # GoFrame gcfg 配置
│       ├── response.go             # 统一响应 + 错误码
│       └── crypto.go               # AES-256-GCM 加密
├── manifest/config/
│   └── config.yaml                 # 应用配置文件
├── docs/
│   └── swagger.yaml                # OpenAPI 3.0 文档
└── go.mod                          # Go 1.24
```

### 3.2 前端结构

```
webapp/
├── src/
│   ├── App.vue                     # 根组件（侧边栏+顶栏布局）
│   ├── main.js                     # 入口文件
│   ├── api/                        # API 服务层（12 个文件）
│   │   ├── index.js                # 统一 API 客户端（Fetch）
│   │   ├── client.js               # HTTP 客户端封装
│   │   ├── defiService.js
│   │   ├── nftService.js
│   │   ├── transactionService.js
│   │   ├── strategyService.js
│   │   ├── achievementService.js
│   │   ├── feeService.js
│   │   ├── benchmarkService.js
│   │   ├── annualReportService.js
│   │   ├── marketService.js
│   │   └── mockData.js
│   ├── components/                 # 39 个可复用组件
│   │   ├── StatCard.vue            # 统计卡片
│   │   ├── AddAccountDialog.vue    # 添加账户对话框
│   │   ├── CommandPalette.vue      # 命令面板（Cmd+K）
│   │   ├── DashboardCustomizer.vue # 仪表盘自定义
│   │   ├── OnboardingWizard.vue    # 新手引导
│   │   ├── CalendarHeatmap.vue     # 日历热力图
│   │   ├── DeFiOverview.vue        # DeFi 概览
│   │   ├── NFTGallery.vue          # NFT 画廊
│   │   ├── TransactionTimeline.vue # 交易时间线
│   │   ├── StrategyPanel.vue       # 策略面板
│   │   ├── AchievementPanel.vue    # 成就面板
│   │   ├── BenchmarkPanel.vue      # 基准对比
│   │   ├── BottomNav.vue           # 移动端底部导航
│   │   ├── PullToRefresh.vue       # 下拉刷新
│   │   └── ...
│   ├── pages/                      # 9 个页面
│   │   ├── Dashboard.vue           # 仪表盘
│   │   ├── Accounts.vue            # 账户管理
│   │   ├── History.vue             # 历史记录
│   │   ├── Analytics.vue           # 数据分析
│   │   ├── Reports.vue             # 资产报告
│   │   ├── Settings.vue            # 系统设置
│   │   ├── Login.vue               # PIN 登录
│   │   ├── Register.vue            # PIN 注册
│   │   └── TwoFactorAuth.vue       # 2FA 验证
│   ├── stores/                     # 12 个 Pinia Store
│   │   ├── assetStore.js           # 资产状态
│   │   ├── accountStore.js         # 账户状态
│   │   ├── authStore.js            # 认证状态
│   │   ├── themeStore.js           # 主题 + 语言状态
│   │   ├── dashboardStore.js       # 仪表盘布局
│   │   ├── commandStore.js         # 命令面板
│   │   ├── notificationStore.js    # 通知
│   │   ├── goalStore.js            # 目标追踪
│   │   ├── nftStore.js             # NFT
│   │   ├── transactionStore.js     # 交易记录
│   │   ├── strategyStore.js        # 投资策略
│   │   └── achievementStore.js     # 成就系统
│   ├── composables/                # 8 个组合式函数
│   │   ├── useI18n.js              # 多语言
│   │   ├── useFormatters.js        # 格式化工具
│   │   ├── useToast.js             # Toast 提示
│   │   ├── useExportData.js        # 数据导出
│   │   ├── useHealthScore.js       # 健康评分
│   │   ├── useShareImage.js        # 分享图片
│   │   ├── useWalletConnect.js     # 钱包连接
│   │   └── usePullToRefresh.js     # 下拉刷新
│   ├── i18n/index.js               # 翻译定义（800+ key x 3 语言）
│   ├── data/                       # Mock 数据
│   │   ├── institutions.js         # 金融机构数据
│   │   └── mock*.js                # 各模块 Mock 数据
│   ├── router/index.js             # 路由配置
│   └── styles/main.css             # 全局样式 + CSS 变量
├── index.html
└── package.json                    # Vue 3.5 + Vite 7 + Pinia 3
```

---

## 4. 开发流程

### 4.1 功能开发

```bash
# 1. 从 master 创建功能分支
git checkout master && git pull origin master
git checkout -b feature/your-feature

# 2. 后端开发（如需要）
#    创建 handler → 注册路由 → 实现 service → 添加 repository

# 3. 前端开发（如需要）
#    创建 API 服务 → 创建/修改组件 → 集成到页面 → 添加 i18n 翻译

# 4. 测试
cd core && go test ./...
cd webapp && pnpm build  # 确保构建通过

# 5. 提交
git add <相关文件>
git commit -m "添加你的功能"
git push origin feature/your-feature
```

### 4.2 新增后端接口的标准步骤

1. **在 `service/interfaces.go` 中定义接口方法**
2. **在对应 service 文件中实现**
3. **在 `handlers/` 中创建 handler**
4. **在 `router.go` 的 `Setup()` 中注册路由**
5. **更新 `docs/swagger.yaml`**

### 4.3 新增前端页面的标准步骤

1. **在 `api/` 中创建 API 服务**
2. **在 `stores/` 中创建 Pinia Store（如需状态管理）**
3. **在 `components/` 中创建组件**
4. **在 `pages/` 中创建页面**
5. **在 `router/index.js` 中注册路由**
6. **在 `i18n/index.js` 中添加三语翻译**

---

## 5. 测试指南

### 5.1 后端测试

```bash
# 运行所有测试
cd core && go test ./...

# 显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

测试使用 `stretchr/testify` 断言库。

### 5.2 前端构建验证

```bash
cd webapp

# 构建（会检测编译错误）
pnpm build

# 预览构建结果
pnpm preview
```

---

## 6. 贡献指南

### 6.1 如何贡献

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/your-feature`
3. 开发功能（遵循代码规范）
4. 提交代码：`git commit -m "添加你的功能"`
5. 创建 Pull Request

### 6.2 代码审查标准

**必须满足**：
- 代码符合规范（Go 命名、Vue Composition API）
- 注释使用中文
- 所有测试通过
- 更新了相关文档
- 无安全隐患

**推荐**：
- 测试覆盖率 > 70%
- 性能无明显下降
- 遵循现有架构模式

### 6.3 添加新的交易所集成

1. 在 `core/internal/integrations/` 下创建目录
2. 实现 `GetBalances()` 和 `TestConnection()` 方法
3. 在 `ExchangeService` 中注册
4. 在 `models/base.go` 的 `ValidExchanges()` 中添加
5. 更新前端交易所选择列表

当前支持的交易所：**Binance**、**OKX**、**Coinbase**。

### 6.4 文档更新

**每次功能变更后必须同步更新**：
1. `README.md`（功能列表）
2. `docs/` 下对应文档
3. `i18n/index.js`（三语翻译）

---

## 相关文档

- [功能全景](../product/feature-overview.md)
- [用户指南](./user-guide.md)
- [API 接口文档](../tech/api-reference.md)
- [部署指南](./deployment-guide.md)
- [代理配置指南](./proxy-guide.md)
- [技术基线](../tech/tech-baseline.md)
- [前端需求规格](../specs/frontend-spec.md)
- [后端需求规格](../specs/backend-spec.md)

---

## 参考资源

- [Go 官方文档](https://golang.org/doc/)
- [Vue 3 官方文档](https://vuejs.org/)
- [Vite 文档](https://vitejs.dev/)
- [Tailwind CSS 文档](https://tailwindcss.com/)
- [Pinia 文档](https://pinia.vuejs.org/)
- [GoFrame 文档](https://goframe.org/)（配置管理部分）

---

**文档维护者**: @allfi
**最后更新**: 2026-02-13
