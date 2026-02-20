#!/bin/bash
# AllFi Docker 部署脚本
# 用法: bash deploy/docker-deploy.sh
# 远程部署: curl -sSL https://raw.githubusercontent.com/your-finance/allfi/main/deploy/docker-deploy.sh | bash

set -e

# 颜色定义
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'

# 构建模式：auto/docker/local
BUILD_MODE="${BUILD_MODE:-auto}"

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
    echo -e "${GREEN}  ✓ docker-compose (standalone) detected${RESET}"
else
    echo -e "${RED}  ✗ 未检测到 Docker Compose，请先安装${RESET}"
    echo "    安装指南: https://docs.docker.com/compose/install/"
    exit 1
fi

# 检查 openssl
if ! command -v openssl >/dev/null 2>&1; then
    echo -e "${RED}  ✗ 未检测到 openssl，请先安装${RESET}"
    exit 1
fi
echo -e "${GREEN}  ✓ openssl detected${RESET}"

# 本地构建所需的工具
if [ "$BUILD_MODE" = "local" ] || [ "$BUILD_MODE" = "auto" ]; then
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${YELLOW}  ⊘ 未检测到 Go，将使用 Docker 构建${RESET}"
        BUILD_MODE="docker"
    else
        echo -e "${GREEN}  ✓ Go $(go version | grep -oE '[0-9]+\.[0-9]+' || echo 'detected')${RESET}"
    fi
fi

if [ "$BUILD_MODE" = "local" ] || [ "$BUILD_MODE" = "auto" ]; then
    if command -v pnpm >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ pnpm detected${RESET}"
    elif command -v npm >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ npm detected${RESET}"
    else
        echo -e "${YELLOW}  ⊘ 未检测到 pnpm/npm，将使用 Docker 构建${RESET}"
        BUILD_MODE="docker"
    fi
fi

echo ""

# ==================== 判断运行模式 ====================

# 如果已在 allfi 项目根目录（有 docker-compose.yml 和 core/ webapp/），直接使用
if [ -f "docker-compose.yml" ] && [ -d "core" ] && [ -d "webapp" ]; then
    echo -e "${CYAN}>>> 检测到已在 AllFi 项目目录，使用本地文件...${RESET}"
    PROJECT_DIR="."
else
    # 远程模式：克隆仓库
    echo -e "${CYAN}>>> 克隆 AllFi 仓库...${RESET}"
    if [ -d "allfi" ]; then
        echo -e "${YELLOW}  ⊘ allfi 目录已存在，跳过克隆${RESET}"
    else
        git clone https://github.com/your-finance/allfi.git
        echo -e "${GREEN}  ✓ 仓库克隆完成${RESET}"
    fi
    PROJECT_DIR="allfi"
fi

cd "$PROJECT_DIR"

echo ""

# ==================== 生成 .env ====================

echo -e "${CYAN}>>> 配置环境变量...${RESET}"

if [ -f ".env" ]; then
    echo -e "${YELLOW}  ⊘ .env 已存在，跳过（如需重新生成，请先删除 .env）${RESET}"
else
    if [ ! -f ".env.example" ]; then
        echo -e "${RED}  ✗ 缺少 .env.example 文件${RESET}"
        exit 1
    fi

    cp .env.example .env

    # 自动生成 MASTER_KEY
    MASTER_KEY=$(openssl rand -base64 32)
    if [ "$(uname)" = "Darwin" ]; then
        sed -i '' "s|CHANGE_ME_USE_openssl_rand_base64_32|${MASTER_KEY}|" .env
    else
        sed -i "s|CHANGE_ME_USE_openssl_rand_base64_32|${MASTER_KEY}|" .env
    fi

    echo -e "${GREEN}  ✓ 已生成 .env${RESET}"
    echo -e "${GREEN}  ✓ ALLFI_MASTER_KEY 已自动生成${RESET}"
fi

echo ""

# ==================== 选择构建模式 ====================

# 获取 Docker 内存信息
DOCKER_MEMORY=$(docker info 2>/dev/null | grep "Total Memory:" | grep -oE '[0-9.]+GiB|[0-9.]+Gi' || echo "8GiB")

# 解析内存数值（单位 GiB）
MEMORY_NUM=$(echo "$DOCKER_MEMORY" | grep -oE '[0-9.]+' | awk '{print int($1)}')

# 自动判断：内存小于 8GiB 或显式指定本地模式时使用本地构建
if [ "$BUILD_MODE" = "auto" ] && [ -n "$MEMORY_NUM" ] && [ "$MEMORY_NUM" -lt 8 ]; then
    echo -e "${YELLOW}>>> Docker 内存不足（${DOCKER_MEMORY}），使用本地构建避免 OOM...${RESET}"
    BUILD_MODE="local"
elif [ "$BUILD_MODE" = "auto" ]; then
    BUILD_MODE="docker"
fi

if [ "$BUILD_MODE" = "local" ]; then
    echo -e "${CYAN}>>> 使用本地构建模式（Docker 仅打包）...${RESET}"
else
    echo -e "${CYAN}>>> 使用 Docker 构建模式...${RESET}"
fi
echo ""

# ==================== 本地构建 ====================

if [ "$BUILD_MODE" = "local" ]; then
    # 获取版本信息
    if [ -f VERSION ]; then
        ALLFI_VERSION=$(cat VERSION)
        export ALLFI_VERSION
    elif git describe --tags --abbrev=0 >/dev/null 2>&1; then
        ALLFI_VERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')
        export ALLFI_VERSION
    else
        ALLFI_VERSION="dev"
        export ALLFI_VERSION
    fi
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
    export GIT_COMMIT
    BUILD_TIME=$(date +%Y-%m-%dT%H:%M:%S)
    export BUILD_TIME

    echo -e "  版本: ${GREEN}${ALLFI_VERSION}${RESET}"
    echo -e "  提交: ${GREEN}${GIT_COMMIT}${RESET}"
    echo ""

    # 构建前端
    echo -e "${CYAN}>>> 构建前端...${RESET}"
    cd webapp
    if command -v pnpm >/dev/null 2>&1; then
        pnpm install && pnpm build
    else
        npm install && npm run build
    fi
    cd ..
    echo -e "${GREEN}  ✓ 前端构建完成${RESET}"
    echo ""

    # 复制前端产物
    echo -e "${CYAN}>>> 复制前端产物（用于 Go embed）...${RESET}"
    rm -rf core/internal/statics/dist/*
    cp -R webapp/dist/* core/internal/statics/dist/
    echo -e "${GREEN}  ✓ 前端产物已复制到 core/internal/statics/dist/${RESET}"
    echo ""

    # 构建后端（交叉编译到 Linux）
    echo -e "${CYAN}>>> 构建后端...${RESET}"
    cd core
    CGO_ENABLED=0 GOOS=linux GOARCH=$(uname -m | sed 's/arm64/arm64/;s/x86_64/amd64/') GOMAXPROCS=2 GOGC=20 go build \
      -p=1 \
      -ldflags="-s -w \
      -X your-finance/allfi/internal/version.Version=${ALLFI_VERSION} \
      -X your-finance/allfi/internal/version.BuildTime=${BUILD_TIME} \
      -X your-finance/allfi/internal/version.GitCommit=${GIT_COMMIT}" \
      -o allfi cmd/server/main.go
    cd ..
    echo -e "${GREEN}  ✓ 后端构建完成${RESET}"
    echo ""

    # 构建 Docker 镜像（仅打包）
    echo -e "${CYAN}>>> 构建 Docker 镜像（仅打包）...${RESET}"
    docker build -t allfi-backend:latest -f - . <<'DOCKERFILE'
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata wget
RUN addgroup -S allfi && adduser -S allfi -G allfi
WORKDIR /app
COPY core/allfi ./allfi
COPY core/manifest/config/config.yaml manifest/config/config.yaml
RUN mkdir -p /app/data /app/logs && chown -R allfi:allfi /app
USER allfi
EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["./allfi"]
DOCKERFILE
    echo -e "${GREEN}  ✓ Docker 镜像构建完成${RESET}"
    echo ""
fi

# ==================== 启动服务 ====================

echo -e "${CYAN}>>> 启动 Docker 服务...${RESET}"
if [ "$BUILD_MODE" = "local" ]; then
    echo -e "  ${GREEN}使用预构建镜像启动...${RESET}"
else
    echo -e "  ${YELLOW}首次启动需要构建镜像，可能需要几分钟...${RESET}"
    echo ""
fi

$COMPOSE_CMD up -d ${BUILD_MODE:+--build}

# 从 .env 文件读取实际端口配置（docker-compose 使用 .env，脚本也需要同步）
if [ -f ".env" ]; then
    _ALLFI_PORT=$(grep -E '^ALLFI_PORT=' .env | cut -d= -f2 | tr -d '[:space:]')
fi
_ALLFI_PORT="${_ALLFI_PORT:-5173}"

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${GREEN}║           AllFi 部署完成！                   ║${RESET}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${RESET}"
echo ""
echo -e "  访问地址: ${CYAN}http://localhost:${_ALLFI_PORT}${RESET}"
echo -e "  API 文档: ${CYAN}http://localhost:${_ALLFI_PORT}/swagger/${RESET}"
echo ""
echo -e "  首次访问需设置 ${YELLOW}PIN 码${RESET}（4-8 位数字）"
echo ""
echo -e "${CYAN}常用命令:${RESET}"
echo -e "  查看日志:   ${GREEN}$COMPOSE_CMD logs -f${RESET}"
echo -e "  停止服务:   ${GREEN}$COMPOSE_CMD down${RESET}"
echo -e "  重启服务:   ${GREEN}$COMPOSE_CMD restart${RESET}"
echo -e "  重新构建:   ${GREEN}$COMPOSE_CMD up -d --build${RESET}"
echo ""
