# AllFi 开发指南

> 审计日期：2026-03-31
> 范围：本地开发、构建、测试与当前注意事项。

---

## 1. 开发环境

### 必需工具

```bash
go version
node --version
pnpm --version
git --version
```

推荐版本：

- Go 1.24.x
- Node.js 20+
- pnpm 最新稳定版

---

## 2. 快速开始

### 方式一：Makefile

```bash
make setup
make dev
```

### 方式二：脚本

```bash
bash scripts/quickstart.sh
bash scripts/quickstart.sh --mock
bash scripts/quickstart.sh --check
```

### 启动后的地址

- 前端：`http://localhost:3000`
- 后端：`http://localhost:8080`
- Swagger：`http://localhost:8080/swagger/`

---

## 3. 常用命令

| 命令 | 说明 |
|------|------|
| `make help` | 查看全部命令 |
| `make setup` | 生成 `.env` 并安装依赖 |
| `make dev` | 同时启动前后端 |
| `make dev-mock` | 仅前端 mock |
| `make dev-backend` | 仅后端 |
| `make dev-frontend` | 仅前端 |
| `make build` | 构建前后端 |
| `make health` | 健康检查 |
| `make swagger` | 打开 Swagger |
| `make clean` | 清理构建产物 |

### 关于 `make docker`

- 当前根目录 `docker-compose.yml` 默认依赖本地镜像 `allfi-backend:latest`
- 因此 `make docker` 更适合维护者，而不是“拿源码就能直接起”

---

## 4. 手动启动

### 后端

```bash
cd core
go mod download && go mod verify
go run cmd/server/main.go
```

### 前端

```bash
cd webapp
pnpm install
pnpm dev
```

### Mock 模式

```bash
cd webapp
pnpm dev:mock
```

---

## 5. 当前目录认知

### 后端

- `core/api/v1/`：接口定义与 OpenAPI 元数据
- `core/internal/app/`：29 个业务模块
- `core/internal/cron/`：10 个定时任务
- `core/internal/database/`：建表与迁移
- `core/internal/integrations/`：外部数据接入

### 前端

- `webapp/src/pages/`：12 个页面
- `webapp/src/components/`：58 个组件
- `webapp/src/stores/`：13 个 Store
- `webapp/src/composables/`：9 个组合式函数
- `webapp/src/api/`：mock / real API 双路径封装

---

## 6. 测试

### 前端测试

```bash
cd webapp && pnpm test --run
```

当前状态：

- 27 个测试全部通过

### 后端测试

```bash
cd core && go test ./... -timeout 60s
```

当前状态：

- 默认测试可通过
- `exchange_rate/provider` 中依赖外网的集成测试默认跳过
- 如需执行这类测试，可设置 `ALLFI_RUN_ONLINE_TESTS=1`

---

## 7. 开发注意事项

- Swagger 地址是 `/swagger/`，不是 `/api/v1/docs`
- OpenAPI JSON 地址是 `/api.json`
- 当前后端主数据访问层是 GoFrame `gdb`，不是 GORM
- 默认数据库连接是 `sqlite::@file(./data/allfi.db)`
- 发布模式前端通过 `embed.FS` 内嵌到后端

---

## 8. 建议的阅读顺序

1. `docs/project-status.md`
2. `docs/tech/tech-baseline.md`
3. `docs/specs/frontend-spec.md` 或 `docs/specs/backend-spec.md`
4. 需要部署时再看 `docs/guides/deployment-guide.md`
