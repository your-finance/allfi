#!/bin/bash
# AllFi 本地构建脚本（用于 Docker 内存不足的情况）
# 用法: bash deploy/build-local.sh

set -e

# 颜色定义
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${CYAN}║       AllFi 本地构建（避免 Docker OOM）      ║${RESET}"
echo -e "${CYAN}╚══════════════════════════════════════════════╝${RESET}"
echo ""

# ==================== 前端构建 ====================

echo -e "${CYAN}>>> 构建前端...${RESET}"
cd webapp
if command -v pnpm >/dev/null 2>&1; then
    pnpm install && pnpm build
elif command -v npm >/dev/null 2>&1; then
    npm install && npm run build
else
    echo -e "${RED}  ✗ 未检测到 pnpm 或 npm${RESET}"
    exit 1
fi
cd ..
echo -e "${GREEN}  ✓ 前端构建完成${RESET}"
echo ""

# ==================== 复制前端产物（用于 embed）====================

echo -e "${CYAN}>>> 复制前端产物（用于 Go embed）...${RESET}"
rm -rf core/internal/statics/dist/*
cp -R webapp/dist/* core/internal/statics/dist/
echo -e "${GREEN}  ✓ 前端产物已复制到 core/internal/statics/dist/${RESET}"
echo ""

# ==================== 后端构建 ====================

echo -e "${CYAN}>>> 构建后端...${RESET}"

# 获取版本信息
if [ -f VERSION ]; then
    export ALLFI_VERSION=$(cat VERSION)
elif git describe --tags --abbrev=0 >/dev/null 2>&1; then
    export ALLFI_VERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')
else
    export ALLFI_VERSION="dev"
fi
export GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
export BUILD_TIME=$(date +%Y-%m-%dT%H:%M:%S)

cd core

# 使用较低内存的构建参数，交叉编译到 Linux ARM64
echo -e "  版本: ${GREEN}${ALLFI_VERSION}${RESET}"
echo -e "  提交: ${GREEN}${GIT_COMMIT}${RESET}"
echo ""

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOMAXPROCS=2 GOGC=20 go build \
  -p=1 \
  -ldflags="-s -w \
  -X your-finance/allfi/internal/version.Version=${ALLFI_VERSION} \
  -X your-finance/allfi/internal/version.BuildTime=${BUILD_TIME} \
  -X your-finance/allfi/internal/version.GitCommit=${GIT_COMMIT}" \
  -o allfi cmd/server/main.go

cd ..
echo -e "${GREEN}  ✓ 后端构建完成${RESET}"
echo ""

# ==================== Docker 镜像构建 ====================

echo -e "${CYAN}>>> 构建 Docker 镜像（仅打包）...${RESET}"

# 使用本地构建的二进制
docker build -t allfi-backend:latest -f - . <<'DOCKERFILE'
# 使用本地构建的二进制打包到镜像
FROM alpine:3.21

# 安装基础工具 (证书、时区)
RUN apk add --no-cache ca-certificates tzdata wget

# 创建非特权用户
RUN addgroup -S allfi && adduser -S allfi -G allfi

WORKDIR /app

# 从本地复制二进制（包含嵌入的前端资源）
COPY core/allfi ./allfi
# 复制配置模板
COPY core/manifest/config/config.yaml manifest/config/config.yaml

# 设置权限
RUN mkdir -p /app/data /app/logs && chown -R allfi:allfi /app

# 切换用户
USER allfi

# 暴露端口
EXPOSE 8080

# 数据持久化
VOLUME ["/app/data"]

ENTRYPOINT ["./allfi"]
DOCKERFILE

echo -e "${GREEN}  ✓ Docker 镜像构建完成${RESET}"
echo ""

echo -e "${GREEN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${GREEN}║           构建完成！                           ║${RESET}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${RESET}"
echo ""
echo -e "${CYAN}下一步: docker compose up -d${RESET}"
echo ""
