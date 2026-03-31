# AllFi 前端规格

> 审计日期：2026-03-31
> 范围：描述当前已实现的前端，而不是未落地需求。

---

## 1. 技术栈

| 类别 | 当前实现 |
|------|---------|
| 框架 | Vue 3.5（Composition API + `<script setup>`） |
| 构建 | Vite 7.3 |
| 状态管理 | Pinia 3 |
| 路由 | Vue Router 4，`createWebHistory()` |
| 图表 | Chart.js 4 + vue-chartjs 5 |
| 图标 | `@phosphor-icons/vue` |
| 样式 | CSS 变量主题系统 + Tailwind CSS 4 引入 |
| 工具库 | `@vueuse/core`、`html2canvas`、`qrcode` |
| PWA | `vite-plugin-pwa` |
| 测试 | Vitest |

---

## 2. 当前目录结构

```text
webapp/src/
├── api/           15 个 API 文件（含 client.js / index.js / 分模块 service）
├── components/    58 个 Vue 组件
├── composables/   9 个组合式函数
├── data/          mock 数据
├── i18n/          翻译字典与翻译函数
├── pages/         12 个页面
├── router/        路由与守卫
├── stores/        13 个 Pinia Store
├── styles/        全局样式与设计令牌
├── themes/        4 套主题配置
├── App.vue
└── main.js
```

### 关键更正

- 当前不存在 `socialService.js`
- 当前不存在 `useTheme.js`
- 主题与语言统一由 `themeStore.js` 管理
- 组件实际数量为 58，不是 57

---

## 3. 页面与路由

### 认证页

| 路径 | 页面 |
|------|------|
| `/login` | 登录 |
| `/register` | 首次注册 / 设置密码 |
| `/2fa` | 双因素认证 |

### 应用页

| 路径 | 页面 |
|------|------|
| `/dashboard` | 仪表盘 |
| `/accounts` | 账户管理 |
| `/history` | 历史记录 |
| `/analytics` | 数据分析 |
| `/reports` | 报告 |
| `/risk` | 风险 |
| `/settings` | 设置 |
| `/defi` | DeFi 专题页 |
| `/nft` | NFT 专题页 |

### 路由守卫

- 认证状态由 `authStore` 管理
- 已登录用户访问 `/login` 或 `/register` 会被重定向到 `/dashboard`
- `requires2FA` 为真时，仅允许访问 `/2fa`
- `/swagger` 与 `/api.json` 在开发模式下会被直接透传到后端

---

## 4. 状态管理

### Store 列表（13 个）

- `accountStore`
- `achievementStore`
- `assetStore`
- `authStore`
- `commandStore`
- `dashboardStore`
- `goalStore`
- `nftStore`
- `notificationStore`
- `strategyStore`
- `systemStore`
- `themeStore`
- `transactionStore`

### 核心职责

| Store | 说明 |
|------|------|
| `authStore` | 密码模式、JWT、2FA 流程 |
| `themeStore` | 主题、语言、隐私模式、引导状态 |
| `systemStore` | 版本信息、更新检查、回滚目标 |
| `accountStore` | CEX / 钱包 / 手动资产入口 |
| `assetStore` | 聚合资产摘要与分布 |
| `transactionStore` | 交易记录与筛选 |

---

## 5. 接口层

### 接口模式

- 默认通过 `client.js` + `index.js` 统一封装请求
- 支持 mock / real 双路径
- 真实接口走 `/api/v1`
- 开发模式下由 Vite 代理到 `http://127.0.0.1:8080`

### 主要服务文件

- `achievementService.js`
- `annualReportService.js`
- `benchmarkService.js`
- `cross-chain.js`
- `defiService.js`
- `feeService.js`
- `marketService.js`
- `nftService.js`
- `riskService.js`
- `strategyService.js`
- `systemService.js`
- `transactionService.js`

---

## 6. 认证与安全流

### 密码模式

- PIN：4-20 位数字
- 复杂密码：8-20 位，必须包含大小写字母和数字

### 登录返回值

`authStore.login()` 当前返回：

```js
{ success: boolean, requires2FA: boolean }
```

### 会话恢复

- 前端将 `token`、`isAuthenticated`、`requires2FA` 写入 `localStorage`
- 路由守卫进入受保护页面前会尝试恢复会话

---

## 7. UI 与交互基线

- 主题：4 套
- 语言：3 套
- 隐私模式：支持
- PWA：支持
- Web Push：支持浏览器订阅
- 图像分享：年度报告 / 资产分享图已实现
- 移动端：底部导航、下拉刷新、响应式布局已实现

---

## 8. 测试现状

执行命令：

```bash
cd webapp && pnpm test --run
```

当前结果：

- 27 个测试全部通过
