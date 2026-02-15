# AllFi 部署指南

> **版本**：v2.0 | **更新时间**：2026-02-11

本文档详细说明 AllFi 的部署方式，包括 Docker 一键部署和手动部署两种方案。

---

## 目录

- [环境要求](#环境要求)
- [Docker 部署（推荐）](#docker-部署推荐)
- [手动部署](#手动部署)
- [HTTPS/TLS 配置](#httpstls-配置)
- [数据备份与恢复](#数据备份与恢复)
- [系统服务配置（Linux）](#系统服务配置linux)
- [常见问题](#常见问题)

---

## 环境要求

### Docker 部署

| 依赖 | 最低版本 | 说明 |
|------|---------|------|
| Docker | 20.10+ | 容器运行时 |
| Docker Compose | 2.0+ | 容器编排 |
| 磁盘空间 | 1 GB+ | 镜像 + 数据存储 |
| 内存 | 512 MB+ | 推荐 1 GB |

### 手动部署

| 依赖 | 版本要求 | 说明 |
|------|---------|------|
| Go | **1.24** | 后端编译（项目固定版本） |
| Node.js | 18+ | 前端构建 |
| pnpm | 9+ | 前端包管理器 |
| Nginx | 1.20+ | 前端静态文件托管（可选） |

> **注意**：后端使用 `modernc.org/sqlite`（纯 Go SQLite 驱动），无需 CGO 和 GCC。

### 外部 API 密钥（可选但推荐）

| 服务 | 用途 | 免费申请地址 |
|------|------|-------------|
| Etherscan | 以太坊链上数据 | https://etherscan.io/apis |
| BscScan | BSC 链上数据 | https://bscscan.com/apis |
| CoinGecko | 加密货币价格 | https://www.coingecko.com/en/api |

> 后端 Etherscan 统一客户端还支持 Polygon/Arbitrum/Optimism/Base，对应 API Key 可复用 Etherscan Key（均属 Etherscan 系列服务）。

---

## Docker 部署（推荐）

Docker Compose 一键部署是最简单的方式，适合大多数用户。

### 步骤一：克隆仓库

```bash
git clone https://github.com/your-finance/allfi.git
cd allfi
```

### 步骤二：配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env
```

编辑 `.env` 文件，至少需要设置加密主密钥：

```bash
# 生成随机密钥（必须修改！）
openssl rand -base64 32
```

将生成的密钥填入 `.env`：

```env
# 加密主密钥（必须修改！用于 AES-256-GCM 加密 API Key）
ALLFI_MASTER_KEY=你生成的随机密钥

# Etherscan API Key（可选，链上功能需要）
ETHERSCAN_API_KEY=你的密钥

# BscScan API Key（可选，BSC 链功能需要）
BSCSCAN_API_KEY=你的密钥

# CoinGecko API Key（可选）
COINGECKO_API_KEY=

# 服务端口（默认即可）
SERVER_PORT=8080
FRONTEND_PORT=3174
```

### 步骤三：启动服务

```bash
# 构建并启动（后台运行）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 步骤四：验证部署

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 预期响应：
# {"code":0,"message":"AllFi is running","data":{"status":"ok","version":"0.2.0"}}
```

访问前端页面：http://localhost:3174

首次访问会进入 PIN 码设置流程（4-8 位数字），设置后自动登录。

### 常用 Docker 命令

```bash
# 停止服务
docker-compose down

# 重新构建（代码更新后）
docker-compose up -d --build

# 查看后端日志
docker logs -f allfi-backend

# 查看前端日志
docker logs -f allfi-frontend

# 进入后端容器排查问题
docker exec -it allfi-backend sh
```

### Docker 架构说明

```
                    用户浏览器
                        |
                   :3174 (HTTP)
                        |
              +---------v---------+
              |  allfi-frontend   |
              |  (Nginx + Vue)    |
              +------|-----|------+
                     |     |
              静态文件  API 代理
                     |     |
                     |  :8080
                     |     |
              +------v-----------+
              |  allfi-backend   |
              |  (Go API Server) |
              +------|-----|-----+
                     |     |
              +------v-----v-----+
              |  allfi-data      |
              |  (Docker Volume) |
              |  SQLite3 数据库    |
              +------------------+
```

- **allfi-frontend**：Nginx 托管 Vue 构建产物，同时反向代理 API 请求到后端
- **allfi-backend**：Go API 服务器（net/http），处理业务逻辑和数据存储
- **allfi-data**：Docker 卷，持久化 SQLite 数据库文件

### 安全特性

Docker 部署默认启用了以下安全措施：

- `no-new-privileges: true` -- 禁止容器内进程提权
- `read_only: true` -- 根文件系统只读
- tmpfs 挂载 `/tmp` 和日志目录
- 健康检查自动重启

---

## 手动部署

适合需要更灵活控制的用户，或没有 Docker 环境的场景。

### 步骤一：编译后端

```bash
cd core

# 安装 Go 依赖
go mod download

# 编译生产二进制文件（纯 Go，无需 CGO）
CGO_ENABLED=0 go build -ldflags="-s -w" -o allfi cmd/server/main.go
```

编译产物为 `core/allfi` 二进制文件，约 20-30 MB。

> **提示**：项目使用 `modernc.org/sqlite`（纯 Go SQLite 驱动），编译时无需 CGO，也不需要安装 GCC。

### 步骤二：构建前端

```bash
cd webapp

# 安装依赖
pnpm install

# 构建生产版本
pnpm build
```

构建产物位于 `webapp/dist/` 目录。

### 步骤三：配置后端

后端使用 GoFrame gcfg 管理配置，配置文件位于 `core/manifest/config/config.yaml`。

关键配置项：

```yaml
# 服务器配置
server:
  address: "0.0.0.0:8080"

# 应用信息
app:
  name: "AllFi"
  version: "0.2.0"
  env: "production"

# 数据库配置
database:
  default:
    link: "sqlite::@file(./data/allfi.db)"
    debug: false

# 安全配置
security:
  masterKey: "你的加密主密钥"  # 或设置环境变量 ALLFI_MASTER_KEY

# 外部 API
externalAPIs:
  etherscan:
    apiKey: ""
  bscscan:
    apiKey: ""
  coingecko:
    apiKey: ""

# 定时任务
cron:
  snapshotInterval: 60    # 快照间隔（分钟）
  priceCacheTTL: 5        # 价格缓存（分钟）
  rateCacheTTL: 60        # 汇率缓存（分钟）
```

所有配置项均支持环境变量覆盖：

| 环境变量 | 配置路径 | 默认值 |
|---------|---------|--------|
| `ALLFI_PORT` | `server.port` | 8080 |
| `ALLFI_MODE` | `server.mode` | development |
| `ALLFI_DB_TYPE` | `database.type` | sqlite |
| `ALLFI_DB_PATH` | `database.sqlite.path` | ../data/allfi.db |
| `ALLFI_MASTER_KEY` | `security.master_key` | - |
| `ETHERSCAN_API_KEY` | `external_apis.etherscan.api_key` | - |
| `BSCSCAN_API_KEY` | `external_apis.bscscan.api_key` | - |
| `COINGECKO_API_KEY` | `external_apis.coingecko.api_key` | - |

### 步骤四：配置 Nginx

创建 Nginx 配置文件 `/etc/nginx/conf.d/allfi.conf`：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 或 localhost

    # 前端静态文件
    root /path/to/allfi/webapp/dist;
    index index.html;

    # 安全 Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml;
    gzip_min_length 256;

    # API 反向代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # SPA 路由回退
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 静态资源长期缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 禁止访问隐藏文件
    location ~ /\. {
        deny all;
    }
}
```

重新加载 Nginx：

```bash
sudo nginx -t          # 测试配置
sudo systemctl reload nginx
```

### 步骤五：启动后端

```bash
cd core
./allfi
```

或使用 nohup 后台运行：

```bash
nohup ./allfi > /var/log/allfi.log 2>&1 &
```

### 步骤六：验证部署

```bash
# 后端健康检查
curl http://localhost:8080/api/v1/health

# 前端页面（通过 Nginx）
curl http://localhost/
```

---

## HTTPS/TLS 配置

生产环境强烈建议启用 HTTPS。以下提供两种方式。

### 方式一：自签名证书（开发/内网）

```bash
# 生成自签名证书
mkdir -p certs
openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -subj "/CN=allfi.local"
```

### 方式二：Let's Encrypt（公网部署）

使用 certbot 自动申请和续期证书：

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 申请证书（需要域名已解析到服务器）
sudo certbot --nginx -d your-domain.com

# 证书自动续期（certbot 会自动设置 cron）
sudo certbot renew --dry-run
```

### Nginx HTTPS 配置

将以下配置添加到 Nginx 站点配置：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # TLS 安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # 其余配置同 HTTP 方案
    root /path/to/allfi/webapp/dist;
    index index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}

# HTTP 自动跳转 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}
```

### Docker 环境启用 HTTPS

1. 将证书文件放入 `certs/` 目录
2. 修改 `docker-compose.yml`，开放 443 端口：

```yaml
frontend:
  ports:
    - "80:80"
    - "443:443"
  volumes:
    - ./certs:/etc/nginx/certs:ro
```

3. 重新构建并启动：

```bash
docker-compose up -d --build
```

---

## 数据备份与恢复

### 数据存储位置

| 部署方式 | 数据路径 | 说明 |
|---------|---------|------|
| Docker | Docker 卷 `allfi-data` | 挂载到容器 `/app/data` |
| 手动部署 | `core/data/allfi.db` | SQLite 数据库文件 |

AllFi 使用 SQLite 单文件数据库，备份非常简单。

### 备份

#### Docker 部署备份

```bash
# 方式一：直接从容器复制
docker cp allfi-backend:/app/data/allfi.db ./backup/allfi-$(date +%Y%m%d).db

# 方式二：停止服务后备份卷目录
docker-compose down
sudo cp /var/lib/docker/volumes/allfi_allfi-data/_data/allfi.db ./backup/
docker-compose up -d
```

#### 手动部署备份

```bash
# 直接复制数据库文件（推荐先停止服务）
cp core/data/allfi.db backup/allfi-$(date +%Y%m%d).db
```

#### 自动备份脚本

创建 `scripts/backup.sh`：

```bash
#!/bin/bash
# AllFi 自动备份脚本
BACKUP_DIR="/path/to/backup"
DATE=$(date +%Y%m%d_%H%M%S)
DB_FILE="core/data/allfi.db"

mkdir -p "$BACKUP_DIR"

# 使用 SQLite 在线备份（无需停止服务）
sqlite3 "$DB_FILE" ".backup '$BACKUP_DIR/allfi-$DATE.db'"

# 保留最近 30 天的备份
find "$BACKUP_DIR" -name "allfi-*.db" -mtime +30 -delete

echo "备份完成: $BACKUP_DIR/allfi-$DATE.db"
```

设置定时任务（每天凌晨 3 点自动备份）：

```bash
chmod +x scripts/backup.sh
crontab -e
# 添加：
# 0 3 * * * /path/to/allfi/scripts/backup.sh
```

### 恢复

```bash
# 停止服务
docker-compose down  # 或 kill 后端进程

# 替换数据库文件
cp backup/allfi-20260210.db core/data/allfi.db

# 重新启动
docker-compose up -d  # 或 ./allfi
```

> **注意**：恢复前请确认备份文件的完整性。可以用 `sqlite3 backup.db "PRAGMA integrity_check;"` 验证。

---

## 系统服务配置（Linux）

在 Linux 服务器上，推荐使用 systemd 管理 AllFi 后端服务。

### 创建 systemd 服务文件

创建 `/etc/systemd/system/allfi.service`：

```ini
[Unit]
Description=AllFi Backend API Server
After=network.target

[Service]
Type=simple
User=allfi
Group=allfi
WorkingDirectory=/opt/allfi/core
ExecStart=/opt/allfi/core/allfi
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

# 环境变量
Environment=ALLFI_MASTER_KEY=你的密钥
Environment=ALLFI_MODE=production

# 安全加固
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/allfi/core/data /opt/allfi/core/logs

[Install]
WantedBy=multi-user.target
```

### 启用服务

```bash
# 创建专用用户
sudo useradd -r -s /bin/false allfi

# 设置目录权限
sudo chown -R allfi:allfi /opt/allfi

# 启用并启动
sudo systemctl daemon-reload
sudo systemctl enable allfi
sudo systemctl start allfi

# 查看状态
sudo systemctl status allfi

# 查看日志
sudo journalctl -u allfi -f
```

---

## 代理配置（国内用户）

如果你在国内或使用需要代理的网络环境，可能需要配置以下代理：

- **Go 模块代理**：`GOPROXY=https://goproxy.cn,direct`
- **pnpm 镜像源**：`pnpm config set registry https://registry.npmmirror.com`
- **Docker 镜像加速**：配置 Docker daemon 的 registry-mirrors
- **外部 API 代理**：在 `.env` 中设置 `HTTP_PROXY` / `HTTPS_PROXY`

> 详细配置方案见 [代理配置指南](./proxy-guide.md)。

---

## 常见问题

### 1. Docker 构建失败

**说明**：项目使用 `modernc.org/sqlite`（纯 Go SQLite 驱动），已不再需要 CGO。构建时设置 `CGO_ENABLED=0`，无需安装 GCC 或其他 C 编译工具。

如果遇到构建问题，请确认 Dockerfile 中 `CGO_ENABLED=0` 已正确设置。

### 2. 前端无法连接后端 API

**现象**：页面加载成功但数据为空，浏览器控制台出现 CORS 或网络错误。

**解决**：
- Docker 部署：确认 `allfi-backend` 容器健康（`docker ps` 检查状态）
- 手动部署：确认 Nginx 反向代理配置正确（`/api/` 代理到 `127.0.0.1:8080`）
- 检查后端是否正常运行：`curl http://localhost:8080/api/v1/health`

### 3. 端口被占用

**现象**：启动时提示端口 `8080` 或 `3174` 已被占用。

**解决**：修改 `.env` 中的端口配置：

```env
SERVER_PORT=9090
FRONTEND_PORT=3000
```

或者查找并关闭占用端口的进程：

```bash
lsof -i :8080
```

### 4. 数据库文件权限问题

**现象**：后端启动时提示无法读写数据库文件。

**解决**：

```bash
# 确保数据目录存在
mkdir -p core/data

# 设置正确权限
chmod 700 core/data
chmod 600 core/data/allfi.db  # 如果文件已存在
```

### 5. 链上数据获取失败

**现象**：钱包资产显示为空或报错。

**原因**：缺少 Etherscan/BscScan API 密钥。

**解决**：
1. 前往 https://etherscan.io/apis 免费注册获取 API Key
2. 在 `.env` 或 `config.yaml` 中配置对应密钥
3. 重启服务

### 6. Docker 卷数据丢失

**说明**：`docker-compose down` 不会删除命名卷。数据只会在以下情况丢失：

```bash
# 会删除数据卷（危险！）
docker-compose down -v

# 不会删除数据卷（安全）
docker-compose down
```

### 7. 如何更新到新版本？

```bash
# 拉取最新代码
git pull origin master

# Docker 部署：重新构建
docker-compose up -d --build

# 手动部署：重新编译
cd core && CGO_ENABLED=0 go build -ldflags="-s -w" -o allfi cmd/server/main.go
cd ../webapp && pnpm build
# 重启服务
```

数据库迁移会在启动时自动执行（GORM AutoMigrate）。

---

## 相关文档

- [技术基线](../tech/tech-baseline.md)
- [API 接口文档](../tech/api-reference.md)
- [开发指南](./dev-guide.md)
- [代理配置指南](./proxy-guide.md)
- [用户指南](./user-guide.md)

---

**文档维护者**: @allfi
**最后更新**: 2026-02-13
