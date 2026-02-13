# 代理与镜像源配置指南

> **版本**：v1.0 | **更新时间**：2026-02-13

本文档为国内用户或需要通过代理访问外部服务的用户提供完整的配置方案。

---

## 目录

- [概览](#概览)
- [Go 模块代理](#go-模块代理)
- [pnpm 镜像源](#pnpm-镜像源)
- [Docker 镜像加速](#docker-镜像加速)
- [HTTP 代理（外部 API）](#http-代理外部-api)
- [Docker Compose 代理配置](#docker-compose-代理配置)
- [一键配置方案](#一键配置方案)
- [常见问题](#常见问题)

---

## 概览

AllFi 开发和运行中涉及以下外部网络访问：

| 场景 | 访问目标 | 配置方式 |
|------|---------|---------|
| Go 依赖下载 | proxy.golang.org | GOPROXY 环境变量 |
| 前端依赖下载 | registry.npmjs.org | pnpm registry 配置 |
| Docker 镜像拉取 | docker.io / ghcr.io | Docker daemon 镜像加速 |
| 后端调用外部 API | etherscan.io / coingecko.com | HTTP_PROXY 环境变量 |
| Docker 构建时下载依赖 | 同上 | docker-compose build args |

---

## Go 模块代理

### 推荐方案：使用国内 Go 代理

```bash
# 设置 Go 模块代理（推荐 goproxy.cn）
export GOPROXY=https://goproxy.cn,direct

# 可选：跳过校验和检查（遇到私有模块时需要）
# export GONOSUMCHECK=*
# export GONOSUMDB=*
```

### 持久化配置

将以上配置写入 shell 配置文件（`.bashrc` / `.zshrc`）：

```bash
# ~/.zshrc 或 ~/.bashrc
export GOPROXY=https://goproxy.cn,direct
```

### 其他可选代理

| 代理 | 地址 | 说明 |
|------|------|------|
| goproxy.cn | `https://goproxy.cn` | 七牛云维护，推荐 |
| goproxy.io | `https://goproxy.io` | 社区维护 |
| mirrors.aliyun.com | `https://mirrors.aliyun.com/goproxy/` | 阿里云镜像 |

### 验证配置

```bash
go env GOPROXY
# 预期输出: https://goproxy.cn,direct

# 测试下载
cd core && go mod download
```

---

## pnpm 镜像源

### 推荐方案：使用 npmmirror

```bash
# 全局设置 pnpm 镜像源
pnpm config set registry https://registry.npmmirror.com
```

### 项目级配置（推荐）

在项目根目录创建 `.npmrc` 文件，对所有开发者生效：

```ini
# .npmrc
registry=https://registry.npmmirror.com
```

> 注意：本项目暂不提交 `.npmrc` 到仓库，避免影响海外用户。如需要，可自行创建。

### 验证配置

```bash
pnpm config get registry
# 预期输出: https://registry.npmmirror.com

# 测试安装
cd webapp && pnpm install
```

### 其他可选镜像源

| 镜像 | 地址 |
|------|------|
| npmmirror（原淘宝源）| `https://registry.npmmirror.com` |
| 腾讯云 | `https://mirrors.cloud.tencent.com/npm/` |
| 华为云 | `https://repo.huaweicloud.com/repository/npm/` |

---

## Docker 镜像加速

### Linux 配置

编辑 Docker daemon 配置文件：

```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
EOF

# 重启 Docker
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### macOS（Docker Desktop）配置

1. 打开 Docker Desktop → Settings → Docker Engine
2. 在 JSON 配置中添加：

```json
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
```

3. 点击 Apply & Restart

### 验证配置

```bash
docker info | grep -A 5 "Registry Mirrors"
```

### Docker 构建时使用代理

如果 Docker 构建过程中需要下载依赖（如 Go 模块、npm 包），可通过 build args 传递代理：

```bash
# 带代理构建
docker-compose build \
  --build-arg HTTP_PROXY=http://host.docker.internal:7890 \
  --build-arg HTTPS_PROXY=http://host.docker.internal:7890
```

> 说明：`host.docker.internal` 是 Docker Desktop 中宿主机的地址。Linux 可能需要使用宿主机实际 IP。

---

## HTTP 代理（外部 API）

AllFi 后端需要访问以下外部 API：

| 服务 | 地址 | 用途 |
|------|------|------|
| Etherscan | api.etherscan.io | 以太坊链上数据 |
| BscScan | api.bscscan.com | BSC 链上数据 |
| CoinGecko | api.coingecko.com | 加密货币价格 |
| Yahoo Finance | query1.finance.yahoo.com | 汇率数据 |

### Go 原生代理支持

Go 标准库 `net/http` **自动尊重** `HTTP_PROXY`、`HTTPS_PROXY`、`NO_PROXY` 环境变量，AllFi 无需额外代码修改。

### 本地开发配置

```bash
# 设置代理（替换为你的代理地址和端口）
export HTTP_PROXY=http://127.0.0.1:7890
export HTTPS_PROXY=http://127.0.0.1:7890
export NO_PROXY=localhost,127.0.0.1

# 启动后端
cd core && go run cmd/server/main.go
```

### 使用 .env 配置

编辑项目根目录 `.env`：

```env
# 代理配置（取消注释并修改为你的代理地址）
HTTP_PROXY=http://127.0.0.1:7890
HTTPS_PROXY=http://127.0.0.1:7890
NO_PROXY=localhost,127.0.0.1,backend
```

> 注意：`NO_PROXY` 中包含 `backend` 是为了确保 Docker 容器间通信不走代理。

---

## Docker Compose 代理配置

`docker-compose.yml` 已支持通过 `.env` 文件传递代理环境变量到容器：

```yaml
# docker-compose.yml（已配置）
environment:
  - HTTP_PROXY=${HTTP_PROXY:-}
  - HTTPS_PROXY=${HTTPS_PROXY:-}
  - NO_PROXY=${NO_PROXY:-localhost,127.0.0.1}
```

使用方法：

1. 编辑 `.env`，取消代理变量的注释并填入你的代理地址
2. 启动服务：`docker-compose up -d`
3. 代理变量会自动传递到后端容器

### 构建时代理

Docker 构建阶段（下载 Go 模块 / npm 包）也可能需要代理：

```bash
# 方法一：通过环境变量（推荐）
HTTP_PROXY=http://host.docker.internal:7890 \
HTTPS_PROXY=http://host.docker.internal:7890 \
docker-compose build

# 方法二：修改 Dockerfile 添加 ARG（不推荐提交到仓库）
```

---

## 一键配置方案

### 方案 A：环境变量汇总（开发环境）

在 `~/.zshrc` 或 `~/.bashrc` 中一次性配置：

```bash
# === AllFi 开发环境代理配置 ===

# Go 模块代理
export GOPROXY=https://goproxy.cn,direct

# HTTP 代理（如果需要）
# export HTTP_PROXY=http://127.0.0.1:7890
# export HTTPS_PROXY=http://127.0.0.1:7890
# export NO_PROXY=localhost,127.0.0.1
```

pnpm 镜像源单独配置（只需运行一次）：

```bash
pnpm config set registry https://registry.npmmirror.com
```

### 方案 B：使用快速启动脚本检测

```bash
# 检测当前环境和代理配置
bash scripts/quickstart.sh --check
```

脚本会自动检测 GOPROXY 和 pnpm registry 的配置状态，并给出建议。

---

## 常见问题

### 1. Go 模块下载超时

**现象**：`go mod download` 长时间无响应或报错 `dial tcp: i/o timeout`

**解决**：

```bash
export GOPROXY=https://goproxy.cn,direct
go mod download
```

### 2. pnpm install 下载慢

**现象**：`pnpm install` 进度缓慢或出现网络错误

**解决**：

```bash
pnpm config set registry https://registry.npmmirror.com
pnpm install
```

### 3. Docker 镜像拉取失败

**现象**：`docker pull` 或 `docker-compose build` 报错 `TLS handshake timeout`

**解决**：配置 Docker 镜像加速器（见上方 Docker 章节）。

### 4. 后端无法获取链上数据

**现象**：后端启动正常，但钱包资产显示为空

**原因**：
1. 缺少 API Key（Etherscan/BscScan）
2. 网络无法访问外部 API

**解决**：
1. 配置 API Key（见 `.env.example`）
2. 设置 HTTP_PROXY 环境变量

### 5. Docker 容器内无法访问宿主机代理

**现象**：容器中的 HTTP_PROXY 指向 `127.0.0.1`，但无法连接

**解决**：
- Docker Desktop（macOS/Windows）：使用 `host.docker.internal` 替代 `127.0.0.1`
- Linux：使用宿主机实际 IP（如 `172.17.0.1`）

```env
# .env（Docker Desktop）
HTTP_PROXY=http://host.docker.internal:7890
HTTPS_PROXY=http://host.docker.internal:7890
NO_PROXY=localhost,127.0.0.1,backend
```

### 6. 代理环境下 Docker 构建失败

**现象**：Dockerfile 中的 `go mod download` 或 `pnpm install` 失败

**解决**：在 `docker-compose build` 时传递代理参数：

```bash
docker-compose build \
  --build-arg HTTP_PROXY=http://host.docker.internal:7890 \
  --build-arg HTTPS_PROXY=http://host.docker.internal:7890
```

---

## 相关文档

- [部署指南](./deployment-guide.md)
- [开发指南](./dev-guide.md)
- [技术基线](../tech/tech-baseline.md)

---

**文档维护者**: @allfi
**最后更新**: 2026-02-13
