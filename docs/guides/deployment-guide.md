# AllFi 部署指南

> 审计日期：2026-03-31
> 重点：本指南只保留当前真实可执行的部署路径。

---

## 1. 先选部署方式

| 场景 | 推荐方式 | 说明 |
|------|---------|------|
| 最终用户，不想克隆源码 | `deploy/docker-deploy.sh` | 下载 Release 二进制并生成独立部署目录 |
| 维护者，需要跑仓库内 Docker Compose | 根目录 `docker-compose.yml` | 需要先手动构建本地镜像 |
| 本地开发 | `make dev` | 前后端联调 |
| 宿主机直接运行二进制 | GitHub Release | 单二进制运行，默认监听 8080 |

---

## 2. 最终用户部署：一键脚本

### 前置条件

- Docker 20.10+
- Docker Compose v2 或 `docker-compose`
- `curl`
- `tar`

### 部署命令

```bash
curl -sSL https://raw.githubusercontent.com/your-finance/allfi/master/deploy/docker-deploy.sh | bash
```

### 脚本会做什么

1. 检测系统架构
2. 拉取最新 GitHub Release
3. 生成独立部署目录（默认 `allfi-docker/`）
4. 生成 `.env`
5. 生成最小化 `Dockerfile` 和 `docker-compose.yml`
6. 构建并启动服务

### 默认访问地址

- 页面与 API：`http://localhost:3000`
- Swagger：`http://localhost:3000/swagger/`

### 当前限制

- 该最小化部署当前**不包含** `updater` sidecar
- 因此不应依赖页面内“一键更新”
- 升级建议：重新执行部署脚本或手动拉取新 Release 重新部署

---

## 3. 维护者部署：仓库根目录 Docker Compose

### 重要说明

根目录 `docker-compose.yml` 中的 `backend` 服务默认写的是：

```yaml
image: allfi-backend:latest
```

这意味着直接执行 `docker compose up -d` 并不会从源码自动构建后端镜像。

### 推荐步骤

```bash
git clone https://github.com/your-finance/allfi.git
cd allfi
cp .env.example .env
```

生成主密钥并写入 `.env`：

```bash
openssl rand -base64 32
```

先构建本地镜像：

```bash
docker build -f core/Dockerfile -t allfi-backend:latest .
```

再启动 Compose：

```bash
docker compose up -d --build
```

### 默认地址

- 页面与 API：`http://localhost:3000`
- 容器内服务端口：`8080`

### 这个模式的特点

- 包含 `updater` sidecar
- 更适合维护者、演示环境或需要验证 Docker 更新流程的场景

---

## 4. 本地开发模式

### 依赖

- Go 1.24
- Node.js 20+
- pnpm

### 命令

```bash
make setup
make dev
```

### 地址

- 前端：`http://localhost:3000`
- 后端 API：`http://localhost:8080/api/v1`
- Swagger：`http://localhost:8080/swagger/`

---

## 5. 宿主机二进制运行

### 适用场景

- 不想装 Go / Node.js
- 不想使用 Docker
- 需要验证宿主机 OTA 更新逻辑

### 步骤

1. 从 GitHub Release 下载对应平台压缩包
2. 解压得到 `allfi` 二进制
3. 准备配置文件
4. 运行：

```bash
./allfi
```

### 默认地址

- 页面与 API：`http://localhost:8080`
- Swagger：`http://localhost:8080/swagger/`

---

## 6. 端口说明

| 模式 | 默认外部端口 | 备注 |
|------|-------------|------|
| 本地开发 | 3000 | Vite dev server |
| 仓库根目录 Docker Compose | 3000 | 由根目录 `.env.example` 控制 |
| 一键脚本生成部署 | 3000 | 由脚本生成的 `.env` 控制 |
| 宿主机二进制 | 8080 | 由配置文件 `server.address` 控制 |

---

## 7. 备份与恢复

### Docker 卷备份

对于 Docker 部署，核心数据保存在容器的 `/app/data`。

可以直接导出命名卷：

```bash
docker run --rm -v allfi-data:/data -v "$(pwd)":/backup alpine:3.21 \
  tar czf /backup/allfi-data.tar.gz -C /data .
```

恢复时反向解压：

```bash
docker run --rm -v allfi-data:/data -v "$(pwd)":/backup alpine:3.21 \
  tar xzf /backup/allfi-data.tar.gz -C /data
```

### 宿主机二进制

- 直接备份运行目录下的 `data/`
- 如有自定义配置，也一起备份 `manifest/config/`

---

## 8. HTTPS 与反向代理

AllFi 生产环境通常建议放在反向代理后面。

### Caddy 示例

```caddy
allfi.example.com {
    reverse_proxy 127.0.0.1:3000
}
```

### Nginx 示例

```nginx
server {
    listen 443 ssl http2;
    server_name allfi.example.com;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 9. 常用维护命令

```bash
docker compose logs -f
docker compose ps
docker compose restart
docker compose down
```

对于一键脚本部署，请先进入生成的部署目录再执行这些命令。
