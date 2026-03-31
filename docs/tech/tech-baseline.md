# AllFi 技术基线

> 审计日期：2026-03-31
> 目标：描述当前仓库真实实现，而不是历史方案。

---

## 1. 架构摘要

```text
浏览器 / PWA
    │
    ├── 开发模式：Vite dev server (:3000) 代理后端 API (:8080)
    └── 发布模式：后端内嵌前端静态资源，单端口提供页面与 API

前端
    Vue 3 + Vite 7 + Pinia + Vue Router + Chart.js + Tailwind CSS 4
    12 页面 / 58 组件 / 13 Store / 9 Composable

后端
    Go 1.24.11 + GoFrame v2.10
    controller -> logic -> service -> dao/model
    29 个业务模块 / 10 个 cron / 78 个 API 定义

存储
    SQLite 默认配置（GoFrame gdb）
    启动时由 `core/internal/database/migrate.go` 执行建表与补充迁移

外部集成
    Binance / OKX / Coinbase
    Etherscan 系列 / Alchemy / CoinGecko / Yahoo / Frankfurter / GateIO
    DeFi: Lido / Rocket Pool / Aave / Compound / Uniswap V2/V3 / Curve
```

---

## 2. 前端基线

| 维度 | 当前实现 |
|------|---------|
| 框架 | Vue 3.5（Composition API + `<script setup>`） |
| 构建 | Vite 7.3 |
| 状态管理 | Pinia 3 |
| 路由 | Vue Router 4，`createWebHistory()` |
| UI | CSS 变量主题系统 + `@import "tailwindcss"` |
| 主题 | 4 套主题，主题与语言统一由 `themeStore` 管理 |
| 国际化 | `zh-CN` / `zh-TW` / `en-US` |
| 离线能力 | `vite-plugin-pwa` + `push-sw.js` |
| 图表 | Chart.js + vue-chartjs |
| 接口层 | 原生 `fetch`，支持 mock / real 双路径 |

### 当前目录统计

- 页面：12 个
- 组件：58 个
- Store：13 个
- Composable：9 个
- API 文件：15 个

### 关键实现约定

- 开发模式下前端运行在 `http://localhost:3000`
- `/api`、`/swagger`、`/api.json` 会代理到后端 `:8080`
- 主题、语言、隐私模式、引导状态由 `webapp/src/stores/themeStore.js` 统一管理
- 登录流程支持 PIN / 复杂密码 + 可选 2FA

---

## 3. 后端基线

| 维度 | 当前实现 |
|------|---------|
| 语言 | Go 1.24.11 |
| 框架 | GoFrame v2.10 |
| Web 层 | `ghttp.Server` + 路由分组 + OpenAPI |
| 配置 | `gcfg` + YAML |
| 数据访问 | GoFrame `gdb` + 代码生成 DAO |
| 文档 | `goai` 生成 OpenAPI，Swagger UI 暴露于 `/swagger/` |
| 存储 | SQLite 默认配置 |
| 认证 | bcrypt + JWT + 2FA/TOTP |
| 加密 | AES-256-GCM |
| 更新 | 宿主机 OTA / Docker sidecar 更新 |

### 模块与任务

- 业务模块：29 个 `core/internal/app/*`
- 定时任务：10 个 `core/internal/cron/*.go`
- API 定义：78 个 `g.Meta path`

### 分层方式

- `core/api/v1/`：请求与响应结构、OpenAPI 元数据
- `core/internal/app/<module>/controller`：路由注册与参数接收
- `core/internal/app/<module>/logic`：业务逻辑
- `core/internal/app/<module>/service`：服务接口
- `core/internal/app/<module>/dao|model`：DAO、DO、Entity

### 数据层现状

- 仓库默认与示例配置都使用 SQLite
- `core/internal/database/migrate.go` 负责启动建表
- 仓库中仍保留部分 MySQL 辅助配置代码，但并不是当前默认文档化主路径

---

## 4. 部署与运行基线

### 开发模式

- `make dev`
- 前端：`http://localhost:3000`
- 后端：`http://localhost:8080`

### 发布模式

- 前端 `dist` 会复制到 `core/internal/statics/dist`
- 后端通过 `embed.FS` 提供静态页面
- 对用户暴露单端口页面 + API

### Docker 模式

存在两条不同路径：

1. 仓库根目录 `docker-compose.yml`
   - 面向维护者
   - 依赖本地镜像 `allfi-backend:latest`
   - 带 `updater` sidecar

2. `deploy/docker-deploy.sh`
   - 面向最终用户
   - 下载 Release 二进制并生成独立部署目录
   - 当前生成的是最小化部署，不带 `updater` sidecar

---

## 5. 版本与更新基线

| 场景 | 当前实现 |
|------|---------|
| 版本注入 | 通过构建参数写入 `core/internal/version` 与前端 `__APP_VERSION__` |
| Release | `.github/workflows/release.yml` 构建多平台二进制 |
| 宿主机更新 | 下载 GitHub Release tarball，`go-update` 原位替换 |
| Docker 更新 | `updater` sidecar 走 git checkout + docker build + restart |

---

## 6. 测试与已知注意事项

- 文档曾长期把 Swagger 路径写成 `/api/v1/docs`，当前应统一为 `/swagger/`。
- 部分旧文档把后端写成 `net/http + GORM`，与当前 GoFrame 实现不符。
- 默认前后端测试当前可通过：
  - `cd core && go test ./... -timeout 60s`
  - `cd webapp && pnpm test --run`
- `exchange_rate/provider` 中依赖外网的 Provider 测试默认跳过，需显式设置 `ALLFI_RUN_ONLINE_TESTS=1` 才执行。
