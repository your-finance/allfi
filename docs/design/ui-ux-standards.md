# AllFi UI/UX 设计规范

> 版本：v2.0 | 更新时间：2026-02-11

---

## 一、设计理念

采用**专业金融级设计语言**（Bloomberg / Wind / DeBank 风格）：

- **信息密集** — 高数据密度，紧凑布局，减少无意义留白
- **低装饰性** — 无玻璃拟态、无发光效果、无渐变文字
- **高对比度** — 深色背景 + 高亮数字，长时间阅读不疲劳
- **克制交互** — Hover 最多改变 2 个属性，无浮起动画

### 禁止的视觉效果

| 效果 | 说明 |
|------|------|
| `backdrop-filter: blur()` | 玻璃模糊效果 |
| `glow-card` / 发光边框 | 装饰性发光 |
| `gradient-text` | 渐变文字 |
| `translateY(-2px)` hover | 悬浮浮起动画 |
| 彩色 `box-shadow` | 彩色阴影光晕 |

---

## 二、主题系统

### 2.1 四套主题

| 主题 | ID | 模式 | 强调色 | 风格 |
|------|-----|------|--------|------|
| **Nexus Pro** | `nexus-pro` | 深色 | `#4B83F0` 金融蓝 | 默认，专业金融终端 |
| **Vestia** | `vestia` | 深色 | `#3B82D9` GitHub 蓝 | GitHub/开发者风格 |
| **XChange** | `xchange` | 深色 | `#3EA87A` 交易所绿 | 交易所风格 |
| **Aurora** | `aurora` | 浅色 | `#3574D4` 浅蓝 | 浅色模式 |

### 2.2 主题切换机制

配置文件：`webapp/src/themes/index.js`
状态管理：`webapp/src/stores/themeStore.js`

```javascript
// 主题应用流程
function applyTheme(themeId) {
  const cssVars = themeToCssVars(theme)
  // 将 27 个 CSS 变量注入 document.documentElement
  Object.entries(cssVars).forEach(([key, value]) => {
    document.documentElement.style.setProperty(key, value)
  })
  document.documentElement.setAttribute('data-theme', themeId)
  document.documentElement.setAttribute('data-theme-mode', theme.mode)
}
```

持久化 key：`allfi-theme`（localStorage）。

---

## 三、设计令牌 (Design Tokens)

### 3.1 颜色系统

所有颜色通过 CSS 变量定义，由当前主题动态注入。

#### 背景层级

| 变量 | 用途 | Nexus Pro 值 |
|------|------|-------------|
| `--color-bg-primary` | 主背景 | `#0D1117` |
| `--color-bg-secondary` | 卡片背景 | `#161B22` |
| `--color-bg-tertiary` | 悬浮层 / hover 背景 | `#21262D` |
| `--color-bg-elevated` | 最高层（dropdown, modal） | `#30363D` |

#### 强调色

| 变量 | 用途 |
|------|------|
| `--color-accent-primary` | 主强调色（按钮、链接、活动态） |
| `--color-accent-secondary` | 次强调色（按钮 hover、辅助标记） |
| `--color-accent-tertiary` | 第三强调色（图表配色等） |

#### 语义色

| 变量 | 用途 | 典型值 |
|------|------|--------|
| `--color-success` | 涨幅、盈利 | `#2EBD85` |
| `--color-error` | 跌幅、亏损 | `#E25C5C` |
| `--color-warning` | 警告、待处理 | `#D4A843` |
| `--color-info` | 中性信息 | `#4B83F0` |

#### 文字层级

| 变量 | 用途 |
|------|------|
| `--color-text-primary` | 主要文字（高对比） |
| `--color-text-secondary` | 次要文字 |
| `--color-text-muted` | 弱化文字（标签、描述） |
| `--color-text-inverse` | 反色文字（用于亮色按钮上） |

#### 边框

| 变量 | 用途 |
|------|------|
| `--color-border` | 常规边框 |
| `--color-border-hover` | 悬停边框 |
| `--color-border-active` | 焦点/活动边框 |

#### 图表

| 变量 | 用途 |
|------|------|
| `--color-chart-line` | 图表线条颜色 |
| `--color-chart-gradient` | 图表渐变填充 |

### 3.2 字体系统

| 用途 | CSS 变量 | 字体 |
|------|----------|------|
| 标题 | `--font-heading` | DM Sans |
| 正文 | `--font-body` | IBM Plex Sans |
| 数值 | `--font-mono` | IBM Plex Mono |

**规则**：
- 所有数字（资产价值、百分比、数量）必须使用 `.font-mono` 类
- 字重上限 600，禁止使用 700 及以上

### 3.3 间距系统

| 变量 | 值 | 用途 |
|------|-----|------|
| `--gap-xs` | `4px` | 最小间距（图标与文字） |
| `--gap-sm` | `8px` | 小间距（列表项内） |
| `--gap-md` | `12px` | 中间距（卡片内边距） |
| `--gap-lg` | `16px` | 大间距（区块间） |
| `--gap-xl` | `24px` | 特大间距（页面分区） |
| `--gap-2xl` | `32px` | 最大间距（页面顶部） |

### 3.4 圆角系统

| 变量 | 值 | 用途 |
|------|-----|------|
| `--radius-xs` | `2px` | 最小（badge 内标记） |
| `--radius-sm` | `4px` | 小元素（badge, tag） |
| `--radius-md` | `6px` | 常规（button, input, card） |
| `--radius-lg` | `8px` | 最大值（modal, 大卡片） |

**规则**：圆角上限 8px，禁止更大的圆角。

### 3.5 阴影系统

| 变量 | 值 | 用途 |
|------|-----|------|
| `--shadow-sm` | `0 1px 2px rgba(0,0,0,0.2)` | 轻微浮起 |
| `--shadow-md` | `0 4px 12px rgba(0,0,0,0.3)` | 卡片默认 |
| `--shadow-lg` | `0 8px 24px rgba(0,0,0,0.4)` | 弹窗、抽屉 |

仅使用黑色阴影，不使用彩色阴影。

### 3.6 动画时长

| 变量 | 值 | 用途 |
|------|-----|------|
| `--transition-fast` | `120ms` | hover、active |
| `--transition-base` | `180ms` | 颜色、透明度 |
| `--transition-slow` | `280ms` | 布局、展开 |

---

## 四、组件规范

### 4.1 按钮

5 种按钮变体：

| 类名 | 用途 | 样式 |
|------|------|------|
| `.btn-primary` | 主操作 | 实心强调色背景 + 白色文字 |
| `.btn-secondary` | 次操作 | tertiary 背景 + 边框 |
| `.btn-ghost` | 弱操作 | 透明背景 + 次要文字色 |
| `.btn-danger` | 危险操作 | error 色背景 |
| `.btn-success` | 确认操作 | success 色背景 |

```css
.btn {
  padding: 8px 16px;
  border-radius: var(--radius-md);
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-primary {
  background: var(--color-accent-primary);
  color: #FFFFFF;
}

/* Hover：仅改变背景色，无浮起 */
.btn-primary:hover:not(:disabled) {
  background: var(--color-accent-secondary);
}
```

### 4.2 卡片

```css
.glass-card {
  background: var(--color-bg-secondary);  /* 纯色，非渐变 */
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  /* 无 backdrop-filter */
  /* 无 ::before 高光线 */
  /* 无彩色阴影 */
}

/* 可交互卡片的 Hover */
.glass-card-hover:hover {
  border-color: var(--color-border-hover);
  background: var(--color-bg-tertiary);
}
```

### 4.3 统计卡片 (StatCard)

```
┌─────────────────────────┐
│ [图标]  标签            │
│         $123,456.78     │  ← .font-mono
│         +2.35% (24h)    │  ← .positive / .negative
└─────────────────────────┘
```

支持 5 种 variant：`default`、`success`、`warning`、`error`、`info`。

### 4.4 表格

```css
.holdings-table th {
  padding: 8px 12px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-size: 0.75rem;
  color: var(--color-text-muted);
  border-bottom: 1px solid var(--color-border);
}

.holdings-table tr:hover {
  background: rgba(100, 116, 139, 0.1);
}

.holdings-table td {
  padding: 8px 12px;
  color: var(--color-text-secondary);
}
```

### 4.5 徽章

```css
.badge-success { background: rgba(var(--color-success-rgb), 0.15); color: var(--color-success); }
.badge-warning { background: rgba(var(--color-warning-rgb), 0.15); color: var(--color-warning); }
.badge-error   { background: rgba(var(--color-error-rgb), 0.15);   color: var(--color-error); }
.badge-info    { background: rgba(var(--color-info-rgb), 0.15);    color: var(--color-info); }
```

---

## 五、图标

### 5.1 图标库

- **Phosphor Icons**（`@phosphor-icons/vue`）
- 统一使用 **duotone** weight
- 标准尺寸：`20px`（紧凑处）/ `24px`（常规）/ `32px`（大图标）

### 5.2 使用方式

```vue
<script setup>
import { PhChartPieSlice, PhWallet, PhGear } from '@phosphor-icons/vue'
</script>

<template>
  <PhChartPieSlice :size="20" weight="duotone" />
</template>
```

---

## 六、数据可视化

### 6.1 图表技术栈

- **Chart.js 4** + **vue-chartjs**
- 图表颜色动态绑定当前主题

### 6.2 配色方案

**折线图**（趋势）：
- 线条色：`themeStore.currentTheme.colors.accentPrimary`
- 填充色：accentPrimary 的 15% 不透明度渐变
- 线宽：`1.5px`，曲线张力 `0.4`

**饼图**（资产分布）：
| 分类 | 颜色来源 |
|------|---------|
| CEX | `accentPrimary` |
| Blockchain | `accentSecondary` |
| Manual | `accentTertiary` |
| DeFi | 固定 `#8B5CF6` |
| NFT | 固定 `#EC4899` |

**基准对比线**：
- vs BTC：`#F7931A`（Bitcoin 官方橙）
- vs ETH：`#627EEA`（Ethereum 官方紫）

### 6.3 图表类型

| 场景 | 图表类型 | 说明 |
|------|---------|------|
| 资产趋势 | Line/Area | 渐变填充 15% opacity |
| 资产分布 | Doughnut | 内径 60% |
| 基准对比 | 双 Y 轴 Line | 资产值 + 基准线 |

---

## 七、布局规范

### 7.1 Dashboard 布局

```
┌──────────────────────────────────────────────────────┐
│  摘要栏 (summary-bar)                                 │
│  [总资产 $xxx,xxx] [今日 PnL] [分类芯片] [分享] [预警] │
├──────────────────────────────────────────────────────┤
│  图表区 (charts-row) - 3:2 Grid                      │
│  ┌──────────────────┐ ┌────────────────┐             │
│  │  资产趋势图 60%   │ │ 分布饼图 40%   │             │
│  │  [计价] [基准]    │ │ [分类图例]     │             │
│  │  [时间范围选择]    │ │                │             │
│  └──────────────────┘ └────────────────┘             │
├──────────────────────────────────────────────────────┤
│  洞察行 (insight-row) - 1:1 Grid                     │
│  [健康评分卡片]     [目标追踪面板]                      │
├──────────────────────────────────────────────────────┤
│  可选功能行（由 Widget 配置控制显示/隐藏）              │
│  [DeFi 概览] [NFT 概览] [费用分析] [策略面板]          │
├──────────────────────────────────────────────────────┤
│  持仓明细 (holdings-panel)                             │
│  [视图切换: 平铺/分组] [搜索框]                        │
│  [资产] [来源] [价格] [24h] [余额] [价值] [占比]       │
│  [分页控件]                                           │
└──────────────────────────────────────────────────────┘
```

### 7.2 持仓表两种视图

**平铺视图**：所有资产扁平列表，可按各列排序，支持展开显示来源详情。

**分组视图**：按智能分类分组（弹药库 / 核心持仓 / 风险资产），可折叠。

### 7.3 响应式断点

| 断点 | 范围 | 布局变化 |
|------|------|---------|
| Desktop | ≥ 1025px | 3 列 Grid，完整表格列 |
| Tablet | 768–1024px | 2 列 Grid，图表满宽 |
| Mobile | < 768px | 1 列 Stack，隐藏来源列，底部导航栏 |
| 小屏 | < 480px | 隐藏绝对变化值和时间标签 |

```css
@media (max-width: 768px) {
  .summary-bar { flex-direction: column; }
  .charts-row { grid-template-columns: 1fr; }
  .btn { min-height: 44px; }         /* 触摸目标最小 44px */
}
```

---

## 八、可访问性

### 8.1 对比度

- 大文字（18px+）：最低 3:1
- 正文：最低 4.5:1
- 图标：最低 3:1

### 8.2 交互

- 所有交互元素支持 `:focus-visible`
- 表单元素关联 `<label>`
- 移动端触摸目标最小 44px
- Tab 顺序符合逻辑流程

### 8.3 隐私模式

- 快捷键 `Ctrl+H` 一键隐藏所有金额
- 所有数值显示为 `••••`
- 通过 `useFormatters.js` composable 统一处理

---

## 九、CSS 类速查

```css
/* 按钮 */
.btn  .btn-primary  .btn-secondary  .btn-ghost  .btn-danger  .btn-success

/* 徽章 */
.badge-success  .badge-warning  .badge-error  .badge-info

/* 卡片 */
.glass-card  .glass-card-hover

/* 字体 */
.font-mono  .font-heading  .font-body

/* 文字色 */
.text-primary  .text-secondary  .text-muted

/* 语义色 */
.positive  .negative

/* 布局 */
.summary-bar  .charts-row  .insight-row  .holdings-panel
```

---

## 十、组件清单

### 已实现组件（39 个）

**基础组件**：StatCard、CryptoIcon、ToastContainer、NotificationPanel、BottomNav

**对话框/抽屉**：AddAccountDialog、PriceAlertDialog、BatchImportDialog、AddGoalDialog、AddStrategyDialog、WalletConnectDialog、AssetDetailDrawer、PortfolioShareDialog

**数据展示**：GoalCard、HealthScoreCard、ShareCard、TransactionItem、TransactionTimeline、TransactionFilter、CalendarHeatmap

**DeFi/NFT**：DeFiOverview、DeFiPositionCard、NFTOverview、NFTGallery、NFTCard

**高级功能**：StrategyPanel、RebalanceView、FeeAnalytics、DashboardCustomizer、CommandPalette、OnboardingWizard、BenchmarkPanel、AnnualReport、AnnualReportShare

**成就系统**：AchievementPanel、AchievementBadge、AchievementUnlock

**其他**：WalletConnectButton、PullToRefresh

---

**文档维护者**: @allfi
**最后更新**: 2026-02-11
