#!/bin/bash
# AllFi 零依赖一键 Docker 部署脚本（完全免源码）
# 用法:
#   curl -sSL https://raw.githubusercontent.com/your-finance/allfi/master/deploy/docker-deploy.sh | bash
#
# 前置条件: Docker 20.10+, Docker Compose v2+, curl, tar
# 本脚本会自动:
#   1. 检测系统架构（amd64/arm64）
#   2. 从 GitHub Releases 下载对应平台的预编译二进制包
#   3. 生成 .env + 安全密钥
#   4. 生成 Dockerfile 和 docker-compose.yml
#   5. 构建轻量级 Docker 镜像并启动服务
# 全程无需 git / go / node / pnpm 等开发工具

set -e

# 颜色定义
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${CYAN}║       AllFi — 全资产聚合平台 一键部署        ║${RESET}"
echo -e "${CYAN}╚══════════════════════════════════════════════╝${RESET}"
echo ""

# ==================== 前置检查 ====================

echo -e "${CYAN}>>> 检查前置条件...${RESET}"

# 检查 Docker
if ! command -v docker >/dev/null 2>&1; then
    echo -e "${RED}  ✗ 未检测到 Docker，请先安装 Docker${RESET}"
    echo "    安装指南: https://docs.docker.com/engine/install/"
    exit 1
fi
echo -e "${GREEN}  ✓ Docker $(docker --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo 'detected')${RESET}"

# 检查 Docker Compose (v2 插件优先，回退 docker-compose)
COMPOSE_CMD=""
if docker compose version >/dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
    echo -e "${GREEN}  ✓ Docker Compose (plugin) detected${RESET}"
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_CMD="docker-compose"
    echo -e "${GREEN}  ✓ docker-compose (standalone) detected${RESET}"
else
    echo -e "${RED}  ✗ 未检测到 Docker Compose，请先安装${RESET}"
    echo "    安装指南: https://docs.docker.com/compose/install/"
    exit 1
fi

# 检查 curl
if ! command -v curl >/dev/null 2>&1; then
    echo -e "${RED}  ✗ 未检测到 curl，请先安装${RESET}"
    exit 1
fi
echo -e "${GREEN}  ✓ curl detected${RESET}"

# 检查 tar
if ! command -v tar >/dev/null 2>&1; then
    echo -e "${RED}  ✗ 未检测到 tar，请先安装${RESET}"
    exit 1
fi
echo -e "${GREEN}  ✓ tar detected${RESET}"

# 检测系统架构
ARCH=$(uname -m)
case $ARCH in
    x86_64)         PLATFORM="linux-amd64" ;;
    aarch64|arm64)  PLATFORM="linux-arm64" ;;
    *)
        echo -e "${RED}  ✗ 不支持的系统架构: $ARCH${RESET}"
        exit 1
        ;;
esac
echo -e "${GREEN}  ✓ 检测到架构: $PLATFORM${RESET}"

echo ""

# ==================== 初始化工作目录 ====================

DEPLOY_DIR="${ALLFI_DEPLOY_DIR:-allfi-docker}"

if [ ! -d "$DEPLOY_DIR" ]; then
    mkdir -p "$DEPLOY_DIR"
    echo -e "${GREEN}  ✓ 已创建部署目录: $DEPLOY_DIR${RESET}"
fi
cd "$DEPLOY_DIR"

# ==================== 获取最新版本 ====================

echo -e "${CYAN}>>> 获取最新版本信息...${RESET}"

LATEST_JSON=$(curl -sSL https://api.github.com/repos/your-finance/allfi/releases/latest 2>/dev/null || echo "")
if [ -z "$LATEST_JSON" ]; then
    echo -e "${RED}  ✗ 无法连接 GitHub API，请检查网络${RESET}"
    echo "    如果使用代理，请设置 HTTP_PROXY / HTTPS_PROXY 环境变量"
    exit 1
fi

VERSION=$(echo "$LATEST_JSON" | grep -o '"tag_name": "[^"]*"' | head -1 | cut -d'"' -f4 | sed 's/^v//')
if [ -z "$VERSION" ]; then
    echo -e "${RED}  ✗ 获取最新版本失败，请检查网络（如果受限可配置 HTTP_PROXY）${RESET}"
    exit 1
fi
echo -e "${GREEN}  ✓ 最新版本: v$VERSION${RESET}"

echo ""

# ==================== 下载二进制发布包 ====================

TARBALL="allfi-${VERSION}-${PLATFORM}.tar.gz"
DOWNLOAD_URL="https://github.com/your-finance/allfi/releases/download/v${VERSION}/${TARBALL}"

echo -e "${CYAN}>>> 下载二进制发布包: ${TARBALL}...${RESET}"
if ! curl -# -fSL -o "$TARBALL" "$DOWNLOAD_URL"; then
    echo -e "${RED}  ✗ 下载失败，请检查网络或在浏览器中访问:${RESET}"
    echo "    $DOWNLOAD_URL"
    exit 1
fi
echo -e "${GREEN}  ✓ 下载完成${RESET}"

echo ""

# ==================== 解压发布包 ====================

echo -e "${CYAN}>>> 解压发布包...${RESET}"

tar -xzf "$TARBALL"

# release 的目录结构为 allfi-{version}-{platform}/allfi + config.yaml
DIR_NAME="allfi-${VERSION}-${PLATFORM}"
if [ -d "$DIR_NAME" ]; then
    # 移动二进制文件和配置文件到当前目录
    cp "$DIR_NAME/allfi" ./allfi
    chmod +x ./allfi
    mkdir -p manifest/config
    if [ -f "$DIR_NAME/config.yaml" ]; then
        cp "$DIR_NAME/config.yaml" manifest/config/config.yaml
    fi
    # 清理解压目录和压缩包
    rm -rf "$DIR_NAME" "$TARBALL"
else
    # 兼容其他目录结构
    echo -e "${YELLOW}  ⚠ 目录结构与预期不同，请检查解压内容${RESET}"
    ls -la
    rm -f "$TARBALL"
fi

echo -e "${GREEN}  ✓ 解压完成${RESET}"

echo ""

# ==================== 生成 .env ====================

if [ ! -f ".env" ]; then
    echo -e "${CYAN}>>> 生成配置文件 .env...${RESET}"

    # 生成随机安全密钥
    MASTER_KEY=$(head -c 32 /dev/urandom | base64 | tr -d '\n' | head -c 44)

    cat > .env << ENV_EOF
# AllFi 配置文件（由部署脚本自动生成）
# 端口配置
ALLFI_PORT=3174

# 加密密钥（自动生成，请勿外泄）
ALLFI_MASTER_KEY=${MASTER_KEY}

# 时区配置
TZ=Asia/Shanghai

# 第三方 API Key（可选）
# ETHERSCAN_API_KEY=
# BSCSCAN_API_KEY=
# COINGECKO_API_KEY=

# 代理配置（可选，国内用户按需启用）
# HTTP_PROXY=http://127.0.0.1:7890
# HTTPS_PROXY=http://127.0.0.1:7890
ENV_EOF

    echo -e "${GREEN}  ✓ .env 已生成（ALLFI_MASTER_KEY 已自动配置）${RESET}"
else
    echo -e "${YELLOW}  ⊘ .env 已存在，跳过（如需重新生成请先删除 .env）${RESET}"
fi

echo ""

# ==================== 生成 Dockerfile ====================

echo -e "${CYAN}>>> 生成 Dockerfile...${RESET}"

cat > Dockerfile << 'DOCKER_EOF'
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata wget
RUN addgroup -S allfi && adduser -S allfi -G allfi
WORKDIR /app
COPY allfi /app/allfi
COPY manifest/config/config.yaml manifest/config/config.yaml
RUN mkdir -p /app/data /app/logs && chown -R allfi:allfi /app
USER allfi
EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["/app/allfi"]
DOCKER_EOF

echo -e "${GREEN}  ✓ Dockerfile 已生成${RESET}"

echo ""

# ==================== 生成 docker-compose.yml ====================

echo -e "${CYAN}>>> 生成 docker-compose.yml...${RESET}"

cat > docker-compose.yml << 'COMPOSE_EOF'
# AllFi Docker Compose 编排（由部署脚本自动生成）
services:
  backend:
    build: .
    image: allfi-backend:latest
    container_name: allfi-backend
    restart: unless-stopped
    ports:
      - "${ALLFI_PORT:-3174}:8080"
    volumes:
      - allfi-data:/app/data
      - /etc/localtime:/etc/localtime:ro
    environment:
      - TZ=${TZ:-Asia/Shanghai}
      - ALLFI_MASTER_KEY=${ALLFI_MASTER_KEY}
      - ETHERSCAN_API_KEY=${ETHERSCAN_API_KEY:-}
      - BSCSCAN_API_KEY=${BSCSCAN_API_KEY:-}
      - COINGECKO_API_KEY=${COINGECKO_API_KEY:-}
      - HTTP_PROXY=${HTTP_PROXY:-}
      - HTTPS_PROXY=${HTTPS_PROXY:-}
      - NO_PROXY=${NO_PROXY:-localhost,127.0.0.1}
    security_opt:
      - no-new-privileges:true
    mem_limit: 512m
    memswap_limit: 512m
    healthcheck:
      test: [ "CMD", "wget", "-qO", "/dev/null", "http://localhost:8080/api/v1/health" ]
      interval: 30s
      timeout: 5s
      retries: 3

volumes:
  allfi-data:
    driver: local
COMPOSE_EOF

echo -e "${GREEN}  ✓ docker-compose.yml 已生成${RESET}"

echo ""

# ==================== 构建并启动 ====================

echo -e "${CYAN}>>> 构建 Docker 镜像并启动服务...${RESET}"
echo -e "  ${YELLOW}首次构建可能需要 1-2 分钟（仅打包，无需编译）...${RESET}"
echo ""

$COMPOSE_CMD up -d --build

# 清理临时构建文件（已打包进镜像）
rm -f ./allfi
rm -rf ./manifest
rm -f ./Dockerfile

# 读取实际端口
PORT=$(grep -E '^ALLFI_PORT=' .env 2>/dev/null | cut -d= -f2 | tr -d '[:space:]' || echo "3174")

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${GREEN}║           AllFi 部署完成！ 🎉                ║${RESET}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${RESET}"
echo ""
echo -e "  🌐 访问地址: ${CYAN}http://localhost:${PORT}${RESET}"
echo -e "  📖 API 文档: ${CYAN}http://localhost:${PORT}/swagger/${RESET}"
echo ""
echo -e "  首次访问需设置 ${YELLOW}PIN 码${RESET}（4-8 位数字）"
echo ""
echo -e "${CYAN}后续维护 (进入部署目录 cd ${DEPLOY_DIR}):${RESET}"
echo -e "  查看日志:   ${GREEN}$COMPOSE_CMD logs -f${RESET}"
echo -e "  停止服务:   ${GREEN}$COMPOSE_CMD down${RESET}"
echo -e "  启动服务:   ${GREEN}$COMPOSE_CMD up -d${RESET}"
echo -e "  重启服务:   ${GREEN}$COMPOSE_CMD restart${RESET}"
echo -e "  数据目录:   ${GREEN}存储在 Docker 数据卷 allfi-data 中${RESET}"
echo ""
echo -e "${CYAN}版本更新:${RESET}"
echo -e "  登录后在 ${YELLOW}系统设置${RESET} 页面点击 ${GREEN}「检查更新」${RESET} 即可 OTA 一键升级到最新版本。"
echo ""
