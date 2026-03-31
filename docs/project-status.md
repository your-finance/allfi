# AllFi 项目现状

> 审计日期：2026-03-31
> 结论：项目主干能力完整，但文档与当前实现存在明显漂移，测试也并非全绿。

---

## 1. 仓库现状概览

### 前端

| 项目 | 当前实现 |
|------|---------|
| 页面 | 12 个 `webapp/src/pages/*.vue` |
| 组件 | 58 个 `webapp/src/components/**/*.vue` |
| Store | 13 个 `webapp/src/stores/*.js` |
| Composable | 9 个 `webapp/src/composables/*.js` |
| API 文件 | 15 个 `webapp/src/api/*.js` |
| 路由模式 | `createWebHistory()` |
| 运行端口 | 开发环境默认 `3000` |

### 后端

| 项目 | 当前实现 |
|------|---------|
| Go 版本 | `go 1.24.11` |
| 框架 | GoFrame v2.10（`ghttp` / `gdb` / `gcfg` / `goai`） |
| 业务模块 | 29 个 `core/internal/app/*` |
| 定时任务 | 10 个 `core/internal/cron/*.go` |
| API 定义 | 78 个 `core/api/v1/**/*.go` 中的 `g.Meta path` |
| 默认数据库 | SQLite（`sqlite::@file(./data/allfi.db)`） |
| OpenAPI | `/api.json` |
| Swagger UI | `/swagger/` |

---

## 2. 当前部署形态

### A. 仓库内开发模式

- `make dev` 同时启动前端和后端。
- 前端开发地址：`http://localhost:3000`
- 后端 API 地址：`http://localhost:8080/api/v1`
- Swagger 地址：`http://localhost:8080/swagger/`

### B. 仓库根目录 Docker Compose

- 文件：`docker-compose.yml`
- 组成：`backend` + `updater`
- 现状：`backend` 默认引用本地镜像 `allfi-backend:latest`，不是开箱即用的源码构建。
- 含义：更适合维护者或已经手动构建过镜像的场景。

### C. 一键部署脚本

- 文件：`deploy/docker-deploy.sh`
- 用途：面向最终用户，下载 GitHub Release 二进制并生成独立部署目录。
- 默认访问端口：`3000`
- 现状：脚本生成的是最小化部署目录，不包含仓库根目录里的 `updater` sidecar。

### D. 宿主机二进制运行

- 发布产物来自 `.github/workflows/release.yml`
- 宿主机模式下，后端 `system` 模块已实现 OTA 二进制替换与重启逻辑。

---

## 3. 版本更新能力现状

| 运行模式 | 当前实现 |
|---------|---------|
| 宿主机二进制 | 通过 GitHub Release 下载对应平台 tarball，并使用 `go-update` 原位替换 |
| 仓库根目录 Docker Compose | 通过 `updater` sidecar 执行 `git fetch` / `git checkout` / `docker build` / `docker compose up` |
| 一键部署脚本生成的最小部署 | 当前不包含 `updater` sidecar，不应宣称支持页面内一键更新 |

---

## 4. 认证与安全现状

- 支持两种密码模式：
  - PIN：4-20 位数字
  - 复杂密码：8-20 位，必须包含大小写字母和数字
- 支持 2FA / TOTP
- JWT 用于会话认证
- API Key 本地 AES-256-GCM 加密存储
- 前端登录流程已经返回对象结果：`{ success, requires2FA }`

---

## 5. 与旧文档的主要偏差

- Swagger 入口不是 `/api/v1/docs`，而是 `/swagger/`，OpenAPI JSON 是 `/api.json`。
- 后端不是 GORM / `net/http` 自研路由结构，而是 GoFrame `ghttp` + `gdb` + 代码生成 DAO。
- 根目录 Docker Compose 不是“直接 `docker compose up -d` 就能跑”的纯源码构建模式。
- 一键部署脚本当前不会生成 `updater` sidecar，因此不能继续在文档里宣称页面内 OTA 更新可用。
- 前端组件数已是 58，不再是 57。
- 后端业务模块是 29 个，不再是 26 个。

---

## 6. 测试状态

### 后端

执行命令：

```bash
cd core && go test ./... -timeout 60s
```

当前结果：

- 默认测试通过。
- `asset/logic` 与 `manual_asset/logic` 的历史签名漂移测试已修复。
- `exchange_rate/provider` 中依赖公网的用例已改为显式 opt-in。
- 如需执行外网集成测试，请设置：

```bash
ALLFI_RUN_ONLINE_TESTS=1 go test ./internal/app/exchange_rate/provider -timeout 60s
```

### 前端

执行命令：

```bash
cd webapp && pnpm test --run
```

当前结果：

- 27 个测试全部通过。
- `authStore.login()` 的测试断言已同步到当前返回结构 `{ success, requires2FA }`。

---

## 7. 本轮文档整理原则

- 保留“当前实现说明”和“长期有效设计原则”。
- 删除已经完成且内容明显过时的计划文档。
- 将部署、架构、端口、Swagger、更新机制统一到当前实现。
