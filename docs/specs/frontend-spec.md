# AllFi 前端需求规格

> 版本：v2.1 | 更新时间：2026-02-28

---

## 1. 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | Vue 3（Composition API + `<script setup>`） | ^3.5 |
| 构建工具 | Vite | ^7.3 |
| 状态管理 | Pinia（Composition API 风格） | ^3.0 |
| 路由 | Vue Router | ^4.6 |
| 图表 | Chart.js + vue-chartjs | ^4.5 / ^5.3 |
| 图标 | Phosphor Icons（`@phosphor-icons/vue`） | ^2.2 |
| CSS | CSS 变量 + Tailwind CSS | ^4.1 |
| 工具库 | @vueuse/core | ^14.1 |
| PWA | vite-plugin-pwa | ^1.2 |
| 截图 | html2canvas | ^1.4 |
| 包管理 | pnpm | - |

**不使用**：Axios（用原生 Fetch）、TypeScript（纯 JavaScript）、Lucide（用 Phosphor）、ESLint。

---

## 2. 项目结构

```
webapp/src/
├── api/                    # API 服务模块（原生 Fetch）
│   ├── index.js            # 基础 HTTP 客户端 + 请求封装
│   ├── achievementService.js
│   ├── annualReportService.js
│   ├── benchmarkService.js
│   ├── defiService.js
│   ├── feeService.js
│   ├── marketService.js
│   ├── nftService.js
│   ├── riskService.js
│   ├── socialService.js
│   ├── strategyService.js
│   ├── systemService.js
│   ├── transactionService.js
│   └── cross-chain.js      # 跨链交易服务
├── components/             # 57 个可复用组件
├── composables/            # 9 个组合式函数
│   ├── useI18n.js          # 多语言翻译
│   ├── useToast.js         # 消息通知
│   ├── useFormatters.js    # 格式化（日期/数字/货币）
│   ├── useHealthScore.js   # 健康评分计算
│   ├── useShareImage.js    # 分享图片生成
│   ├── useExportData.js    # 数据导出
│   ├── useWalletConnect.js # WalletConnect 集成
│   ├── usePullToRefresh.js # 下拉刷新
│   └── useTheme.js         # 主题切换
├── data/                   # Mock 数据（开发/演示用）
├── i18n/
│   └── index.js            # 3 语言翻译（zh-CN/zh-TW/en-US），900+ key
├── pages/                  # 12 个页面组件
├── router/
│   └── index.js            # 路由配置 + 认证守卫
├── stores/                 # 13 个 Pinia Store
├── styles/
│   └── main.css            # 全局样式 + 设计令牌
├── themes/
│   └── index.js            # 4 套主题配置（27 个 CSS 变量）
├── App.vue                 # 根组件（侧边栏 + 顶栏布局）
└── main.js                 # 入口文件
```

---

## 3. 路由定义

```javascript
// src/router/index.js
const routes = [
  // 认证页面（无需登录）
  { path: '/login',    name: 'Login',         component: Login,         meta: { requiresAuth: false } },
  { path: '/register', name: 'Register',      component: Register,      meta: { requiresAuth: false } },
  { path: '/2fa',      name: 'TwoFactorAuth', component: TwoFactorAuth, meta: { requiresAuth: false } },

  // 核心应用页面（需登录）
  { path: '/',          redirect: '/dashboard' },
  { path: '/dashboard', name: 'Dashboard', component: Dashboard, meta: { requiresAuth: true, titleKey: 'nav.dashboard' } },
  { path: '/accounts',  name: 'Accounts',  component: Accounts,  meta: { requiresAuth: true, titleKey: 'nav.accounts' } },
  { path: '/history',   name: 'History',   component: History,   meta: { requiresAuth: true, titleKey: 'nav.history' } },
  { path: '/analytics', name: 'Analytics', component: Analytics, meta: { requiresAuth: true, titleKey: 'nav.analytics' } },
  { path: '/reports',   name: 'Reports',   component: Reports,   meta: { requiresAuth: true, titleKey: 'nav.reports' } },
  { path: '/risk',      name: 'Risk',      component: Risk,      meta: { requiresAuth: true, titleKey: 'nav.risk' } },
  { path: '/settings',  name: 'Settings',  component: Settings,  meta: { requiresAuth: true, titleKey: 'nav.settings' } },

  // DeFi 和 NFT 专题页面（需登录）
  { path: '/defi',      name: 'DeFi',      component: DeFi,      meta: { requiresAuth: true, titleKey: 'nav.defi' } },
  { path: '/nft',       name: 'NFT',       component: NFT,       meta: { requiresAuth: true, titleKey: 'nav.nft' } },

  // 404 重定向
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' }
]
```

路由守卫：`beforeEach` 检查 `authStore.isLoggedIn`，未登录重定向到 `/login`。

---

## 4. 页面需求

### 4.1 登录页（Login.vue）

**认证方式：灵活密码模式**

支持两种密码模式，用户可在设置页面随时切换：
- **PIN 码模式**：4-20 位纯数字，显示格子输入界面
- **复杂密码模式**：8-20 位，必须包含大小写字母和数字，显示普通密码输入框

**首次设置流程**：
1. 系统检测未设置密码，显示设置界面
2. 用户可选择 PIN 或复杂密码模式
3. 根据选择显示对应的输入界面
4. 设置成功后自动获取 JWT Token 并登录

**登录流程**：
- 已有密码登录，返回 JWT Token
- Token 存储在 `localStorage['allfi-auth']`
- 连续 5 次失败锁定 15 分钟
- 若启用 2FA，登录后需验证 2FA 码

```javascript
// authStore 核心方法
setupPIN(pin)              // 首次设置（自动检测类型）
login(password)            // 验证登录
changePassword(old, new)   // 修改密码
switchPasswordType(currentPassword, newType, newPassword, twoFACode)  // 切换密码类型
logout()                   // 登出
restoreSession()           // 恢复会话
```

### 4.2 仪表盘（Dashboard.vue）

| 区域 | 组件 | 功能 |
|------|------|------|
| 摘要栏 | summary-bar | 总资产值 + 今日 PnL + 分类芯片 + 分享/预警按钮 |
| 图表区 | charts-row（3:2 Grid） | 左：资产趋势图（计价/基准/时间范围） 右：资产分布饼图 |
| 洞察行 | insight-row（1:1 Grid） | 健康评分卡片 + 目标追踪面板 |
| 可选功能行 | Widget 控制 | DeFi 概览 / NFT 概览 / 费用分析 |
| 持仓明细 | holdings-panel | 平铺/分组视图 + 搜索 + 排序 + 分页 |

**数据接口**：
- `GET /assets/summary` — 总资产 + 24h 变化 + 分类占比
- `GET /assets/history?range=30d` — 历史趋势
- `GET /assets/details` — 资产明细列表
- `GET /portfolio/health` — 健康评分

**交互特性**：
- 4 种计价货币切换（USDC/BTC/ETH/CNY）
- 基准对比（vs BTC/vs ETH）
- 时间范围（7D/30D/90D/1Y/ALL）
- Widget 配置（DashboardCustomizer 组件）
- 隐私模式（Ctrl+H 隐藏金额）

### 4.3 账户管理（Accounts.vue）

| Tab | 功能 | 接口 |
|-----|------|------|
| CEX 交易所 | 添加/编辑/删除/测试连接/刷新余额 | `/exchanges/accounts` |
| 区块链钱包 | 添加/删除/批量导入/查看余额 | `/wallets/addresses` |
| 手动资产 | 4 类（银行/现金/股票/基金）CRUD | `/assets/manual` |

**对话框组件**：AddAccountDialog、BatchImportDialog、AssetDetailDrawer

### 4.4 历史记录（History.vue）

- 资产趋势图（多时间范围）
- 日历热力图（CalendarHeatmap，GitHub 风格）
- 交易记录时间线（TransactionTimeline + TransactionFilter）
- 支持 CEX + 链上交易的统一视图

### 4.5 数据分析（Analytics.vue）

| 功能 | 说明 | 接口 |
|------|------|------|
| 每日盈亏 | PnL 折线图 + 摘要统计 | `/analytics/pnl/daily` |
| 归因分析 | 价格效应 vs 数量效应 | `/analytics/attribution` |
| 趋势预测 | 线性回归 + R² 置信度 | `/analytics/forecast` |
| 基准对比 | vs BTC/ETH/S&P500 | `/benchmark` |
| 集中度分析 | HHI 指数 + Top 5 占比 | 本地计算 |

### 4.6 报告（Reports.vue）

- 日报/周报/月报/年报自动生成
- 年度报告分享（AnnualReport + AnnualReportShare）
- 报告列表 + 详情查看

### 4.7 设置（Settings.vue）

| 分类 | 设置项 |
|------|--------|
| 通用 | 默认计价货币、语言（3 种）、主题（4 套） |
| 刷新 | 自动刷新间隔、价格缓存 TTL |
| 安全 | 密码管理（修改密码、切换密码类型）、2FA 设置 |
| 通知 | 通知偏好、WebPush 订阅 |
| 数据 | 导出数据、清除缓存、重置设置 |
| 交易同步 | 同步设置、增量同步触发 |
| 关于 | 版本号、GitHub 链接 |

**密码类型切换功能**：
- 支持在 PIN 码和复杂密码之间切换
- 切换流程：输入当前密码 → 选择新类型 → 输入新密码 → （若启用 2FA）验证 2FA
- 前端根据密码类型动态显示输入界面

---

## 5. 状态管理

12 个 Pinia Store，全部使用 Composition API 风格（`setup` 函数）：

| Store | 文件 | 职责 |
|-------|------|------|
| `useAccountStore` | accountStore.js | CEX 账户 + 钱包地址 + 手动资产 CRUD |
| `useAssetStore` | assetStore.js | 资产总览 + 详情 + 历史 + 刷新 |
| `useAuthStore` | authStore.js | PIN 认证 + JWT Token + 会话管理 |
| `useThemeStore` | themeStore.js | 主题切换 + 语言切换 + 隐私模式 + 引导状态 |
| `useGoalStore` | goalStore.js | 目标追踪 CRUD |
| `useNotificationStore` | notificationStore.js | 通知列表 + 未读计数 + 偏好 |
| `useAchievementStore` | achievementStore.js | 成就列表 + 解锁检查 |
| `useNftStore` | nftStore.js | NFT 资产列表 |

| `useTransactionStore` | transactionStore.js | 交易记录 + 同步 + 统计 |
| `useDashboardStore` | dashboardStore.js | Widget 配置 + 布局偏好 |
| `useCommandStore` | commandStore.js | 命令面板状态 |

---

## 6. API 集成

### 6.1 HTTP 客户端

使用原生 `Fetch API`，封装在 `api/index.js`：

- 基础 URL：`http://localhost:8080/api/v1`（可通过 `VITE_API_BASE_URL` 覆盖）
- 超时：30 秒
- 认证：自动附加 `Authorization: Bearer {token}`
- 401 响应自动清除 Token 并重定向登录

### 6.2 Mock 模式

- 环境变量 `VITE_USE_MOCK_API` 控制（默认启用）
- 后端不可用时自动返回模拟数据
- 模拟网络延迟（200-2000ms）

### 6.3 API 服务模块

| 模块 | 文件 | 端点前缀 |
|------|------|---------|
| 认证 | api/index.js | `/auth/*` |
| 资产 | api/index.js | `/assets/*` |
| 交易所 | api/index.js | `/exchanges/*` |
| 钱包 | api/index.js | `/wallets/*` |
| 汇率 | api/index.js | `/rates/*` |
| 设置 | api/index.js | `/users/*` |
| 通知 | api/index.js | `/notifications/*` |
| 预警 | api/index.js | `/alerts/*` |
| 报告 | api/index.js | `/reports/*` |
| 分析 | api/index.js | `/analytics/*` |
| DeFi | defiService.js | `/defi/*` |
| NFT | nftService.js | `/nft/*` |
| 交易记录 | transactionService.js | `/transactions/*` |
| 费用 | feeService.js | `/analytics/fees` |

| 成就 | achievementService.js | `/achievements/*` |

---

## 7. 组件清单（57 个）

### 7.1 基础组件（7 个）

| 组件 | 功能 |
|------|------|
| StatCard.vue | 统计卡片（多种 variant） |
| CryptoIcon.vue | 加密资产图标 |
| ToastContainer.vue | Toast 消息容器 |
| NotificationPanel.vue | 通知面板 |
| BottomNav.vue | 移动端底部导航 |
| PullToRefresh.vue | 下拉刷新 |
| VersionBadge.vue | 版本徽章 |

### 7.2 对话框/抽屉（8 个）

| 组件 | 功能 |
|------|------|
| AddAccountDialog.vue | 添加 CEX/钱包/手动资产 |
| AssetDetailDrawer.vue | 资产详情侧边抽屉 |
| PriceAlertDialog.vue | 价格预警设置 |
| BatchImportDialog.vue | 批量导入钱包地址 |
| AddGoalDialog.vue | 添加投资目标 |
| AddStrategyDialog.vue | 添加策略 |
| WalletConnectDialog.vue | WalletConnect 连接 |
| PortfolioShareDialog.vue | 投资组合分享 |

### 7.3 数据展示（5 个）

| 组件 | 功能 |
|------|------|
| GoalCard.vue | 目标进度卡片 |
| HealthScoreCard.vue | 健康评分卡片 |
| ShareCard.vue | 分享卡片 |
| CalendarHeatmap.vue | 日历热力图 |
| BenchmarkPanel.vue | 基准对比面板 |

### 7.4 交易记录（3 个）

| 组件 | 功能 |
|------|------|
| TransactionTimeline.vue | 交易时间线 |
| TransactionItem.vue | 单条交易项 |
| TransactionFilter.vue | 交易筛选器 |

### 7.5 DeFi / NFT（11 个）

| 组件 | 功能 |
|------|------|
| DeFiOverview.vue | DeFi 仓位概览 |
| DeFiPositionCard.vue | 单个 DeFi 仓位卡片 |
| DeFiMiniCard.vue | DeFi 迷你卡片 |
| LendingPositionCard.vue | 借贷仓位卡片 |
| LendingOptimizer.vue | 借贷优化器 |
| LendingRateChart.vue | 借贷利率图表 |
| HealthFactorGauge.vue | 健康因子仪表盘 |
| NFTOverview.vue | NFT 资产概览 |
| NFTGallery.vue | NFT 画廊 |
| NFTCard.vue | 单个 NFT 卡片 |
| NFTMiniCard.vue | NFT 迷你卡片 |

### 7.6 分析与风险（7 个）

| 组件 | 功能 |
|------|------|
| FeeAnalytics.vue | 费用分析 |
| AttributionPanel.vue | 盈亏归因分析 |
| DrawdownChart.vue | 回撤图表 |
| ForecastPanel.vue | 收益预测 |
| RiskAlertPanel.vue | 风险预警面板 |
| RiskMetricsChart.vue | 风险指标图表 |
| RiskOverviewCard.vue | 风险概览卡片 |
| BetaComparisonCard.vue | Beta 对比卡片 |

### 7.7 高级功能（8 个）

| 组件 | 功能 |
|------|------|
| RebalanceView.vue | 再平衡视图 |
| DashboardCustomizer.vue | Widget 配置 |
| CommandPalette.vue | 命令面板（Cmd+K） |
| OnboardingWizard.vue | 首次使用引导 |
| AnnualReport.vue | 年度报告 |
| AnnualReportShare.vue | 年度报告分享 |
| WalletConnectButton.vue | 钱包连接按钮 |
| StrategyPanel.vue | 策略面板 |

### 7.8 成就系统（3 个）

| 组件 | 功能 |
|------|------|
| AchievementPanel.vue | 成就面板 |
| AchievementBadge.vue | 成就徽章 |
| AchievementUnlock.vue | 成就解锁动画 |

### 7.9 新增组件（5 个）

| 组件 | 功能 |
|------|------|
| AttributionPanel.vue | 盈亏归因分析面板 |
| DrawdownChart.vue | 回撤图表 |
| ForecastPanel.vue | 收益预测面板 |
| RiskAlertPanel.vue | 风险预警面板 |
| RiskMetricsChart.vue | 风险指标图表 |

---

## 8. PWA 配置

```javascript
// vite.config.js
VitePWA({
  registerType: 'autoUpdate',
  manifest: {
    name: 'AllFi - 全资产聚合平台',
    short_name: 'AllFi',
    theme_color: '#0C0E12',
    display: 'standalone'
  },
  workbox: {
    // 预缓存：JS/CSS/HTML/SVG/PNG/ICO/WOFF2
    // 运行时缓存：Google Fonts（1 年） + API 缓存（24 小时）
  }
})
```

开发服务器端口：`3173`，API 代理 `/api` → `http://localhost:8080`。

---

## 9. 响应式设计

| 断点 | 范围 | 布局变化 |
|------|------|---------|
| Desktop | ≥ 1025px | 侧边栏 + 3 列 Grid |
| Tablet | 768–1024px | 折叠侧边栏 + 2 列 Grid |
| Mobile | < 768px | 底部导航 + 1 列 Stack + 隐藏部分表格列 |
| 小屏 | < 480px | 进一步简化显示 |

触摸目标最小 44px。

---

## 10. 国际化

- 3 种语言：简体中文（zh-CN，默认）/ 繁體中文（zh-TW）/ English（en-US）
- 32 个翻译分区，900+ key
- 使用 `useI18n()` composable，详见 [docs/design/i18n.md](../design/i18n.md)

---

## 11. Store 清单（13 个）

| Store | 功能 |
|-------|------|
| authStore | 认证状态（PIN/2FA） |
| assetStore | 资产总览与明细 |
| accountStore | CEX/钱包/手动资产账户 |
| transactionStore | 统一交易记录 |
| nftStore | NFT 资产 |
| goalStore | 投资目标 |
| achievementStore | 成就系统 |
| notificationStore | 通知管理 |
| strategyStore | 策略引擎 |
| systemStore | 系统配置 |
| themeStore | 主题切换 |
| dashboardStore | Dashboard 数据聚合 |
| commandStore | 命令面板（Cmd+K） |

---

**文档维护者**: @allfi
**最后更新**: 2026-02-28
